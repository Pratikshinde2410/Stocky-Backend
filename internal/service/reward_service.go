package service

import (
    "context"
    "fmt"
    "github.com/shopspring/decimal"
    "stocky-backend/internal/models"
    "stocky-backend/internal/repository"
)

type RewardService interface {
    CreateReward(ctx context.Context, req *models.RewardRequest) (*models.RewardResponseData, error)
    GetTodayStocks(ctx context.Context, userID string) (interface{}, error)
}

type rewardService struct {
    rewardRepo    repository.RewardRepository
    ledgerService LedgerService
    pricingService PricingService
}

func NewRewardService(
    rewardRepo repository.RewardRepository,
    ledgerSvc LedgerService,
    pricingSvc PricingService,
) RewardService {
    return &rewardService{
        rewardRepo:     rewardRepo,
        ledgerService:  ledgerSvc,
        pricingService: pricingSvc,
    }
}

func (s *rewardService) CreateReward(ctx context.Context, req *models.RewardRequest) (*models.RewardResponseData, error) {
    existing, err := s.rewardRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
    if err != nil {
        return nil, fmt.Errorf("failed to check idempotency: %w", err)
    }
    if existing != nil {
        return nil, fmt.Errorf("duplicate_reward")
    }

    currentPrice, err := s.pricingService.GetCurrentPrice(ctx, req.StockSymbol)
    if err != nil {
        return nil, fmt.Errorf("failed to get stock price: %w", err)
    }

    totalValue := req.Shares.Mul(currentPrice)
    fees := s.calculateFees(totalValue)

    reward := &models.StockReward{
        UserID:          req.UserID,
        StockSymbol:     req.StockSymbol,
        Shares:          req.Shares,
        RewardType:      req.RewardType,
        PriceAtReward:   currentPrice,
        TotalStockValue: totalValue,
        IdempotencyKey:  req.IdempotencyKey,
        RewardTimestamp: *req.Timestamp,
    }

    if err := s.rewardRepo.Create(ctx, reward); err != nil {
        return nil, fmt.Errorf("failed to create reward: %w", err)
    }

    if err := s.ledgerService.RecordRewardTransaction(ctx, reward, fees); err != nil {
        return nil, fmt.Errorf("failed to record ledger entries: %w", err)
    }

    return &models.RewardResponseData{
        RewardID:     reward.RewardID,
        UserID:       reward.UserID,
        StockSymbol:  reward.StockSymbol,
        Shares:       reward.Shares,
        CurrentPrice: currentPrice,
        TotalValue:   totalValue,
        Fees:         fees,
        Timestamp:    reward.RewardTimestamp,
    }, nil
}

func (s *rewardService) GetTodayStocks(ctx context.Context, userID string) (interface{}, error) {
    return s.rewardRepo.GetTodayRewards(ctx, userID)
}

func (s *rewardService) calculateFees(totalValue decimal.Decimal) models.FeeBreakdown {
    brokerage := totalValue.Mul(decimal.NewFromFloat(0.0005))
    stt := totalValue.Mul(decimal.NewFromFloat(0.000125))
    gst := brokerage.Mul(decimal.NewFromFloat(0.18))
    total := brokerage.Add(stt).Add(gst)
    return models.FeeBreakdown{
        Brokerage: brokerage.Round(2),
        STT:       stt.Round(2),
        GST:       gst.Round(2),
        Total:     total.Round(2),
    }
}



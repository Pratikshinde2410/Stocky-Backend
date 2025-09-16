package service

import (
    "context"
    "stocky-backend/internal/models"
)

type LedgerService interface {
    RecordRewardTransaction(ctx context.Context, reward *models.StockReward, fees models.FeeBreakdown) error
}

type ledgerService struct{}

func NewLedgerService() LedgerService { return &ledgerService{} }

func (l *ledgerService) RecordRewardTransaction(ctx context.Context, reward *models.StockReward, fees models.FeeBreakdown) error {
    return nil
}



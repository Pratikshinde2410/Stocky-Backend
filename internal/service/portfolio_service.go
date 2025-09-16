package service

import (
    "context"
    "stocky-backend/internal/repository"
)

type PortfolioService interface {
    GetHistoricalINR(ctx context.Context, userID string) (interface{}, error)
    GetStats(ctx context.Context, userID string) (interface{}, error)
    GetPortfolio(ctx context.Context, userID string) (interface{}, error)
}

type portfolioService struct{ repo repository.PortfolioRepository }

func NewPortfolioService(repo repository.PortfolioRepository) PortfolioService { return &portfolioService{repo: repo} }

func (s *portfolioService) GetHistoricalINR(ctx context.Context, userID string) (interface{}, error) {
    return s.repo.GetHistoricalINR(ctx, userID)
}

func (s *portfolioService) GetStats(ctx context.Context, userID string) (interface{}, error) {
    return s.repo.GetStats(ctx, userID)
}

func (s *portfolioService) GetPortfolio(ctx context.Context, userID string) (interface{}, error) {
    return s.repo.GetPortfolio(ctx, userID)
}



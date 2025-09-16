package repository

import "context"

type PortfolioRepository interface {
    GetHistoricalINR(ctx context.Context, userID string) (interface{}, error)
    GetStats(ctx context.Context, userID string) (interface{}, error)
    GetPortfolio(ctx context.Context, userID string) (interface{}, error)
}



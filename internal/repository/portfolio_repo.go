package repository

import (
    "context"
    "github.com/jmoiron/sqlx"
)

type portfolioRepo struct{ db *sqlx.DB }

func NewPortfolioRepository(db *sqlx.DB) PortfolioRepository { return &portfolioRepo{db: db} }

func (r *portfolioRepo) GetHistoricalINR(ctx context.Context, userID string) (interface{}, error) {
    return nil, nil
}

func (r *portfolioRepo) GetStats(ctx context.Context, userID string) (interface{}, error) {
    return nil, nil
}

func (r *portfolioRepo) GetPortfolio(ctx context.Context, userID string) (interface{}, error) {
    return nil, nil
}



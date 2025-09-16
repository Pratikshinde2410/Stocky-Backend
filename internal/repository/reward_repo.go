package repository

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/jmoiron/sqlx"
    "stocky-backend/internal/models"
)

type RewardRepository interface {
    Create(ctx context.Context, reward *models.StockReward) error
    GetByIdempotencyKey(ctx context.Context, key string) (*models.StockReward, error)
    GetTodayRewards(ctx context.Context, userID string) ([]*models.StockReward, error)
}

type rewardRepository struct {
    db *sqlx.DB
}

func NewRewardRepository(db *sqlx.DB) RewardRepository {
    return &rewardRepository{db: db}
}

func (r *rewardRepository) Create(ctx context.Context, reward *models.StockReward) error {
    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()

    query := `
        INSERT INTO stock_rewards 
        (user_id, stock_symbol, shares, reward_type, price_at_reward, 
         total_stock_value, idempotency_key, reward_timestamp)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING reward_id, created_at`

    err = tx.QueryRowxContext(ctx, query,
        reward.UserID, reward.StockSymbol, reward.Shares,
        reward.RewardType, reward.PriceAtReward, reward.TotalStockValue,
        reward.IdempotencyKey, reward.RewardTimestamp,
    ).Scan(&reward.RewardID, &reward.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to insert reward: %w", err)
    }

    return tx.Commit()
}

func (r *rewardRepository) GetByIdempotencyKey(ctx context.Context, key string) (*models.StockReward, error) {
    var reward models.StockReward
    query := `
        SELECT reward_id, user_id, stock_symbol, shares, reward_type,
               price_at_reward, total_stock_value, idempotency_key,
               reward_timestamp, created_at
        FROM stock_rewards 
        WHERE idempotency_key = $1`

    err := r.db.GetContext(ctx, &reward, query, key)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get reward by idempotency key: %w", err)
    }
    return &reward, nil
}

func (r *rewardRepository) GetTodayRewards(ctx context.Context, userID string) ([]*models.StockReward, error) {
    var rewards []*models.StockReward
    query := `
        SELECT sr.reward_id, sr.user_id, sr.stock_symbol, sr.shares,
               sr.reward_type, sr.price_at_reward, sr.total_stock_value,
               sr.reward_timestamp
        FROM stock_rewards sr
        WHERE sr.user_id = $1 
        AND DATE(sr.reward_timestamp) = CURRENT_DATE
        ORDER BY sr.reward_timestamp DESC`

    err := r.db.SelectContext(ctx, &rewards, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get today rewards: %w", err)
    }
    return rewards, nil
}



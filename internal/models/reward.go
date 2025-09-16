package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

type RewardType string

const (
    RewardTypeOnboarding       RewardType = "ONBOARDING"
    RewardTypeReferral         RewardType = "REFERRAL"
    RewardTypeTradingMilestone RewardType = "TRADING_MILESTONE"
    RewardTypeBonus            RewardType = "BONUS"
)

type StockReward struct {
    RewardID        uuid.UUID       `json:"rewardId" db:"reward_id"`
    UserID          string          `json:"userId" db:"user_id"`
    StockSymbol     string          `json:"stockSymbol" db:"stock_symbol"`
    Shares          decimal.Decimal `json:"shares" db:"shares"`
    RewardType      RewardType      `json:"rewardType" db:"reward_type"`
    PriceAtReward   decimal.Decimal `json:"priceAtReward" db:"price_at_reward"`
    TotalStockValue decimal.Decimal `json:"totalStockValue" db:"total_stock_value"`
    IdempotencyKey  string          `json:"idempotencyKey" db:"idempotency_key"`
    RewardTimestamp time.Time       `json:"rewardTimestamp" db:"reward_timestamp"`
    CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
}

type RewardRequest struct {
    UserID         string          `json:"userId" validate:"required"`
    StockSymbol    string          `json:"stockSymbol" validate:"required"`
    Shares         decimal.Decimal `json:"shares" validate:"required,gt=0"`
    RewardType     RewardType      `json:"rewardType" validate:"required"`
    IdempotencyKey string          `json:"idempotencyKey" validate:"required"`
    Timestamp      *time.Time      `json:"timestamp"`
}

type RewardResponse struct {
    Success bool                `json:"success"`
    Data    *RewardResponseData `json:"data,omitempty"`
    Error   string              `json:"error,omitempty"`
    Message string              `json:"message,omitempty"`
}

type RewardResponseData struct {
    RewardID     uuid.UUID       `json:"rewardId"`
    UserID       string          `json:"userId"`
    StockSymbol  string          `json:"stockSymbol"`
    Shares       decimal.Decimal `json:"shares"`
    CurrentPrice decimal.Decimal `json:"currentPrice"`
    TotalValue   decimal.Decimal `json:"totalValue"`
    Fees         FeeBreakdown    `json:"fees"`
    Timestamp    time.Time       `json:"timestamp"`
}

type FeeBreakdown struct {
    Brokerage decimal.Decimal `json:"brokerage"`
    STT       decimal.Decimal `json:"stt"`
    GST       decimal.Decimal `json:"gst"`
    Total     decimal.Decimal `json:"total"`
}



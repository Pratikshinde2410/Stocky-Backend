package service

import (
    "context"
    "github.com/shopspring/decimal"
)

type PricingService interface {
    GetCurrentPrice(ctx context.Context, symbol string) (decimal.Decimal, error)
}

type pricingService struct{}

func NewPricingService() PricingService { return &pricingService{} }

func (p *pricingService) GetCurrentPrice(ctx context.Context, symbol string) (decimal.Decimal, error) {
    return decimal.NewFromFloat(100.00), nil
}



package models

import "github.com/shopspring/decimal"

type PortfolioHolding struct {
    StockSymbol  string          `json:"stockSymbol"`
    TotalShares  decimal.Decimal `json:"totalShares"`
    CurrentPrice decimal.Decimal `json:"currentPrice"`
    TotalValue   decimal.Decimal `json:"totalValue"`
}




package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "stocky-backend/internal/api/handlers"
    "stocky-backend/internal/config"
    "stocky-backend/internal/db"
    "stocky-backend/internal/repository"
    "stocky-backend/internal/service"
    "stocky-backend/pkg/logger"
)

func main() {
    cfg, err := config.Load()
    if err != nil { log.Fatal("Failed to load config:", err) }

    logg := logger.New()

    database, err := db.NewPostgresConnection(cfg.Database)
    if err != nil { log.Fatal("Failed to connect to database:", err) }
    defer database.Close()

    rewardRepo := repository.NewRewardRepository(database)
    portfolioRepo := repository.NewPortfolioRepository(database)

    pricingService := service.NewPricingService()
    ledgerService := service.NewLedgerService()
    rewardService := service.NewRewardService(rewardRepo, ledgerService, pricingService)
    portfolioService := service.NewPortfolioService(portfolioRepo)

    rewardHandler := handlers.NewRewardHandler(rewardService, logg)
    portfolioHandler := handlers.NewPortfolioHandler(portfolioService, logg)

    router := gin.New()
    router.Use(gin.Recovery())

    api := router.Group("/api/v1")
    {
        api.POST("/reward", rewardHandler.CreateReward)
        api.GET("/today-stocks/:userId", rewardHandler.GetTodayStocks)
        api.GET("/historical-inr/:userId", portfolioHandler.GetHistoricalINR)
        api.GET("/stats/:userId", portfolioHandler.GetStats)
        api.GET("/portfolio/:userId", portfolioHandler.GetPortfolio)
    }

    srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: router, ReadTimeout: time.Duration(cfg.Server.ReadTimeout) * time.Second, WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second, MaxHeaderBytes: 1 << 20}

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()

    logg.Info("Server started", "port", cfg.Server.Port)

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logg.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil { log.Fatal("Server forced to shutdown:", err) }
    logg.Info("Server exited")
}



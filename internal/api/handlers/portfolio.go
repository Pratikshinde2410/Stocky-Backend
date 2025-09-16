package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "stocky-backend/internal/service"
    "stocky-backend/pkg/logger"
)

type PortfolioHandler struct {
    service service.PortfolioService
    logger  logger.Logger
}

func NewPortfolioHandler(s service.PortfolioService, l logger.Logger) *PortfolioHandler {
    return &PortfolioHandler{service: s, logger: l}
}

func (h *PortfolioHandler) GetHistoricalINR(c *gin.Context) {
    userID := c.Param("userId")
    data, err := h.service.GetHistoricalINR(c.Request.Context(), userID)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false}); return }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PortfolioHandler) GetStats(c *gin.Context) {
    userID := c.Param("userId")
    data, err := h.service.GetStats(c.Request.Context(), userID)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false}); return }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
    userID := c.Param("userId")
    data, err := h.service.GetPortfolio(c.Request.Context(), userID)
    if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false}); return }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}



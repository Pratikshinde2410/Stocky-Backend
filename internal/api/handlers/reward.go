package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "stocky-backend/internal/models"
    "stocky-backend/internal/service"
    "stocky-backend/pkg/logger"
)

type RewardHandler struct {
    rewardService service.RewardService
    logger        logger.Logger
    validator     *validator.Validate
}

func NewRewardHandler(rs service.RewardService, l logger.Logger) *RewardHandler {
    return &RewardHandler{rewardService: rs, logger: l, validator: validator.New()}
}

func (h *RewardHandler) CreateReward(c *gin.Context) {
    var req models.RewardRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        h.logger.Error("Invalid request payload", "error", err)
        c.JSON(http.StatusBadRequest, models.RewardResponse{Success: false, Error: "INVALID_REQUEST", Message: err.Error()})
        return
    }
    if err := h.validator.Struct(&req); err != nil {
        h.logger.Error("Validation failed", "error", err)
        c.JSON(http.StatusBadRequest, models.RewardResponse{Success: false, Error: "VALIDATION_ERROR", Message: err.Error()})
        return
    }
    if req.Timestamp == nil {
        now := time.Now()
        req.Timestamp = &now
    }
    reward, err := h.rewardService.CreateReward(c.Request.Context(), &req)
    if err != nil {
        switch err.Error() {
        case "duplicate_reward":
            c.JSON(http.StatusConflict, models.RewardResponse{Success: false, Error: "DUPLICATE_REWARD", Message: "Reward already processed for this idempotency key"})
        default:
            h.logger.Error("Failed to create reward", "error", err)
            c.JSON(http.StatusInternalServerError, models.RewardResponse{Success: false, Error: "INTERNAL_ERROR", Message: "Failed to process reward"})
        }
        return
    }
    c.JSON(http.StatusCreated, models.RewardResponse{Success: true, Data: reward, Message: "Reward recorded successfully"})
}

func (h *RewardHandler) GetTodayStocks(c *gin.Context) {
    userID := c.Param("userId")
    if userID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "INVALID_USER_ID", "message": "User ID is required"})
        return
    }
    stocks, err := h.rewardService.GetTodayStocks(c.Request.Context(), userID)
    if err != nil {
        h.logger.Error("Failed to get today stocks", "error", err, "userId", userID)
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "INTERNAL_ERROR", "message": "Failed to retrieve stocks"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": stocks})
}



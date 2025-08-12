package handlers

import (
	"net/http"
	"video-ad-tracker/internal/middleware"
	"video-ad-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AdServiceInterface interface {
	GetAllAds() ([]models.Ad, error)
	GetAdByID(id int) (*models.Ad, error)
}

type AnalyticsServiceInterface interface {
	GetAnalytics(timeFrame string) ([]models.Analytics, error)
	GetHourlyBreakdown() ([]models.Analytics, error)
}

// ClickServiceInterface defines the interface for click operations
type ClickServiceInterface interface {
	RecordClick(req models.ClickRequest, clientIP string) error
}
type Handlers struct {
	adService        AdServiceInterface
	analyticsService AnalyticsServiceInterface
	clickService     ClickServiceInterface
	logger           *logrus.Logger
}

// Create new handlers
func NewHandlers(adService AdServiceInterface, analyticsService AnalyticsServiceInterface, clickService ClickServiceInterface, logger *logrus.Logger) *Handlers {
	return &Handlers{
		adService:        adService,
		analyticsService: analyticsService,
		clickService:     clickService,
		logger:           logger,
	}
}

func Routes(router *gin.Engine, adService AdServiceInterface, analyticsService AnalyticsServiceInterface, clickService ClickServiceInterface) {
	logger := logrus.New()
	handlers := NewHandlers(adService, analyticsService, clickService, logger)

	api := router.Group("/api/v1")
	{
		api.GET("/ads", handlers.GetAds)
		api.POST("/ads/click", handlers.RecordClick)
		api.GET("/ads/analytics", handlers.GetAnalytics)
		api.GET("/ads/analytics/hourly", handlers.GetHourlyAnalytics)
	}

	router.GET("/metrics", middleware.MetricsHandler())
}

// Get all advertisements
func (h *Handlers) GetAds(c *gin.Context) {
	ads, err := h.adService.GetAllAds()
	if err != nil {
		h.logger.Errorf("Failed to get ads: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve ads",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    ads,
	})
}

// Record click event
func (h *Handlers) RecordClick(c *gin.Context) {
	var req models.ClickRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Invalid click request: %v", err)
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = c.GetHeader("X-Forwarded-For")
	}

	err := h.clickService.RecordClick(req, clientIP)
	if err != nil {
		h.logger.Errorf("Failed to record click: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to record click",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    map[string]string{"message": "Click recorded successfully"},
	})
}

// Get analytics data
func (h *Handlers) GetAnalytics(c *gin.Context) {
	timeFrame := c.DefaultQuery("timeframe", "24h")

	validTimeFrames := map[string]bool{"15m": true, "30m": true, "1h": true, "6h": true, "12h": true, "24h": true, "7d": true, "30d": true}
	if !validTimeFrames[timeFrame] {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid timeframe. Use: 15m, 30m, 1h, 6h, 12h, 24h, 7d, or 30d",
		})
		return
	}

	analytics, err := h.analyticsService.GetAnalytics(timeFrame)
	if err != nil {
		h.logger.Errorf("Failed to get analytics: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve analytics",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    analytics,
	})
}

// Get hourly analytics breakdown
func (h *Handlers) GetHourlyAnalytics(c *gin.Context) {
	analytics, err := h.analyticsService.GetHourlyBreakdown()
	if err != nil {
		h.logger.Errorf("Failed to get hourly analytics: %v", err)
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to retrieve hourly analytics",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    analytics,
	})
}

package services

import (
	"database/sql"
	"fmt"
	"time"
	"video-ad-tracker/internal/models"

	"github.com/sirupsen/logrus"
)

type AnalyticsService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// Create new analytics service
func NewAnalyticsService(db *sql.DB, logger *logrus.Logger) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger,
	}
}

// Get real-time analytics for ads
func (s *AnalyticsService) GetAnalytics(timeFrame string) ([]models.Analytics, error) {
	// Process pending clicks first
	if err := s.processUnprocessedClicks(); err != nil {
		s.logger.Errorf("Failed to process unprocessed clicks: %v", err)
	}

	// Get time window
	timeWindow := s.getTimeWindow(timeFrame)

	query := `
		SELECT 
			a.id as ad_id,
			COUNT(ce.id) as total_clicks,
			COALESCE(AVG(ce.video_playback_time), 0.0) as avg_playback_time
		FROM ads a
		LEFT JOIN click_events ce ON a.id = ce.ad_id 
			AND ce.timestamp >= $1::timestamp
		GROUP BY a.id, a.title
		ORDER BY total_clicks DESC NULLS LAST, a.id ASC
	`

	rows, err := s.db.Query(query, timeWindow)
	if err != nil {
		s.logger.Errorf("Failed to query analytics: %v", err)
		return nil, err
	}
	defer rows.Close()

	var analytics []models.Analytics
	for rows.Next() {
		var analytic models.Analytics
		err := rows.Scan(
			&analytic.AdID,
			&analytic.TotalClicks,
			&analytic.AvgPlaybackTime,
		)
		if err != nil {
			s.logger.Errorf("Failed to scan analytics: %v", err)
			return nil, err
		}

		// Calculate click-through rate
		analytic.CTR = s.calculateCTR(analytic.AdID, timeWindow)
		analytic.TimeFrame = timeFrame
		analytic.LastUpdated = time.Now()

		analytics = append(analytics, analytic)
	}

	return analytics, nil
}

// Get hourly breakdown for the last 24 hours
func (s *AnalyticsService) GetHourlyBreakdown() ([]models.Analytics, error) {
	// Process pending clicks first
	if err := s.processUnprocessedClicks(); err != nil {
		s.logger.Errorf("Failed to process unprocessed clicks: %v", err)
	}

	query := `
		WITH hourly_data AS (
			SELECT 
				a.id as ad_id,
				COUNT(ce.id) as total_clicks,
				COALESCE(AVG(ce.video_playback_time), 0.0) as avg_playback_time,
				EXTRACT(hour FROM ce.timestamp) as hour
			FROM ads a
			LEFT JOIN click_events ce ON a.id = ce.ad_id 
				AND ce.timestamp >= NOW() - INTERVAL '24 hours'
			GROUP BY a.id, EXTRACT(hour FROM ce.timestamp)
		)
		SELECT 
			ad_id,
			total_clicks,
			avg_playback_time,
			COALESCE(hour, -1) as hour
		FROM hourly_data
		ORDER BY ad_id ASC, hour ASC NULLS LAST
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Failed to query hourly breakdown: %v", err)
		return nil, err
	}
	defer rows.Close()

	var analytics []models.Analytics
	for rows.Next() {
		var analytic models.Analytics
		var hour int
		err := rows.Scan(
			&analytic.AdID,
			&analytic.TotalClicks,
			&analytic.AvgPlaybackTime,
			&hour,
		)
		if err != nil {
			s.logger.Errorf("Failed to scan hourly breakdown: %v", err)
			return nil, err
		}

		analytic.CTR = s.calculateCTR(analytic.AdID, time.Now().Add(-24*time.Hour))
		if hour == -1 {
			analytic.TimeFrame = "no_clicks"
		} else {
			analytic.TimeFrame = fmt.Sprintf("hour_%d", hour)
		}
		analytic.LastUpdated = time.Now()

		analytics = append(analytics, analytic)
	}

	return analytics, nil
}

// Process unprocessed click events
func (s *AnalyticsService) processUnprocessedClicks() error {
	clicks, err := s.getUnprocessedClicks()
	if err != nil {
		return err
	}

	for _, click := range clicks {
		if err := s.markClickProcessed(click.ID); err != nil {
			s.logger.Errorf("Failed to mark click %d as processed: %v", click.ID, err)
			continue
		}
	}

	return nil
}

// Get unprocessed click events
func (s *AnalyticsService) getUnprocessedClicks() ([]models.ClickEvent, error) {
	query := `
		SELECT id, ad_id, timestamp, ip_address, video_playback_time, user_agent, processed
		FROM click_events 
		WHERE processed = false
		ORDER BY timestamp ASC, id ASC
		LIMIT 1000
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clicks []models.ClickEvent
	for rows.Next() {
		var click models.ClickEvent
		err := rows.Scan(
			&click.ID,
			&click.AdID,
			&click.Timestamp,
			&click.IPAddress,
			&click.VideoPlaybackTime,
			&click.UserAgent,
			&click.Processed,
		)
		if err != nil {
			return nil, err
		}
		clicks = append(clicks, click)
	}

	return clicks, nil
}

// Mark click event as processed
func (s *AnalyticsService) markClickProcessed(clickID int) error {
	query := "UPDATE click_events SET processed = true, updated_at = NOW() WHERE id = $1 AND processed = false"
	_, err := s.db.Exec(query, clickID)
	return err
}

// Calculate click-through rate for an ad
func (s *AnalyticsService) calculateCTR(adID int, timeWindow time.Time) float64 {
	// Basic CTR calculation
	query := `
		SELECT COUNT(*)::bigint
		FROM click_events 
		WHERE ad_id = $1 AND timestamp >= $2::timestamp
	`

	var clicks int
	err := s.db.QueryRow(query, adID, timeWindow).Scan(&clicks)
	if err != nil {
		s.logger.Errorf("Failed to count clicks for ad %d: %v", adID, err)
		return 0
	}

	// Assume 1000 impressions per ad for demo
	impressions := 1000.0
	if impressions > 0 {
		return float64(clicks) / impressions * 100
	}
	return 0
}

// Get time window based on time frame
func (s *AnalyticsService) getTimeWindow(timeFrame string) time.Time {
	now := time.Now()

	switch timeFrame {
	case "15m":
		return now.Add(-15 * time.Minute)
	case "30m":
		return now.Add(-30 * time.Minute)
	case "1h":
		return now.Add(-1 * time.Hour)
	case "6h":
		return now.Add(-6 * time.Hour)
	case "12h":
		return now.Add(-12 * time.Hour)
	case "24h":
		return now.Add(-24 * time.Hour)
	case "7d":
		return now.Add(-7 * 24 * time.Hour)
	case "30d":
		return now.Add(-30 * 24 * time.Hour)
	default:
		return now.Add(-24 * time.Hour) // Default 24 hours
	}
}

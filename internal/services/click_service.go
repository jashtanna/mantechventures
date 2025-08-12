package services

import (
	"database/sql"
	"time"
	"video-ad-tracker/internal/models"

	"github.com/sirupsen/logrus"
)

type ClickService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// Create new click service
func NewClickService(db *sql.DB, logger *logrus.Logger) *ClickService {
	return &ClickService{
		db:     db,
		logger: logger,
	}
}

// Record click event asynchronously
func (s *ClickService) RecordClick(req models.ClickRequest, clientIP string) error {
	// Process click in background
	go s.processClickAsync(req, clientIP)

	// Return immediately
	return nil
}

// Process click event asynchronously
func (s *ClickService) processClickAsync(req models.ClickRequest, clientIP string) {
	// Check if ad exists
	ad, err := s.validateAd(req.AdID)
	if err != nil {
		s.logger.Errorf("Failed to validate ad %d: %v", req.AdID, err)
		return
	}
	if ad == nil {
		s.logger.Errorf("Ad %d not found", req.AdID)
		return
	}

	// Save click event
	query := `
		INSERT INTO click_events (ad_id, timestamp, ip_address, video_playback_time, user_agent, processed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`

	var clickID int
	err = s.db.QueryRow(query,
		req.AdID,
		time.Now(),
		clientIP,
		req.VideoPlaybackTime,
		req.UserAgent,
		false, // Processed by analytics service
	).Scan(&clickID)

	if err != nil {
		s.logger.Errorf("Failed to insert click event: %v", err)
		return
	}

	s.logger.Infof("Click event recorded for ad %d", req.AdID)
}

// Validate ad exists
func (s *ClickService) validateAd(adID int) (*models.Ad, error) {
	query := "SELECT id FROM ads WHERE id = $1"
	var id int
	err := s.db.QueryRow(query, adID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &models.Ad{ID: id}, nil
}

// Get unprocessed click events
func (s *ClickService) GetUnprocessedClicks() ([]models.ClickEvent, error) {
	query := `
		SELECT id, ad_id, timestamp, ip_address, video_playback_time, user_agent, processed
		FROM click_events 
		WHERE processed = false
		ORDER BY timestamp ASC
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
func (s *ClickService) MarkClickProcessed(clickID int) error {
	query := "UPDATE click_events SET processed = true WHERE id = $1"
	_, err := s.db.Exec(query, clickID)
	return err
}

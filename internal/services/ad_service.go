package services

import (
	"database/sql"
	"video-ad-tracker/internal/models"

	"github.com/sirupsen/logrus"
)

type AdService struct {
	db     *sql.DB
	logger *logrus.Logger
}

// Create new ad service
func NewAdService(db *sql.DB, logger *logrus.Logger) *AdService {
	return &AdService{
		db:     db,
		logger: logger,
	}
}

// Get all advertisements
func (s *AdService) GetAllAds() ([]models.Ad, error) {
	query := `
		SELECT id, image_url, target_url, title, description, created_at, updated_at 
		FROM ads 
		ORDER BY created_at DESC, id ASC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		s.logger.Errorf("Failed to query ads: %v", err)
		return nil, err
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		err := rows.Scan(
			&ad.ID,
			&ad.ImageURL,
			&ad.TargetURL,
			&ad.Title,
			&ad.Description,
			&ad.CreatedAt,
			&ad.UpdatedAt,
		)
		if err != nil {
			s.logger.Errorf("Failed to scan ad: %v", err)
			return nil, err
		}
		ads = append(ads, ad)
	}

	return ads, nil
}

// Get advertisement by ID
func (s *AdService) GetAdByID(id int) (*models.Ad, error) {
	query := `
		SELECT id, image_url, target_url, title, description, created_at, updated_at 
		FROM ads 
		WHERE id = $1::integer
	`

	var ad models.Ad
	err := s.db.QueryRow(query, id).Scan(
		&ad.ID,
		&ad.ImageURL,
		&ad.TargetURL,
		&ad.Title,
		&ad.Description,
		&ad.CreatedAt,
		&ad.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		s.logger.Errorf("Failed to get ad by ID %d: %v", id, err)
		return nil, err
	}

	return &ad, nil
}

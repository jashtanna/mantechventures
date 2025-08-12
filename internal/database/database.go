package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Create new database connection
func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// Create database tables
func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS ads (
			id SERIAL PRIMARY KEY,
			image_url VARCHAR(500) NOT NULL,
			target_url VARCHAR(500) NOT NULL,
			title VARCHAR(200) NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS click_events (
			id SERIAL PRIMARY KEY,
			ad_id INTEGER NOT NULL REFERENCES ads(id) ON DELETE CASCADE,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			ip_address VARCHAR(45),
			video_playback_time DECIMAL(10,2) DEFAULT 0,
			user_agent TEXT,
			processed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_ad_id ON click_events(ad_id)`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_timestamp ON click_events(timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_processed ON click_events(processed)`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_ad_timestamp ON click_events(ad_id, timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_click_events_processed_timestamp ON click_events(processed, timestamp)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	if err := insertSampleAds(db); err != nil {
		log.Printf("Warning: failed to insert sample ads: %v", err)
	}

	return nil
}

// Insert sample advertisement data
func insertSampleAds(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ads").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	sampleAds := []struct {
		imageURL, targetURL, title, description string
	}{
		{
			"https://example.com/ad1.jpg",
			"https://example.com/product1",
			"Product 1",
			"Product description",
		},
		{
			"https://example.com/ad2.jpg",
			"https://example.com/product2",
			"Product 2",
			"Product description",
		},
	}

	for _, ad := range sampleAds {
		_, err := db.Exec(
			"INSERT INTO ads (image_url, target_url, title, description) VALUES ($1, $2, $3, $4)",
			ad.imageURL, ad.targetURL, ad.title, ad.description,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

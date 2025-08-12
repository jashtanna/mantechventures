package models

import (
	"time"
)

// Ad represents a video advertisement
type Ad struct {
	ID          int       `json:"id" db:"id"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	TargetURL   string    `json:"target_url" db:"target_url"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ClickEvent represents a user click on an ad
type ClickEvent struct {
	ID                int       `json:"id" db:"id"`
	AdID              int       `json:"ad_id" db:"ad_id"`
	Timestamp         time.Time `json:"timestamp" db:"timestamp"`
	IPAddress         string    `json:"ip_address" db:"ip_address"`
	VideoPlaybackTime float64   `json:"video_playback_time" db:"video_playback_time"`
	UserAgent         string    `json:"user_agent" db:"user_agent"`
	Processed         bool      `json:"processed" db:"processed"`
}

// ClickRequest represents the incoming click data
type ClickRequest struct {
	AdID              int     `json:"ad_id" binding:"required"`
	VideoPlaybackTime float64 `json:"video_playback_time"`
	IPAddress         string  `json:"ip_address"`
	UserAgent         string  `json:"user_agent"`
}

// Analytics represents aggregated ad performance metrics
type Analytics struct {
	AdID            int       `json:"ad_id"`
	TotalClicks     int       `json:"total_clicks"`
	CTR             float64   `json:"ctr"` // Click-through rate
	AvgPlaybackTime float64   `json:"avg_playback_time"`
	TimeFrame       string    `json:"time_frame"`
	LastUpdated     time.Time `json:"last_updated"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

package model

type FeedingLog struct {
	ID              string  `json:"id"`
	ChildID         string  `json:"child_id"`
	FeedType        string  `json:"feed_type"`
	StartTime       string  `json:"start_time"`
	EndTime         *string `json:"end_time"`
	DurationMinutes *int    `json:"duration_minutes"`
	QuantityML      *int    `json:"quantity_ml,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

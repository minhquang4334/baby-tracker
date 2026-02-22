package model

type SleepLog struct {
	ID              string  `json:"id"`
	ChildID         string  `json:"child_id"`
	StartTime       string  `json:"start_time"`
	EndTime         *string `json:"end_time"`
	DurationMinutes *int    `json:"duration_minutes"`
	Notes           string  `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at"`
}

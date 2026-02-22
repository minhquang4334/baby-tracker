package model

type Child struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DateOfBirth string `json:"date_of_birth"`
	Gender      string `json:"gender"`
	PhotoURL    string `json:"photo_url,omitempty"`
	Notes       string `json:"notes,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

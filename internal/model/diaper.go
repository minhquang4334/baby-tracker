package model

type DiaperLog struct {
	ID         string `json:"id"`
	ChildID    string `json:"child_id"`
	DiaperType string `json:"diaper_type"`
	ChangedAt  string `json:"changed_at"`
	Notes      string `json:"notes,omitempty"`
	CreatedAt  string `json:"created_at"`
}

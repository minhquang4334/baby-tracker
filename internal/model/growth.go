package model

type GrowthLog struct {
	ID                    string `json:"id"`
	ChildID               string `json:"child_id"`
	MeasuredOn            string `json:"measured_on"`
	WeightGrams           *int   `json:"weight_grams,omitempty"`
	LengthMM              *int   `json:"length_mm,omitempty"`
	HeadCircumferenceMM   *int   `json:"head_circumference_mm,omitempty"`
	Notes                 string `json:"notes,omitempty"`
	CreatedAt             string `json:"created_at"`
}

package models

import "time"

type SecurityRule struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // low, medium, high, critical
	Category    string    `json:"category"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateRuleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Severity    string `json:"severity" validate:"required,oneof=low medium high critical"`
	Category    string `json:"category"`
	Enabled     bool   `json:"enabled"`
}

type UpdateRuleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Severity    *string `json:"severity" validate:"omitempty,oneof=low medium high critical"`
	Category    *string `json:"category"`
	Enabled     *bool   `json:"enabled"`
}

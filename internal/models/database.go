package models

import "time"

//SecretBase is a secret database model
type SecretBase struct {
	CreatedAt      time.Time  `json:"createdAt"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	RemainingViews int32      `json:"remainingViews"`
	SecretText     string     `json:"secretText"`
}

//Secret is a secret database model with persistence version tag
type Secret struct {
	SecretBase
	Version int64
}

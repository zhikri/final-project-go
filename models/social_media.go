package models

import "time"

type SocialMedia struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"type:varchar(50);not null"`
	SocialMediaURL string `gorm:"type:text;not null" json:"social_media_url"`
	UserID         uint   `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

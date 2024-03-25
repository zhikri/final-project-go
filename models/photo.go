package models

import "time"

type Photo struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"type:varchar(100);not null"`
	Caption   string `gorm:"type:varchar(200);not null"`
	PhotoURL  string `gorm:"type:text;not null" json:"photo_url"`
	UserID    uint   `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

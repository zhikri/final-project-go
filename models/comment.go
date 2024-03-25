package models

import "time"

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	PhotoID   uint   `gorm:"not null" json:"photo_id"`
	Message   string `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

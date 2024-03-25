package models

import "time"

type User struct {
	ID              uint   `gorm:"primaryKey"`
	Username        string `gorm:"type:varchar(50);not null"`
	Email           string `gorm:"type:varchar(150);not null"`
	Password        string `gorm:"type:text;not null"`
	Age             int    `gorm:"not null"`
	ProfileImageURL string `gorm:"type:text;not null" json:"profile_image_url"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	SocialMedias    []SocialMedia `gorm:"foreignKey:UserID;constaint:OnDelete:CASCADE"`
	Photos          []Photo       `gorm:"foreignKey:UserID;constaint:OnDelete:CASCADE"`
	Comments        []Comment     `gorm:"foreignKey:UserID;constaint:OnDelete:CASCADE"`
}

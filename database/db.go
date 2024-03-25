package database

import (
	"final-project-go/models"
	"gorm.io/gorm"
)
import "gorm.io/driver/postgres"

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=postgres dbname=final-project-go port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic("gagal koneksi ke database")
	}

	database.AutoMigrate(&models.User{}, &models.SocialMedia{}, &models.Photo{}, models.Comment{})

	DB = database

}

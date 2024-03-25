package photocontroller

import (
	"final-project-go/config"
	"final-project-go/database"
	"final-project-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
)

func GetAll(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Silakan login terlebih dahulu")
		return
	}
	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}
	// Memeriksa apakah token valid
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	var photos []models.Photo
	if err := database.DB.Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}

	var userIds []uint
	for _, photo := range photos {
		userIds = append(userIds, photo.UserID)
	}

	// Mengambil informasi pengguna terkait
	var users []models.User
	if err := database.DB.Where("id IN (?)", userIds).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Membuat map untuk mencocokkan ID pengguna dengan data pengguna
	userMap := make(map[uint]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Menggabungkan informasi pengguna dengan data foto
	var response []map[string]interface{}
	for _, photo := range photos {
		user, ok := userMap[photo.UserID]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user for photo"})
			return
		}

		response = append(response, map[string]interface{}{
			"id":        photo.ID,
			"title":     photo.Title,
			"caption":   photo.Caption,
			"photo_url": photo.PhotoURL,
			"user": map[string]interface{}{
				"id":       photo.UserID,
				"email":    user.Email,
				"username": user.Username,
			},
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetOne(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Silakan login terlebih dahulu")
		return
	}
	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}
	// Memeriksa apakah token valid
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var photo models.Photo

	// Mengambil foto berdasarkan ID
	if err := database.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Mengambil informasi pengguna terkait
	var user models.User
	if err := database.DB.First(&user, photo.UserID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	response := map[string]interface{}{
		"id":        photo.ID,
		"title":     photo.Title,
		"caption":   photo.Caption,
		"photo_url": photo.PhotoURL,
		"user": map[string]interface{}{
			"id":       photo.UserID,
			"email":    user.Email,
			"username": user.Username,
		},
	}
	c.JSON(http.StatusOK, response)
}

func CreatePhoto(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Silakan login terlebih dahulu")
		return
	}
	claims := &config.JWTClaim{}
	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}
	// Memeriksa apakah token valid
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	var photo models.Photo
	if err := c.BindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Validasi title
	if photo.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is required"})
		return
	}

	// Validasi url
	if photo.PhotoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo Url is Required"})
		return
	}

	photo.UserID = claims.ID
	database.DB.Create(&photo)
	c.JSON(http.StatusCreated, &photo)

}

func UpdatePhoto(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Silakan login terlebih dahulu")
		return
	}

	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	// Memeriksa apakah token valid
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.UserID = claims.ID
	if err := database.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	if database.DB.Save(&photo).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "update photo gagal"})
		return
	}

	c.JSON(http.StatusOK, &photo)

}

func DeletePhoto(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Silakan login terlebih dahulu")
		return
	}

	claims := &config.JWTClaim{}

	token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	if err != nil {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	// Memeriksa apakah token valid
	if !token.Valid {
		c.String(http.StatusUnauthorized, "Unauthorized! Token tidak valid")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid photo ID"})
		return
	}

	var photo models.Photo
	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.First(&photo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	if database.DB.Delete(&photo).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "hapus photo gagal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success delete"})
}

package socialmediacontroller

import (
	"final-project-go/config"
	"final-project-go/database"
	"final-project-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}

	userID := claims.(*config.JWTClaim).ID // Mendapatkan ID pengguna dari klaim JWT

	var socmeds []models.SocialMedia
	// Mengambil semua social media dari database
	if err := database.DB.Find(&socmeds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social media"})
		return
	}

	var userSocmeds []models.SocialMedia
	// Memfilter social media berdasarkan ID pengguna
	for _, sc := range socmeds {
		if sc.UserID == userID {
			userSocmeds = append(userSocmeds, sc)
		}
	}

	c.JSON(http.StatusOK, userSocmeds)
}

func GetOne(c *gin.Context) {
	id := c.Param("id")

	// Mencari social media berdasarkan ID dalam database
	var socmed models.SocialMedia
	if err := database.DB.First(&socmed, id).Error; err != nil {
		// Jika social media tidak ditemukan, kirim respons dengan status Not Found (404)
		c.JSON(http.StatusNotFound, gin.H{"error": "Sosial media tidak ditemukan"})
		return
	}

	// Mengembalikan social media dalam respons JSON
	c.JSON(http.StatusOK, socmed)
}

func CreateSocialMedia(c *gin.Context) {
	var socmed models.SocialMedia
	if err := c.ShouldBindJSON(&socmed); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		// Jika klaim tidak ada dalam konteks, tangani kesalahan sesuai kebutuhan Anda
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}

	//Validasi nama social media
	if socmed.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Social media name is required"})
		return
	}

	// Validasi url
	if socmed.SocialMediaURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Social Media Url is Required"})
		return
	}

	// Mengakses nilai spesifik dalam klaim
	userID := claims.(*config.JWTClaim).ID
	socmed.UserID = userID
	if err := database.DB.Create(&socmed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social media"})
		return
	}

	c.JSON(http.StatusCreated, socmed)
}

func UpdateSocialMedia(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}
	// Mendapatkan ID Sosmed dari parameter URL
	id := c.Param("id")

	// Mencari social media yang akan diupdate dalam database
	var socmed models.SocialMedia
	if err := database.DB.First(&socmed, id).Error; err != nil {
		// Jika social media tidak ditemukan, kirim respons dengan status Not Found (404)
		c.JSON(http.StatusNotFound, gin.H{"error": "social media tidak ditemukan"})
		return
	}

	// Bind JSON dari body permintaan ke struct Comment
	if err := c.ShouldBindJSON(&socmed); err != nil {
		// Jika terjadi kesalahan saat binding JSON, kirim respons dengan status Bad Request (400)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := claims.(*config.JWTClaim).ID
	socmed.UserID = userID
	// Update social media dalam database
	if err := database.DB.Save(&socmed).Error; err != nil {
		// Jika terjadi kesalahan saat menyimpan perubahan ke database, kirim respons dengan status Internal Server Error (500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui social media"})
		return
	}

	// Mengembalikan social media yang telah diupdate dalam respons JSON
	c.JSON(http.StatusOK, socmed)
}

func DeleteSocialMedia(c *gin.Context) {
	id := c.Param("id")

	// Mencari social media yang akan dihapus dalam database
	var socmed models.SocialMedia
	if err := database.DB.First(&socmed, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media tidak ditemukan"})
		return
	}

	// Menghapus social media dari database
	if err := database.DB.Delete(&socmed).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media berhasil dihapus"})
}

package commentcontroller

import (
	"final-project-go/config"
	"final-project-go/database"
	"final-project-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAll(c *gin.Context) {
	var comments []models.Comment
	if err := database.DB.Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}

	var userIds []uint
	for _, cmt := range comments {
		userIds = append(userIds, cmt.UserID)
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
	for _, cmt := range comments {
		user, ok := userMap[cmt.UserID]
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user for comment"})
			return
		}

		response = append(response, map[string]interface{}{
			"id":       cmt.ID,
			"message":  cmt.Message,
			"user_id":  cmt.UserID,
			"photo_id": cmt.PhotoID,
			"user": map[string]interface{}{
				"id":       cmt.UserID,
				"email":    user.Email,
				"username": user.Username,
			},
		})
	}

	// Mengembalikan semua komentar tanpa memfilter
	c.JSON(http.StatusOK, comments)
}

func GetOne(c *gin.Context) {
	// Mendapatkan ID komentar dari parameter URL
	commentID := c.Param("id")

	// Mencari komentar berdasarkan ID dalam database
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		// Jika komentar tidak ditemukan, kirim respons dengan status Not Found (404)
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Mengembalikan komentar dalam respons JSON
	c.JSON(http.StatusOK, comment)
}

func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, exists := c.Get("claims")
	if !exists {
		// Jika klaim tidak ada dalam konteks, tangani kesalahan sesuai kebutuhan Anda
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}

	// Mengakses nilai spesifik dalam klaim
	userID := claims.(*config.JWTClaim).ID
	comment.UserID = userID
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func UpdateComment(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}
	// Mendapatkan ID komentar dari parameter URL
	commentID := c.Param("id")

	// Mencari komentar yang akan diupdate dalam database
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		// Jika komentar tidak ditemukan, kirim respons dengan status Not Found (404)
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Bind JSON dari body permintaan ke struct Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		// Jika terjadi kesalahan saat binding JSON, kirim respons dengan status Bad Request (400)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := claims.(*config.JWTClaim).ID
	comment.UserID = userID
	// Update komentar dalam database
	if err := database.DB.Save(&comment).Error; err != nil {
		// Jika terjadi kesalahan saat menyimpan perubahan ke database, kirim respons dengan status Internal Server Error (500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui komentar"})
		return
	}

	// Mengembalikan komentar yang telah diupdate dalam respons JSON
	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	// Mencari komentar yang akan dihapus dalam database
	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		// Jika komentar tidak ditemukan, kirim respons dengan status Not Found (404)
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Menghapus komentar dari database
	if err := database.DB.Delete(&comment).Error; err != nil {
		// Jika terjadi kesalahan saat menghapus komentar dari database, kirim respons dengan status Internal Server Error (500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar"})
		return
	}

	// Mengembalikan respons berhasil dengan status OK (200)
	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}

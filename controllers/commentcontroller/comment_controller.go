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

	// Get user
	var users []models.User
	if err := database.DB.Where("id IN (?)", userIds).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Membuat map untuk mencocokkan ID user dengan data user
	userMap := make(map[uint]models.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// Menggabungkan data user dengan data foto
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
	c.JSON(http.StatusOK, comments)
}

func GetOne(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Klaim tidak tersedia"})
		return
	}

	if comment.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment Message is required"})
		return
	}

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
	commentID := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := claims.(*config.JWTClaim).ID
	comment.UserID = userID
	// Update komentar dalam database
	if err := database.DB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui komentar"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if err := database.DB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Komentar tidak ditemukan"})
		return
	}

	// Menghapus komentar dari database
	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komentar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Komentar berhasil dihapus"})
}

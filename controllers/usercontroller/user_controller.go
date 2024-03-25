package usercontroller

import (
	"final-project-go/config"
	"final-project-go/database"
	"final-project-go/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func Register(c *gin.Context) {
	var userInput models.User
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
	userInput.Password = string(hashPassword)
	database.DB.Create(&userInput)
	c.JSON(http.StatusCreated, &userInput)
}

func Login(c *gin.Context) {
	var userInput models.User
	if err := c.BindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", userInput.Email).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Username atau password salah"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	//validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username atau password salah"})
		return
	}

	//JWT Implementation
	expTime := time.Now().Add(time.Minute * 1)
	claims := &config.JWTClaim{
		Username: user.Username,
		ID:       user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "jwt-v5",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	//Algoritma
	tokenAlg := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//signed token
	token, err := tokenAlg.SignedString(config.JWT_KEY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	//set token ke cookie
	c.SetCookie("token", token, 60, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Update(c *gin.Context) {
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

	var userInput models.User
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if database.DB.Model(&userInput).Where("id = ?", claims.ID).Updates(&userInput).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "update user gagal"})
	}

	c.JSON(http.StatusOK, &userInput)
}

func Delete(c *gin.Context) {
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

	var userInput models.User
	if err := database.DB.First(&userInput, claims.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&userInput).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil menghapus user"})
}

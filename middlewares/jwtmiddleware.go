package middlewares

import (
	"final-project-go/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! Silakan login terlebih dahulu"})
			return
		}

		claims := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(cookie, claims, func(t *jwt.Token) (interface{}, error) {
			return config.JWT_KEY, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized! Token tidak valid"})
			return
		}

		// Menetapkan klaim JWT ke konteks Gin untuk digunakan di handler selanjutnya
		c.Set("claims", claims)

		// Lanjutkan ke handler berikutnya
		c.Next()
	}
}

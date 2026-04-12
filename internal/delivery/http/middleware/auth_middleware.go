package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "akses ditolak!, token tidak ditemukan"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "format token salah!"})
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("metode penandatanganan tidak valid")
			}

			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token tidak valid atau kadaluarsa!"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user.id", claims["user.id"])
			c.Set("store.id", claims["store.id"])
			c.Set("role", claims["role"])
		} else {	
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "gagal membaca data dari dalam token"})
			return 
		}
		
		c.Next()
	}
}

func OwnerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		
		if !exists || role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Akses ditolak: Fitur ini hanya untuk Owner toko",
			})
			c.Abort() 
			return
		}
		
		c.Next()
	}
}
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fathima-sithara/UserManagment/initalizeres"
	"github.com/fathima-sithara/UserManagment/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization is required"})
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("not expected signing method %v", token.Method)
		}
		return []byte("SECRET"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invaid token"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invaid token"})
		c.Abort()
		return
	}

	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		c.Abort()
		return
	}

	var user models.User
	if err := initalizeres.DB.Find(&user, "id=?", claims["exp"]).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		c.Abort()
		return
	}
	c.Set("user", user)

	c.JSON(http.StatusOK, gin.H{})
}

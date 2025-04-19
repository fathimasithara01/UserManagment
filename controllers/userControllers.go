package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/fathima-sithara/UserManagment/initalizeres"
	"github.com/fathima-sithara/UserManagment/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var validate *validator.Validate

	validate = validator.New()

	var input models.UserInput
	if c.ShouldBindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid input"})
		return
	}

	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	var existUSer models.User
	err := initalizeres.DB.Where("username=?", input.Username).Or("email=?", input.Email).First(&existUSer).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email already exiist"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to hash password"})
		return
	}
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hash),
	}
	if err := initalizeres.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to creatr user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	var input models.UserInput

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	err := initalizeres.DB.Where("username=?", input.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falied to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": tokenString,
	})

}
func LogOut(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "", "localhost", false, true)
	c.JSON(http.StatusAccepted, gin.H{"msg": "successful"})
}

func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}

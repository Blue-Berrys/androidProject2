package Bcrypt

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func EncryptionByPassword(password string) string {
	password2 := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password2, bcrypt.DefaultCost)
	return string(hashedPassword)
}

func EncryptionMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		password := c.Query("password") //Get方式得到
		if password == "" {
			password = c.PostForm("password") //Post方式得到
		}
		c.Set("password", EncryptionByPassword(password))
		c.Next()
	}
}

func QueryEqualEncryptAndPassword(encryptPassword, password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptPassword), []byte(password)); err != nil {
		return false
	}
	return true
}

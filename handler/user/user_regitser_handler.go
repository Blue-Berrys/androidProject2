package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterResponse struct {
	CommonResponse
	*service.RegisterResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	registerResponse, err := service.NewQueryUserRegister(username, password)
	fmt.Println(registerResponse)
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	fmt.Println("9999")
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		RegisterResponse: registerResponse,
	})
}

package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
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
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		CommonResponse:   CommonResponse{StatusCode: 0, StatusMsg: config.SUCCESS_MSG},
		RegisterResponse: registerResponse,
	})
}

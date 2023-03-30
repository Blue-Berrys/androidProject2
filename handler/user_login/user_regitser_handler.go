package handler

import (
	service "androidProject2/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRegisterResponse struct {
	CommonResponse
	*service.RegisterResponse
}

func UserRegisterHandler(c *gin.Context) {
	username := c.Query("username")
	rawVal, _ := c.Get("password")
	password, ok := rawVal.(string)
	if !ok {
		c.JSON(http.StatusOK, UserRegisterResponse{
			CommonResponse: CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析出错",
			},
		})
		return
	}
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
		CommonResponse:   CommonResponse{StatusCode: 0},
		RegisterResponse: registerResponse,
	})
}

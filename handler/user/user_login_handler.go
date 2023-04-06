package handler

import (
	"androidProject2/config"
	"androidProject2/model"
	service "androidProject2/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	model.CommonResponse
	*service.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	userLoginResponse, err := service.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: model.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	fmt.Println(userLoginResponse)
	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: model.CommonResponse{StatusCode: 0, StatusMsg: config.SUCCESS_MSG},
		LoginResponse:  userLoginResponse,
	})
}

package handler

import (
	"androidProject2/config"
	service "androidProject2/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	CommonResponse
	*service.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	rawPassword, _ := c.Get("password")
	password, ok := rawPassword.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
	}
	userLoginResponse, err := service.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: CommonResponse{
			StatusCode: 0,
			StatusMsg:  config.SUCCESS_MSG,
		},
		LoginResponse: userLoginResponse,
	})
}

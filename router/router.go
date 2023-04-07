package router

import (
	handler2 "androidProject2/handler/FriendsChat"
	handler3 "androidProject2/handler/like"
	handler "androidProject2/handler/user"
	"androidProject2/middleware/JWT"
	model "androidProject2/model/db"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	// public directory is used to serve static resources
	model.Init()
	r.Static("/static", "./static")

	apiRouter := r.Group("/androidProject2")

	// basic apis
	//apiRouter.GET("/user/register/", Bcrypt.EncryptionMiddleWare(), handler.UserRegisterHandler)
	apiRouter.POST("/user/register/", handler.UserRegisterHandler)
	apiRouter.POST("/user/login/", handler.UserLoginHandler)
	apiRouter.POST("/user/info/", JWT.JWTMiddleware(), handler.UserInfoHandler)
	apiRouter.POST("/publish/action/", JWT.JWTMiddleware(), handler2.PublishActionHandler)
	apiRouter.POST("/publish/list/", JWT.JWTMiddleware(), handler2.PublishListHandler)

	// extra apis - I
	apiRouter.POST("/favorite/action/", JWT.JWTMiddleware(), handler3.LikeActionHandler)

	apiRouter.POST("/comment/action/", JWT.JWTMiddleware())
	apiRouter.GET("/comment/list/", JWT.JWTMiddleware())

	return r
}

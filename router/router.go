package router

import (
	handler4 "androidProject2/handler/Comment"
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
	apiRouter.POST("/favorite/list/", JWT.JWTMiddleware(), handler3.LikeListHandler)
	apiRouter.POST("/comment/action/", JWT.JWTMiddleware(), handler4.CommentActionHandler)
	apiRouter.POST("/comment/list/", JWT.JWTMiddleware(), handler4.CommentListHandler)

	return r
}

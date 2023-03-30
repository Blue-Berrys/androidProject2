package router

import (
	"androidProject2/middleware/JWT"
	"github.com/gin-gonic/gin"
	"github.com/neverTanking/TiktokByGo/db"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	// public directory is used to serve static resources
	db.Init()
	r.Static("/static", "./static")

	apiRouter := r.Group("/androidProject2")

	// basic apis
	apiRouter.POST("/user/register/")
	apiRouter.POST("/user/login/")
	apiRouter.GET("/user/info/")
	apiRouter.POST("/publish/action/")
	apiRouter.GET("/publish/list/")

	// extra apis - I
	apiRouter.POST("/favorite/action/", JWT.JWTMiddleware())

	apiRouter.POST("/comment/action/", JWT.JWTMiddleware())
	apiRouter.GET("/comment/list/", JWT.JWTMiddleware())

	return r
}

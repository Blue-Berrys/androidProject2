package main

import (
	"androidProject2/router"
	"github.com/gin-gonic/gin"
)

var R *gin.Engine

func main() {

	R = gin.Default()

	router.InitRouter(R)

	err := R.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

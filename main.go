package main

import (
	"androidProject2/router"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	router.InitRouter(r)

	err := r.Run(":8080")
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

package main

import (
	"api_gateway/internal/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	r = routers.SetStaticRouters(r)
	r = routers.SetApiRouters(r)

	err := r.Run(":8000")
	if err != nil {
		log.Fatal(err)
	}
}

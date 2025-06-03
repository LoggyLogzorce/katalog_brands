package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"user_service/internal/db"
	"user_service/internal/routers"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := routers.SetRouters()

	db.Connect()

	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

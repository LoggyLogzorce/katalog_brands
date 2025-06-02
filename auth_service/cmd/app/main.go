package main

import (
	"auth_service/internal/db"
	"auth_service/internal/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := routers.SetRouters()

	db.Connect()

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}

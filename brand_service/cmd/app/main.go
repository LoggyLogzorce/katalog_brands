package main

import (
	"brand_service/internal/db"
	"brand_service/internal/es"
	"brand_service/internal/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := routers.SetRouters()

	db.Connect()
	es.BootstrapIndexing()

	if err := r.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/db"
	"product_service/internal/es"
	"product_service/internal/routers"
)

func main() {
	gin.SetMode(gin.DebugMode)
	r := routers.SetRouters()

	db.Connect()
	es.BootstrapIndexing()

	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}

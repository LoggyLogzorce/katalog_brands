package main

import (
	"brand_service/internal/api"
	"brand_service/internal/db"
	"brand_service/internal/es"
	"brand_service/internal/routers"
	"brand_service/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)

	dbConn := db.Connect()
	es.BootstrapIndexing()

	esClient, _ := es.NewClient()
	esRepo := es.NewIndexer(esClient)
	esService := api.NewServiceEs(esRepo)

	brandRepo := storage.NewBrandRepository(dbConn)
	brandHandler := api.NewBrandHandler(brandRepo, esService)

	r := routers.SetRouters(brandHandler)

	if err := r.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}

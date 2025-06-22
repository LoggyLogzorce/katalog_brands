package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"product_service/internal/api"
	"product_service/internal/db"
	"product_service/internal/es"
	"product_service/internal/routers"
	"product_service/internal/storage"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	dbConn := db.Connect()
	esClient, _ := es.NewClient()

	catRepo := storage.NewCategoryRepository(dbConn)
	esRepo := es.NewIndexer(esClient)
	catHandler := api.NewCategoryHandler(catRepo, esRepo)

	prodRepo := storage.NewProductRepository(dbConn)
	prodHandler := api.NewProductHandler(prodRepo, esRepo)

	r := routers.SetRouters(catHandler, prodHandler)

	es.BootstrapIndexing()

	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}

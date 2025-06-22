package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"user_service/internal/api"
	"user_service/internal/db"
	"user_service/internal/routers"
	"user_service/internal/storage"
)

func main() {
	gin.SetMode(gin.DebugMode)

	dbConn := db.Connect()

	revRepo := storage.NewUserRepository(dbConn)
	revHandler := api.NewReviewHandler(revRepo)

	r := routers.SetRouters(revHandler)

	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}

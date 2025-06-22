package main

import (
	"auth_service/internal/api"
	"auth_service/internal/db"
	"auth_service/internal/routers"
	"auth_service/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.DebugMode)

	dbConn := db.Connect()

	authRepo := storage.NewUserRepository(dbConn)
	authHandler := api.NewAuthHandler(authRepo)

	r := routers.SetRouters(authHandler)

	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}

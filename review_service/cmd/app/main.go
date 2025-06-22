package main

import (
	"log"
	"review_service/internal/api"
	"review_service/internal/db"
	"review_service/internal/routers"
	"review_service/internal/storage"
)

func main() {
	dbConn := db.Connect()

	revRepo := storage.NewReviewRepository(dbConn)
	revHandler := api.NewReviewHandler(revRepo)

	r := routers.SetRouters(revHandler)

	if err := r.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}

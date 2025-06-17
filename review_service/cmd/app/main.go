package main

import (
	"log"
	"review_service/internal/db"
	"review_service/internal/routers"
)

func main() {
	r := routers.SetRouters()

	db.Connect()

	if err := r.Run(":8085"); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/yuanzix/beatify-core/handlers"
	"github.com/yuanzix/beatify-core/utils"
)

func main() {
	store, err1 := utils.NewPostgresStore()
	_, err2 := utils.ReadJWTSecret()
	_, _, err3 := utils.ReadGmailDetails()
	_, err4 := utils.ReadBackendURL()
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		log.Fatal(err1, err2, err3, err4)
	}

	server := handlers.NewAPIServer(":5000", store)
	server.Run()
}

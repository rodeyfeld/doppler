package application

import (
	"doppler/internal/db"
	"doppler/internal/server"
	"log"
)

func Start() {
	db.SetupDb()
	appServer := server.New()
	err := appServer.Start()
	if err != nil {
		log.Print(err)
	}
}

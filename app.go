package application

import (
	"doppler/internal/server"
	"log"
)

func Start() {
	log.Print("spin")
	appServer := server.NewServer()
	err := appServer.Start()
	if err != nil {
		log.Print(err)
	}
}

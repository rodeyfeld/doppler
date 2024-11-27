package application

import (
	"doppler/internal/db"
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"log"
)

func Start() {
	db.SetupDb()
	dopplerServer := server.NewDopplerServer()
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

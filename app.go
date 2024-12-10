package application

import (
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"log"
)

func Start() {
	dopplerServer := server.NewDopplerServer()
	dopplerServer.Echo.Static("/static", "static")
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

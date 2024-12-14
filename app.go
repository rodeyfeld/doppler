package application

import (
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"log"
)

func Start() {
	dopplerServer := server.NewDopplerServer()
	dopplerServer.Echo.Use(session.Middleware(sessions.NewCookieStore([]byte("secret!"))))
	dopplerServer.Echo.Static("/static", "static")
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

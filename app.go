package application

import (
	"doppler/internal/db"
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"log"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

func Start() {
	dopplerServer := server.NewDopplerServer()

	// Run database migrations on startup
	if err := db.RunMigrations(dopplerServer.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	dopplerServer.Echo.Use(session.Middleware(sessions.NewCookieStore([]byte("secret!"))))
	dopplerServer.Echo.Static("/static", "static")
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

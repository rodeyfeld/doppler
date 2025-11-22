package application

import (
	"doppler/internal/db"
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

func Start() {
	dopplerServer := server.NewDopplerServer()

	// Run database migrations on startup
	if err := db.RunMigrations(dopplerServer.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Session secret for cookie encryption/signing
	sessionSecret := os.Getenv("DOPPLER_SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("DOPPLER_SESSION_SECRET environment variable is required")
	}
	dopplerServer.Echo.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))
	dopplerServer.Echo.Static("/static", "static")
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

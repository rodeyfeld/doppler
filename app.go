package application

import (
	"doppler/internal/db"
	"doppler/internal/server"
	"doppler/internal/server/routes"
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Start() {
	dopplerServer := server.NewDopplerServer()

	// TODO: Move this to a callable depends on job
	if err := db.RunMigrations(dopplerServer.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Session secret for cookie encryption/signing
	sessionSecret := os.Getenv("DOPPLER_SESSION_SECRET")
	if sessionSecret == "" {
		log.Fatal("DOPPLER_SESSION_SECRET environment variable is required")
	}
	dopplerServer.Echo.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionSecret))))

	// Static files carry a Last-Modified header, which browsers would otherwise use to
	// cache them heuristically and keep a stale stylesheet after a rebuild. no-cache
	// makes them revalidate (a cheap 304 when nothing changed).
	static := dopplerServer.Echo.Group("/static", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache")
			return next(c)
		}
	})
	static.Static("/", "static")
	routes.Setup(dopplerServer)
	err := dopplerServer.Start()
	if err != nil {
		log.Print(err)
	}
}

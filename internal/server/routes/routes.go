package routes

import (
	"doppler/internal/server"
	"doppler/internal/server/handlers"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(s *server.DopplerServer) {

	homeHandler := handlers.NewHomeHandler(s)
	s.Echo.GET("/", homeHandler.HomeIndex)

	postHandler := handlers.NewPostHandler(s)

	s.Echo.Use(middleware.Logger())
	g := s.Echo.Group("/app")

	g.GET("/", postHandler.Index)
	g.POST("/create", postHandler.Create)
}

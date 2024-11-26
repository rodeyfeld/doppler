package routes

import (
	"doppler/internal/server"
	"doppler/internal/server/handlers"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(s *server.AppServer) {
	postHandler := handlers.NewPostHandler(s)

	s.Echo.Use(middleware.Logger())
	g := s.Echo.Group("posts")
	g.GET("/", postHandler.GetPosts)
}

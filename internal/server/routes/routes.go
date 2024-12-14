package routes

import (
	"doppler/internal/server"
	"doppler/internal/server/handlers"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(s *server.DopplerServer) {

	s.Echo.Use(middleware.Logger())

	homeHandler := handlers.NewHomeHandler(s)
	s.Echo.GET("/", homeHandler.HomeIndex)

	authHandler := handlers.NewAuthHandler(s)
	g := s.Echo.Group("/auth")
	g.GET("/", authHandler.Index)
	g.POST("/login", authHandler.Login)
	//	g.POST("/logout", authHandler.Logout)
	//	g.POST("/register", authHandler.Register)

	postHandler := handlers.NewPostHandler(s)
	g = s.Echo.Group("/doppler")
	g.GET("/", postHandler.Index)
	g.POST("/create", postHandler.Create)
}

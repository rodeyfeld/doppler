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
	postHandler := handlers.NewPostHandler(s)
	g := s.Echo.Group("/doppler")
	g.GET("/", postHandler.Index)
	g.POST("/create", postHandler.Create)
	g.GET("/user-info/:id", postHandler.UserInfo)
	g.GET("/images/:filename", postHandler.GetImage)

	g.GET("/login", authHandler.LoginIndex)
	g.POST("/login", authHandler.Login)
	g.GET("/logout", authHandler.Logout)
	g.GET("/profile", authHandler.ProfileIndex)

	g.GET("/signup", authHandler.SignupIndex)
	g.POST("/signup", authHandler.Signup)
}

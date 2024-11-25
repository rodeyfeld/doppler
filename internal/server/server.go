package server

import (
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
}

func NewServer() *Server {
	return &Server{
		Echo: echo.New(),
	}
}

func (server *Server) Start() error {
	return server.Echo.Start(":1323")
}


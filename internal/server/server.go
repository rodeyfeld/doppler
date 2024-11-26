package server

import (
	"github.com/labstack/echo/v4"
)

type AppServer struct {
	Echo *echo.Echo
}

func New() *AppServer {
	return &AppServer{
		Echo: echo.New(),
	}
}

func (server *AppServer) Start() error {
	return server.Echo.Start(":1323")
}

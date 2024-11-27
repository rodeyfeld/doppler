package server

import (
	"database/sql"
	"doppler/internal/db"
	"github.com/labstack/echo/v4"
)

type DopplerServer struct {
	Echo *echo.Echo
	DB   *sql.DB
}

func NewDopplerServer() *DopplerServer {
	return &DopplerServer{
		Echo: echo.New(),
		DB:   db.Connect(),
	}
}

func (server *DopplerServer) Start() error {
	return server.Echo.Start(":1323")
}

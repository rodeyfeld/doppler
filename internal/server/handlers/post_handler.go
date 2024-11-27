package handlers

import (
	"doppler/internal/server"
	//"doppler/internal/services"

	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	server *server.DopplerServer
}

func NewPostHandler(s *server.DopplerServer) *PostHandler {
	return &PostHandler{server: s}
}

func (pH *PostHandler) GetPosts(c echo.Context) error {
	//posts := services.GetPosts(pH.server.DB)
	return nil
}

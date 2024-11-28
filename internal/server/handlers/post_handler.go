package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	server *server.DopplerServer
}

func NewPostHandler(s *server.DopplerServer) *PostHandler {
	return &PostHandler{server: s}
}

func (h *PostHandler) Index(c echo.Context) error {
	cmp := components.PostIndex(components.CreatePost(), components.ListPost())
	return renderView(c, cmp)
}

package handlers

import (
	"doppler/internal/server"
)

type PostHandler struct {
	server *server.AppServer
}

func NewPostHandler(s *server.AppServer) *PostHandler {
	return &PostHandler{server: s}
}

package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"
	"doppler/internal/services"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type PostHandler struct {
	server *server.DopplerServer
}

func NewPostHandler(s *server.DopplerServer) *PostHandler {
	return &PostHandler{server: s}
}

func (h *PostHandler) Index(c echo.Context) error {
	posts := services.GetPosts(h.server.DB)
	cmp := components.PostIndex(components.ListPosts(posts))
	return renderView(c, cmp)
}

func (h *PostHandler) Create(c echo.Context) error {
	sess, err := session.Get("auth-session", c)
	if err != nil {
		return err
	}
	userID := sess.Values["userID"].(int)
	post := services.CreatePost(h.server.DB, userID, c.FormValue("title"), c.FormValue("content"))
	cmp := components.PostSuccess(post)
	return renderView(c, cmp)
}

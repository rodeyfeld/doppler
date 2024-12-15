package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"
	"doppler/internal/services"
	"log"
	"strconv"

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

func (h *PostHandler) UserInfo(c echo.Context) error {
	p := c.Param("id")
	log.Printf(p)
	id, err := strconv.Atoi(p)
	if err != nil {

		return err
	}
	user, err := services.GetUserByID(h.server.DB, id)
	if err != nil {
		return err
	}
	cmp := components.PostUserInfo(user)
	return renderView(c, cmp)
}

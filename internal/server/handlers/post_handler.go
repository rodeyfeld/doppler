package handlers

import (
	"doppler/internal/components/post"
	"doppler/internal/components/shared"
	"doppler/internal/models"
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
	sess, err := session.Get("auth-session", c)
	if err != nil {
		log.Printf("Failed to get auth session")
		return err
	}

	var user *models.User
	// Check if user is logged in
	if userID, ok := sess.Values["userID"].(int); ok {
		user, err = services.GetUserByID(h.server.DB, userID)
		if err != nil {
			log.Printf("Failed to get user by ID: %v", err)
		}
	}

	cmp := post.PostIndex(post.ListPosts(posts), user)
	return renderView(c, cmp)
}

func (h *PostHandler) Create(c echo.Context) error {
	sess, err := session.Get("auth-session", c)
	if err != nil {
		log.Printf("Failed to get auth session: %v", err)
		return err
	}

	// Check if user is logged in
	userID, ok := sess.Values["userID"].(int)
	if !ok {
		log.Printf("User not authenticated, cannot create post")
		cmp := shared.AuthRequired()
		return renderView(c, cmp)
	}

	createdPost := services.CreatePost(h.server.DB, userID, c.FormValue("title"), c.FormValue("content"))
	cmp := post.PostSuccess(createdPost)
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
	cmp := post.PostUserInfo(user)
	return renderView(c, cmp)
}

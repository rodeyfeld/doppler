package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"
	"doppler/internal/services"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	server *server.DopplerServer
}

func NewAuthHandler(s *server.DopplerServer) *AuthHandler {
	return &AuthHandler{server: s}
}

func (h *AuthHandler) Index(c echo.Context) error {
	cmp := components.AuthIndex("message")
	return renderView(c, cmp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	if services.ValidateUser(h.server.DB, c.FormValue("username"), c.FormValue("password")) {
		return c.Redirect(http.StatusSeeOther, "/app/")
	}
	sess, err := session.Get("auth-session-key", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600,
	}
	return c.Redirect(http.StatusSeeOther, "/app/")
}

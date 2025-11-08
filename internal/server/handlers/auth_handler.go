package handlers

import (
	"doppler/internal/components/auth"
	"doppler/internal/server"
	"doppler/internal/services"
	"log"
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

func (h *AuthHandler) LoginIndex(c echo.Context) error {
	cmp := auth.LoginIndex()
	return renderView(c, cmp)
}

func (h *AuthHandler) Login(c echo.Context) error {
	if !services.ValidateUser(h.server.DB, c.FormValue("username"), c.FormValue("password")) {
		return c.Redirect(http.StatusFound, "/doppler/signup/")
	}
	user, err := services.GetUserByUsername(h.server.DB, c.FormValue("username"))
	if err != nil {
		log.Printf("Failed to get username after authenticating user")
		return err
	}
	sess, err := session.Get("auth-session", c)
	if err != nil {
		log.Printf("Failed to setup gorilla session")
		return err
	}
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600,
	}
	sess.Values["username"] = user.Username
	sess.Values["authed"] = true
	sess.Values["userID"] = user.ID
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusFound, "/doppler/")
}

func (h *AuthHandler) ProfileIndex(c echo.Context) error {
	sess, err := session.Get("auth-session", c)
	if err != nil {
		return c.Redirect(http.StatusInternalServerError, "/doppler/signup/")
	}
	username := sess.Values["username"].(string)
	user, err := services.GetUserByUsername(h.server.DB, username)
	if err != nil {
		return c.Redirect(http.StatusInternalServerError, "/doppler/signup/")
	}
	cmp := auth.ProfileIndex(user)
	return renderView(c, cmp)
}

func (h *AuthHandler) SignupIndex(c echo.Context) error {
	cmp := auth.SignupIndex()
	return renderView(c, cmp)
}

func (h *AuthHandler) Signup(c echo.Context) error {

	user, _ := services.GetUserByUsername(h.server.DB, c.FormValue("username"))
	if user != nil {
		return c.Redirect(http.StatusInternalServerError, "/doppler/signup/")
	}

	user = services.CreateUser(h.server.DB, c.FormValue("username"), c.FormValue("password"), c.FormValue("email"))
	return c.Redirect(http.StatusFound, "/doppler/")
}

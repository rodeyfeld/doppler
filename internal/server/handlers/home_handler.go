package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct {
	server *server.DopplerServer
}

func NewHomeHandler(s *server.DopplerServer) *HomeHandler {
	return &HomeHandler{server: s}
}

func (h *HomeHandler) HomeIndex(c echo.Context) error {
	homeComponent := components.HomeIndex("Doppler")
	return renderView(c, homeComponent)
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

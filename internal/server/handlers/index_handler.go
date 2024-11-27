package handlers

import (
	"doppler/internal/components"
	"doppler/internal/server"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type IndexHandler struct {
	server *server.DopplerServer
}

func NewIndexHandler(s *server.DopplerServer) *IndexHandler {
	return &IndexHandler{server: s}
}

func (iH *IndexHandler) Index(c echo.Context) error {
	indexComponent := components.Layout("Doppler")
	return renderView(c, indexComponent)
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

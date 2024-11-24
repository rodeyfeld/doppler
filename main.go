package main

import (
	"github.com/labstack/echo/v4"
	"github.com/rodeyfeld/doppler/db"
	"github.com/rodeyfeld/doppler/handlers"
	"net/http"
)

const dbName = "user_data.db"

func main() {
	app := echo.New()
	app.HTTPErrorHandler = handlers.DopplerHTTPErrorHandler
	app.Static("/", "assets")

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HELO")
	})

	uStore, err := db.NewUserStore(dbName)
	app.Logger.Fatal(e.Start(":1323"))
}

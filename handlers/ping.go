package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Ping(c echo.Context) error {
	return echo.NewHTTPError(http.StatusOK, "Ping!")
}

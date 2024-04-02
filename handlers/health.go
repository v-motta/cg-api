package handlers

import (
	"cost-guardian-api/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
// @Tags Health
// @Summary Check database connection
// @Success 200
// @Router /health [get]
func Health(c echo.Context) error {
	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to ping database")
	}

	return echo.NewHTTPError(http.StatusOK, "Database is connected successfully!")
}

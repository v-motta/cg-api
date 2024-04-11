package handlers

import (
	"cost-guardian-api/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Health godoc
// @Tags Health
// @Summary Check database connection
// @Description Checks if the database connection is healthy by pinging the database. If the database is connected successfully, it returns a 200 status code with a success message. If the database connection fails, it returns a 500 status code with an error message.
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string "Database is connected successfully!"
// @Failure 500 {object} map[string]string "Failed to connect to database" or "Failed to ping database"
// @Router /health [get]
func Health(c echo.Context) error {

	SendEmail()

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

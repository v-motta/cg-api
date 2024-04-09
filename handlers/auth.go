package handlers

import (
	"cost-guardian-api/db"
	"cost-guardian-api/models"
	"database/sql"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary User login
// @Description Authenticate user and return JWT token if successful
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param   AuthUser     body    models.AuthUser     true        "User credentials"
// @Success 200 {string} string	"Returns the JWT token"
// @Failure 400 {object} string "Invalid request, user credentials are not provided"
// @Failure 500 {object} string "Internal server error, failed to connect to database"
// @Router /login [post]
func Login(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return echo.NewHTTPError(http.StatusMethodNotAllowed, "Method not allowed")
	}

	var user models.AuthUser
	if err := c.Bind(&user); err != nil {
		return err
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")

	}
	defer db.Close()

	var storedUser models.AuthUser
	err = db.QueryRow("SELECT username, password FROM users WHERE username=$1", user.Username).Scan(&storedUser.Username, &storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
		}
		return err
	}

	if user.Password != storedUser.Password {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = storedUser.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t})
}

package handlers

import (
	"cost-guardian-api/db"
	"cost-guardian-api/helpers"
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
// @Param   LoginUser     body    models.LoginUser     true        "User credentials"
// @Success 200 {string} string	"Returns the JWT token"
// @Failure 400 {object} string "Invalid request, user credentials are not provided"
// @Failure 500 {object} string "Internal server error, failed to connect to database"
// @Router /login [post]
func Login(c echo.Context) error {
	var user models.LoginUser
	if err := c.Bind(&user); err != nil {
		return err
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	var storedUser models.User
	err = db.QueryRow("SELECT username, password, role FROM users WHERE username=$1", user.Username).Scan(&storedUser.Username, &storedUser.Password, &storedUser.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
		}
		return err
	}

	check := helpers.CheckPasswordHash(user.Password, storedUser.Password)
	if !check {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = storedUser.Username
	claims["role"] = storedUser.Role
	if user.Remember {
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	}

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": t, "username": storedUser.Username})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.User true "User to be created"
// @Security Bearer
// @Success 200 {object} models.User "Successfully created user"
// @Failure 400 {object} string "Failed to decode request body"
// @Failure 500 {object} string "Failed to connect to database"
// @Failure 500 {object} string "Failed to insert user into database"
// @Router /signup [post]
func Signup(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body")
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	password, err := helpers.HashPassword(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	_, err = db.Exec("INSERT INTO users (name, username, email, role, password) VALUES ($1, $2, $3, $4, $5)", user.Name, user.Username, user.Email, user.Role, password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert user into database")
	}

	return echo.NewHTTPError(http.StatusCreated, "User created successfully")
}

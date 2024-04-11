package handlers

import (
	"cost-guardian-api/db"
	"cost-guardian-api/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetAllUsers godoc
// @Summary Get all users
// @Description Fetch all users from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Security Bearer
// @Success 200 {array} models.User
// @Failure 500 {object} string "Failed to connect to database"
// @Failure 500 {object} string "Failed to fetch users from database"
// @Router /users [get]
func GetAllUsers(c echo.Context) error {
	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, username, email, role FROM users")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch users from database")
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Role); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to scan user row")

		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error iterating over user rows")
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Description Get a specific user from the database by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Security Bearer
// @Failure 400 {object} string "Invalid user ID"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Failed to connect to database"
// @Failure 500 {object} string "Failed to fetch user from database"
// @Router /users/{id} [get]
func GetUserByID(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")

	}
	defer db.Close()

	row := db.QueryRow("SELECT id, name, username, email, role FROM users WHERE id = $1", userID)

	var user models.User

	err = row.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch user from database")
	}

	return c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update an existing user with the provided information
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID to be updated"
// @Security Bearer
// @Param user body models.User true "Updated user information"
// @Success 200 {object} models.User "Successfully updated user"
// @Failure 400 {object} string "Invalid user ID"
// @Failure 400 {object} string "Failed to decode request body"
// @Failure 500 {object} string "Failed to connect to database"
// @Failure 500 {object} string "Failed to update user in database"
// @Router /users/{id} [put]
func UpdateUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	var updatedUser models.User
	if err := c.Bind(&updatedUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body")
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	_, err = db.Exec("UPDATE users SET name = $1, username = $2, email = $3 WHERE id = $4",
		updatedUser.Name, updatedUser.Username, updatedUser.Email, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update user in database")
	}

	return c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user with the provided ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID to be deleted"
// @Security Bearer
// @Success 200 {object} string "User deleted successfully"
// @Failure 400 {object} string "Invalid user ID"
// @Failure 500 {object} string "Failed to connect to database"
// @Failure 500 {object} string "Failed to delete user from database"
// @Router /users/{id} [delete]
func DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete user from database")
	}

	return echo.NewHTTPError(http.StatusOK, "User deleted successfully")
}

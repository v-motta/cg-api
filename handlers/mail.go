package handlers

import (
	"cost-guardian-api/models"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/labstack/echo/v4"
)

func SendEmail(c echo.Context) error {
	emailAddress := os.Getenv("EMAIL_ADDRESS")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	from := emailAddress
	password := emailPassword

	var mail models.Mail
	if err := c.Bind(&mail); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request body")
	}

	to := []string{
		mail.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: Cost Guardian\n\n" + mail.Message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send email")
	}

	return echo.NewHTTPError(http.StatusOK, "Email Sent Successfully!")
}

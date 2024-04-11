package handlers

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"github.com/labstack/echo/v4"
)

func SendEmail(c echo.Context) error {
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	from := "strogonoffbrasileiro12@gmail.com"
	password := emailPassword

	to := []string{
		"vinny2001.2001@gmail.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("This is a test email message.")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	fmt.Println("EmailPass:", emailPassword)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send email")
	}
	fmt.Println("Email Sent Successfully!")
	return echo.NewHTTPError(http.StatusOK, "Email Sent Successfully!")
}

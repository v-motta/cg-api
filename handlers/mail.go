package handlers

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmail() {

	emailPassword := os.Getenv("EMAIL_PASSWORD")
	// Sender data.
	from := "strogonoffbrasileiro12@gmail.com"
	password := emailPassword

	// Receiver email address.
	to := []string{
		"vinny2001.2001@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	message := []byte("This is a test email message.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	fmt.Println("EmailPass:", emailPassword)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

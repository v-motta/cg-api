package handlers

import (
	"cost-guardian-api/db"
	"cost-guardian-api/models"
	"database/sql"
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

	db, err := db.Connect()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to connect to database")
	}
	defer db.Close()

	var storedUser models.User
	err = db.QueryRow("SELECT name, username, email FROM users WHERE email=$1", mail.Email).Scan(&storedUser.Name, &storedUser.Username, &storedUser.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found with this email address"})
		}
		return err
	}

	to := []string{
		mail.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Redefinição de senha"
	body := `<!DOCTYPE html>
	<html lang="pt">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Redefinição de Senha</title>
			<style>
				a {
					padding: 15px;
					background-color: #002952;
					border-radius: 5px;
					color: white;
					text-decoration: none;
				}
				p {
					color: #484848;
				}
			</style>
		</head>
		<body style="font-family: Roboto, sans-serif; background-color: #cfcfcf; padding: 20px;">
			<table cellpadding="0" cellspacing="0" width="100%" style="max-width: 600px; margin: auto; background-color: #fff; border-radius: 5px; overflow: hidden; box-shadow: 0 0 10px rgba(0,0,0,0.4);">
				<tr>
					<td style="padding: 20px;">
						<p>Oi, ` + storedUser.Name + `,</p>
						<p style="padding-bottom: 20px">Você solicitou uma redefinição de senha para sua conta. Clique no botão abaixo para redefinir sua senha:</p>
						<a href="[Inserir Link de Redefinição]">Redefinir senha</a>
						<p style="padding-top: 20px">Se você não solicitou isso, pode ignorar este e-mail com segurança.</p>
					</td>
				</tr>
			</table>
		</body>
	</html>`

	message := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, body))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send email")
	}

	return echo.NewHTTPError(http.StatusOK, "Email sent successfully!")
}

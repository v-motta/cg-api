package main

import (
	"cost-guardian-api/handlers"
	"fmt"
	"os"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "cost-guardian-api/docs"
)

// @title Cost Guardian API
// @version 1.0
// @description This is a sample server Cost Guardian server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.cost.vmotta.dev
// @BasePath /api/v1

// @schemes https
// @Produces  application/json
// @Consumes  application/json

// @SecurityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @Scheme bearer
// @BearerFormat JWT
func main() {
	e := echo.New()

	g := e.Group("/api/v1")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method}  ${host}${path} ${latency_human}` + "\n",
	}))
	g.Use(middleware.Recover())

	g.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	g.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("secret"),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/api/v1/login" || c.Path() == "/api/v1/health" || c.Path() == "/api/v1/swagger/*" || c.Path() == "/api/v1/signup" || c.Path() == "/api/v1/send-email" {
				return true
			}
			return false
		},
	}))

	// @router /login [post]
	g.POST("/login", handlers.Login)

	// @router /signup [post]
	g.POST("/signup", handlers.Signup)

	// @router /health [get]
	g.GET("/health", handlers.Health)

	// @router /swagger/* [get]
	g.GET("/swagger/*", echoSwagger.WrapHandler)

	// @router /send-email [post]
	g.POST("/send-email", handlers.SendEmail)

	g.GET("/users", handlers.GetAllUsers)
	g.GET("/users/:id", handlers.GetUserByID)
	g.PUT("/users/:id", handlers.UpdateUser)
	g.DELETE("/users/:id", handlers.DeleteUser)

	fmt.Println("host: ", os.Getenv("DB_HOST"))
	fmt.Println("port: ", os.Getenv("DB_PORT"))
	fmt.Println("user: ", os.Getenv("DB_USER"))
	fmt.Println("dbname: ", os.Getenv("DB_NAME"))
	fmt.Println("sslmode: ", os.Getenv("DB_SSLMODE"))

	e.Logger.Fatal(e.Start(":5000"))
}

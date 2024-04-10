package main

import (
	"cost-guardian-api/handlers"

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

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}]  ${status}  ${method}  ${host}${path} ${latency_human}` + "\n",
	}))
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// @router /login [post]
	e.POST("/login", handlers.Login)

	// @router /health [get]
	e.GET("/health", handlers.Health)

	// @router /swagger/* [get]
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	r := e.Group("/api/v1")

	r.Use(echojwt.JWT([]byte("secret")))

	r.GET("/users", handlers.GetAllUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	e.Logger.Fatal(e.Start(":5000"))
}

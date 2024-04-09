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

	g.GET("/ping", handlers.Ping)

	// @router /login [post]
	g.POST("/login", handlers.Login)

	// @router /health [get]
	g.GET("/health", handlers.Health)

	// @router /swagger/* [get]
	g.GET("/swagger/*", echoSwagger.WrapHandler)

	g.Use(echojwt.JWT([]byte("secret")))

	g.GET("/users", handlers.GetAllUsers)
	g.GET("/users/:id", handlers.GetUserByID)
	g.POST("/users", handlers.CreateUser)
	g.PUT("/users/:id", handlers.UpdateUser)
	g.DELETE("/users/:id", handlers.DeleteUser)

	e.Logger.Fatal(e.Start(":5000"))
}

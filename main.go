package main 

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/baync180705/low-latency-matching-engine/controllers"
)

func main () {
	e := echo.New()

	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.GET("/", func(ctx echo.Context) error {
		pingRes := map[string]string{
			"status": "success",
		}
		return ctx.JSON(http.StatusOK, pingRes)
	})

	v1 := e.Group("/api/v1")
	v1.POST("/orders", controller.SubmitOrder)

	e.Logger.Fatal(e.Start(":8080"))
}
package main 

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/baync180705/low-latency-matching-engine/routes/api"
	"github.com/baync180705/low-latency-matching-engine/routes"
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

	v1.POST("/orders", api.SubmitOrder)
	v1.DELETE("/orders/:order_id", api.CancelOrder)
	v1.GET("/orderbook/:symbol", api.GetOrderBook)
	v1.GET("/orders/:order_id", api.GetOrderStatus)

	e.GET("/health", routes.HealthCheck)
	e.GET("/metrics", routes.GetMetricsHandler)


	e.Logger.Fatal(e.Start(":8080"))
}
package main

import (
	"github.com/Mr-Bellali/home_storage/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"}, 
		AllowHeaders: []string{
			echo.HeaderOrigin, 
			echo.HeaderContentType, 
			echo.HeaderAccept, 
			echo.HeaderAuthorization,
		},
	}))
	api := e.Group("/api")

	// Health check route
	api.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "API is running"})
	})
	handlers.SetupAuthRoutes(api)
	e.Logger.Fatal(e.Start("0.0.0.0:5050"))
}

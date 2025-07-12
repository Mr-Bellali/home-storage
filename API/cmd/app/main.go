package main

import (
	"github.com/Mr-Bellali/home_storage/internal/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowHeaders: []string{
			echo.HeaderOrigin, 
			echo.HeaderContentType, 
			echo.HeaderAccept, 
			echo.HeaderAuthorization,
		},
	}))
	api := e.Group("/api")
	handlers.SetupAuthRoutes(api)
	e.Logger.Fatal(e.Start(":1323"))
}

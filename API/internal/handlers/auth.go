package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SetupAuthRoutes(g *echo.Group) {
	g.POST("/auth/login", func(c echo.Context) error {
		var data map[string]string

		if err := c.Bind(&data); err != nil {
			log.Println("Error binding data:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
		}

		email := data["email"]
		password := data["password"]

		if email == "" || password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email and password are required"})
		}

		

		return c.JSON(200, map[string]string{"message": "Login endpoint"})
	})
	// e.POST("/change-password/:id", handlers.ChangePasswordHandler, middlewares.AuthMiddleware())
}
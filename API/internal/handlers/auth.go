package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupAuthRoutes(g *echo.Group) {
	g.POST("/auth/login", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "Login endpoint"})
	})
	// e.POST("/change-password/:id", handlers.ChangePasswordHandler, middlewares.AuthMiddleware())
}
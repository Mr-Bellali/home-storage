package handlers

import (
	"github.com/Mr-Bellali/home_storage/internal/models"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func SetupUsersRoutes(g *echo.Group) {
	g.GET("/users", func(c echo.Context) error {
		var users []models.User
		if err := models.DB.Find(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching users"})
		}

		return c.JSON(http.StatusOK, users)
	})

	g.GET("/users/:id", func(c echo.Context) error {
		id := c.Param("id")
		var user models.User
		if err := models.DB.First(&user, id).Error; err != nil {
			log.Println("Error fetching user:", err)
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
		}
		return c.JSON(http.StatusOK, user)
	})
}

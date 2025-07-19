package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Mr-Bellali/home_storage/internal/middlewares"
	"github.com/labstack/echo/v4"
)

func SetupWorkspacesRoutes(g *echo.Group) {
	g.POST("/workspaces", func(c echo.Context) error {
		// Get data from the request
		var data map[string]string

		if err := c.Bind(&data); err != nil {
			fmt.Println("Error binding data:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
		}

		// Validate it
		name := data["name"]
		description := data["description"]
		workspaceType := data["type"]

		fmt.Printf("Creating workspace - Name: %s, Description: %s, Type: %s\n", name, description, workspaceType)

		if name == "" || description == "" || workspaceType == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Name, description and workspace type are required!"})
		}

		// Check if workspaceType is valid
		if workspaceType != "public" && workspaceType != "personal" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Workspace type must be public or personal!"})
		}

		// Use the Docker-mounted path to your Desktop
		desktopPath := "/host-desktop/workspaces"

		// Ensure the workspaces folder exists
		err := os.MkdirAll(desktopPath, os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Error creating workspaces directory on Desktop: %s", err),
			})
		}

		// Create the specific workspace directory inside it
		dir := filepath.Join(desktopPath, name)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": fmt.Sprintf("Error creating workspace directory: %s", err),
			})
		}

		fmt.Printf("Created workspace directory: %s\n", dir)

		// Respond
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Workspace created successfully: %s", dir),
		})
	},
		middlewares.AuthMiddleware())
}

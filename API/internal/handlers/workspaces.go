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

		fmt.Printf("Creating workspace - Name: %s, Description: %s, Type: %s", name, description, workspaceType)

		if name == "" || description == "" || workspaceType == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Name, description and workspace type are required!"})
		}

		// Check if private or public workspace
		if workspaceType != "public" && workspaceType != "personal" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Workspace type must be public or personal!"})
		}

		// Create a folder (relative to current directory or project root)
		// Get current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error getting working directory: %s", err)})
		}

		// Create workspaces directory in the project root
		workspacesDir := filepath.Join(cwd, "workspaces")
		err = os.MkdirAll(workspacesDir, os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error creating workspaces directory: %s", err)})
		}

		// Create the specific workspace directory
		dir := filepath.Join(workspacesDir, name)
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": fmt.Sprintf("Error creating workspace directory: %s", err)})
		}

		fmt.Printf("Created workspace directory: %s", dir)

		// Store the metadata on the database
		// Return the response

		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Workspace created successfully: %s", dir),
		})
	},
		middlewares.AuthMiddleware())

}

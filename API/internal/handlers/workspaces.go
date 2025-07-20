package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Mr-Bellali/home_storage/internal/middlewares"
	"github.com/Mr-Bellali/home_storage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
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

		// Check if already exists
		var existingWorkspace models.Workspace
		if err := models.DB.Where("name = ?", name).First(&existingWorkspace).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Workspace already exists with this name"})
		}

		fmt.Println("existing workspace: ", existingWorkspace)

		// Create workspace
		workspace := models.Workspace{
			Name:        name,
			Description: description,
			Type:        workspaceType,
			UserId: c.Get("user_id").(uint),
		}

		if err := models.DB.Create(&workspace).Error; err != nil {
			fmt.Println(color.Red("Error creating worspace: "), err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating workspace"})
		}

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

		// Respond
		return c.JSON(http.StatusOK, map[string]string{
			"message": fmt.Sprintf("Workspace created successfully: %s", dir),
		})
	}, middlewares.AuthMiddleware())

	g.GET("/workspaces/:id", func(c echo.Context) error {
		// Get the id param 
		workspaceId, err := strconv.Atoi(c.Param("id")) 
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Currapted Id must be a valid number"})
		}

		// Get the workspace by its id
		var workspace models.Workspace
		models.DB.First(&workspace, workspaceId)

		return c.JSON(http.StatusAccepted, map[string]string{"message":workspace.Name})
	})
}

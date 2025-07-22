package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Mr-Bellali/home_storage/internal/middlewares"
	"github.com/Mr-Bellali/home_storage/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
	"gorm.io/gorm"
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
			UserId:      c.Get("user_id").(uint),
		}

		if err := models.DB.Create(&workspace).Error; err != nil {
			fmt.Println(color.Red("Error creating worspace: "), err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating workspace"})
		}

		fmt.Printf("Creating workspace - Name: %s, Description: %s\n", name, description)

		if name == "" || description == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Name, description and workspace type are required!"})
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

	// Route to get all personal workspaces
	g.GET("/workspaces", func(c echo.Context) error {
		// Retrive the user ID from context
		userId := c.Get("user_id").(uint)

		// Get all user's workspaces
		var workspaces []models.Workspace
		result := models.DB.Where("user_id = ?", userId).Find(&workspaces)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message":"Unexpected error"})
		}
		return c.JSON(http.StatusAccepted, map[string]any{"worspaces":workspaces})
	}, middlewares.AuthMiddleware())

	// Route to upload a media to a workspace
	g.POST("/workspaces/:id/upload", func(c echo.Context) error {
		// Get the id param
		workspaceId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Currapted Id must be a valid number"})
		}

		// Retrive the user id from context
		userId := c.Get("user_id").(uint)

		// Get the workspace by its id and the user ID
		var workspace models.Workspace
		result := models.DB.Where("id = ? AND user_id = ?", workspaceId, userId).First(&workspace)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"message": "Workspace not found"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve workspace"})
		}

		// Read file from the reauest
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to read file"})

		}
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to open uploaded file"})
		}
		defer src.Close()

		// Destination
		workspaceDir := filepath.Join("/host-desktop/workspaces", workspace.Name)

		// Making sure that the workspace folder still exists
		err = os.MkdirAll(workspaceDir, os.ModePerm)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to ensure workspace directory exists"})
		}

		// Create destination file
		dstPath := filepath.Join(workspaceDir, file.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create file on your workspace"})
		}
		defer dst.Close()

		// Copy uploaded file to destination
		size, err := io.Copy(dst, src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save file"})
		}

		// Detect MIME type
		mimeType := file.Header.Get("Content-Type")

		// Save metadata in database
		fileRecord := models.File{
			Filename:    file.Filename,
			Filepath:    dstPath,
			Size:        size,
			MIMEType:    mimeType,
			UserID:      userId,
			WorkspaceID: workspace.ID,
		}

		if err := models.DB.Create(&fileRecord).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save file metadata"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "File uploaded and metadata saved successfully",
			"file":    fileRecord,
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

		return c.JSON(http.StatusAccepted, map[string]string{"message": workspace.Name})
	})
}

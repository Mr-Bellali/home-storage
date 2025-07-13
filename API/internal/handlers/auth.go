package handlers

import (
	"log"
	"net/http"

	"github.com/Mr-Bellali/home_storage/internal/models"
	"github.com/Mr-Bellali/home_storage/pkg"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func SetupAuthRoutes(g *echo.Group) {
	g.POST("/auth/register", func(c echo.Context) error {
		var data map[string]string

		if err := c.Bind(&data); err != nil {
			log.Println("Error binding data:", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid data"})
		}

		name := data["name"]
		email := data["email"]
		password := data["password"]

		if name == "" || email == "" || password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Name, email and password are required"})
		}

		// Check if user already exists
		var existingUser models.User
		if err := models.DB.Where("email = ?", email).First(&existingUser).Error; err == nil {
			return c.JSON(http.StatusConflict, map[string]string{"message": "User already exists"})
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error hashing password"})
		}

		// Create user
		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashedPassword),
		}

		if err := models.DB.Create(&user).Error; err != nil {
			log.Println("Error creating user:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating user"})
		}

		return c.JSON(http.StatusCreated, map[string]string{"message": "User created successfully"})
	})

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

		// Find user by email
		var user models.User
		if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		}

		t, err := pkg.GenerateJWT(user.Name, user.Email, user.ID)
		if err != nil {
			log.Println("JWT generation failed:", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to generate token"})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"user": map[string]interface{}{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
			"token": t,
		})
	})
}

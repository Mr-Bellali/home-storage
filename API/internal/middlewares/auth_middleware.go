package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/Mr-Bellali/home_storage/pkg"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Missing authorization header"})
			}

			// Check if it's a Bearer token
			if !strings.HasPrefix(auth, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid authorization format"})
			}

			// Extract token from "Bearer <token>"
			tokenString := strings.TrimPrefix(auth, "Bearer ")

			// Parse and validate the token
			token, err := jwt.ParseWithClaims(tokenString, &pkg.JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return pkg.JWTSecret, nil
			})

			if err != nil || !token.Valid {
				log.Println("Token validation error:", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid or expired token"})
			}

			// Extract claims
			if claims, ok := token.Claims.(*pkg.JwtCustomClaims); ok {
				// Set user information in context
				c.Set("user_id", claims.UserID)
				c.Set("user_name", claims.Name)
				c.Set("user_email", claims.Email)
				c.Set("claims", claims)
			}

			return next(c)
		}
	}
}

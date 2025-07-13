package pkg

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("secret")

type JwtCustomClaims struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	UserID  uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(name, email string, userID uint) (string, error) {
	claims := &JwtCustomClaims{
		Name:    name,
		Email:   email,
		UserID:  userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 9)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(JWTSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return t, nil
}

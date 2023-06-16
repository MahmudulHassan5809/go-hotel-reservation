package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(ctx *fiber.Ctx) error {
	fmt.Println("JWT middleware")
	token, ok := ctx.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if err := parseToken(token); err != nil {
		return err
	}
	fmt.Println("token:", token)
	return nil
}  

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthenticated")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return fmt.Errorf("unauthenticated")
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
	return fmt.Errorf("unauthenticated")
}
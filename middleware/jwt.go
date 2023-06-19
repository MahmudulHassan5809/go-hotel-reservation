package middleware

import (
	"fmt"
	"hotel-reservation/db"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Println("JWT middleware")
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrUnAuthorized()
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)
		if time.Now().Unix() > expires {
			return NewError(http.StatusUnauthorized, "token expired")
		}

		UserID := claims["id"].(string)
		user, err := userStore.GetUseById(ctx.Context(), UserID)
		if err != nil {
			return ErrUnAuthorized()
		}
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}  

func validateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Unexpected signing method", token.Header["alg"])
			return nil, ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, ErrUnAuthorized()
	}
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else {
		fmt.Println(err)
	}
	return nil, ErrUnAuthorized()
}
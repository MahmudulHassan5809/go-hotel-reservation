package api

import (
	"errors"
	"fmt"
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore	
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(ctx *fiber.Ctx) error {
	var params AuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return err 
	}

	user, err := h.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("INVALID CREDENTIALS")
		}
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		return fmt.Errorf("INVALID CREDENTIALS")
	}
	fmt.Println(user)
	return nil
}
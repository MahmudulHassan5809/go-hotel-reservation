package api

import (
	"errors"
	"hotel-reservation/db"
	"hotel-reservation/types"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}


func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	if err := h.userStore.DeleteUser(ctx.Context(), userId); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"message": "deleted"})
}

func (h *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	return nil
	// var (
	// 	params  map[string]any
	// 	userID  ctx.Params("id") 
	// )
	// if err := ctx.BodyParser(&params); err != nil {
	// 	return err
	// }
	// filter := db.Map{"_id": userID}
	// if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
	// 	return err
	// }
	// return c.JSON(map[string]string{"updated": userID})
}

func (h *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.CreateUser(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}


func (h *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)
	
	user, err := h.userStore.GetUseById(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return ctx.JSON(map[string]string{"message": "not found"})
		}
		return err
	}
	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}


package api

import (
	"hotel-reservation/middleware"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if apiError, ok := err.(Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	if apiError, ok := err.(middleware.Error); ok {
		return ctx.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(http.StatusInternalServerError, err.Error())
	return ctx.Status(apiError.Code).JSON(apiError)
}

type Error struct {
	Code int `json:"code"`
	Err  string `json:"error"`
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "invalid id given",
	}
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err: "unauthorized request",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: http.StatusBadRequest,
		Err: "unauthorized request",
	}
}

func ErrResourceNotFound() Error {
	return Error{
		Code: http.StatusNotFound,
		Err: "resource not found",
	}
}
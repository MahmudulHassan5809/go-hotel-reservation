package api

import (
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store db.Store
}

func NewBookingHandler(store db.Store) *BookingHandler{
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(ctx *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}
	if booking.UserID != user.ID{
		return ctx.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg: "not authorized",
		})
	}
	return ctx.JSON(booking)
}
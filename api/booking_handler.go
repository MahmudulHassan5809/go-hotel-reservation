package api

import (
	"hotel-reservation/db"

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
	return ctx.JSON(booking)
}
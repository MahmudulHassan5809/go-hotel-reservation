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
		return ErrResourceNotFound()
	}
	return ctx.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return ErrResourceNotFound()
	}
	user, err := getAuthUser(ctx)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID{
		return ErrUnAuthorized()
	}
	return ctx.JSON(booking)
}


func (h *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingByID(ctx.Context(), id)
	if err != nil {
		return err
	}
	user, err := getAuthUser(ctx)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID{
		return ErrUnAuthorized()
	}
	if err := h.store.Booking.UpdateBooking(ctx.Context(), ctx.Params("id"), bson.M{"cancelled": true}); err != nil {
		return err
	}
	return ctx.JSON(genericResp{Type: "msg", Msg: "updated"})
}
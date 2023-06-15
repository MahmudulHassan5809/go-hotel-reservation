package api

import (
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore: rs,
	}
}

func (h *HotelHandler) HandleGetHotels(ctx *fiber.Ctx) error {
	hotels, err := h.hotelStore.GetHotels(ctx.Context(), nil)
	if err != nil {
		return err
	}
	return ctx.JSON(hotels)
}
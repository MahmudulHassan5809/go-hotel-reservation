package api

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type BookRoomParams struct {
	FromDate time.Time `bson:"fromDate" json:"fromDate"`
	TillDate time.Time `bson:"tillDate" json:"tillDate"`
	NumPersons int `bson:"numPersons" json:"numPersons"`
}

func (p BookRoomParams) Validate() map[string]string {
	errors := map[string]string{}
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		errors["FromDate/TillDate"] = "cannot book a room in the past"
	}

	return errors
}

type RoomHandler struct {
	store db.Store
}

func NewRoomHandler(store db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	var params BookRoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}
	roomID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}
	user, ok := ctx.Context().Value("user").(*types.User)
	if !ok {
		return ctx.Status(http.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg: "internal server error",
		})
	}

	ok, err = h.isRoomAvailableForBooking(ctx.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return ctx.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", ctx.Params("id")),
		})
	}


	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	}
	fmt.Println(booking)
	inserted, err := h.store.Booking.InsertBooking(ctx.Context(), &booking)
	if err != nil {
		return err
	}
	return ctx.JSON(inserted)
}


func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil
}
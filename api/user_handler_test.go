package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const (
	testMongoUri = "mongodb://localhost:27017"
	dbname = "hotel-reservation-test"
)

type testDb struct {
	db.UserStore
}

func (tdb *testDb) tearDown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testMongoUri))
	if err != nil {
		log.Fatal(err)
	}

	return &testDb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestCreateUser(t *testing.T){
	tdb := setup(t)
	defer tdb.tearDown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email: "test@gmail.com",
		FirstName: "Test",
		LastName: "Last",
		Password: "1234567",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	// bb, _ := io.ReadAll(resp.Body)
	// fmt.Println(user)

	if len(user.ID) == 0 {
		t.Errorf("Expecting a use id to be set")
	}

	if len(user.Password) > 0 {
		t.Errorf("Expecting the password not to be included in the response")
	}

	if user.FirstName != params.FirstName{
		t.Errorf("expected FirstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName{
		t.Errorf("expected LastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email{
		t.Errorf("expected Email %s but got %s", params.Email, user.Email)
	}
}
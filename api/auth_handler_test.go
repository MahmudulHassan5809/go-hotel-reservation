package api

import (
	"bytes"
	"context"
	"encoding/json"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)


func makeTestUser(t *testing.T,userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			FirstName: "Mahmudul",
			LastName: "Hassan",
			Email: "mahmudul@gmail.com",
			Password: "1234567",
		},
	)

	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}


func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	insertedUser := makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "mahmudul@gmail.com",
		Password: "1234567",
	}
	b, _ := json.Marshal(params)
	
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK{
		t.Fatalf("expected http status of 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	
	if len(authResp.Token) == 0 {
		t.Fatalf("expected the JWT token in response")
	}

	insertedUser.Password = ""
	if !reflect.DeepEqual(insertedUser, authResp.User){
		t.Fatalf("expected the user to be equal to be the equal inserted user")
	}
}


func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.tearDown(t)

	makeTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email: "mahmudul@gmail.com",
		Password: "123456755",
	}
	b, _ := json.Marshal(params)
	
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest{
		t.Fatalf("expected http status of 400 but got %d", resp.StatusCode)
	}

	var genResp genericResp
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		t.Fatal(err)
	}

	if genResp.Type != "error"{
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)
	}

	if genResp.Msg != "invalid credentials"{
		t.Fatalf("expected gen response type to be <invalid credentials> but got %s", genResp.Type)
	}
}
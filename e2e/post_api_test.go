package e2e

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

//
// Utils
//

func happyPathData() (username, body string) {
	username = "my_username_test"
	body = "My Body Description Test"
	return
}

//
// E2E Testing
//

func TestCreatePost_ShouldReturnStatusBadRequest_WhenItHasInvalidBody(t *testing.T) {
	api := NewApiClient()
	params := []map[string]string{
		nil,
		{},
		{"other": "value"},
		{"username": "user", "body": "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since"},
	}

	for _, p := range params {
		resp, err := api.Post("/posts", p)
		if err != nil {
			t.Fatal(err.Error())
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf(
				"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
				http.StatusBadRequest,
				resp.Status,
			)
		}
	}
}

func TestDeletePost_ShouldReturnStatusNotFound_WhenPostIdIsNotOnDatabase(t *testing.T) {
	api := NewApiClient()
	id := uuid.NewString()

	resp, err := api.Delete("/posts/" + id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusNotFound,
			resp.Status,
		)
	}
}

func TestGetPost_ShouldReturnStatusNotFound_WhenPostIdIsNotOnDatabase(t *testing.T) {
	api := NewApiClient()
	id := uuid.NewString()

	resp, err := api.Get("/posts/" + id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusNotFound,
			resp.Status,
		)
	}
}

func TestPostApi_HappyPath(t *testing.T) {
	t.Log("*** Start Post Successful")

	id := createSuccessfully(t)
	readSuccessfully(id, t)
	// updateSuccessfully(id, t)
	deleteSuccessfully(id, t)

	t.Log("*** End Post Successful")
}

func createSuccessfully(t *testing.T) string {
	t.Log("*** Create Post")
	api := NewApiClient()
	username, body := happyPathData()
	payload := map[string]string{
		"username": username,
		"body":     body,
	}

	resp, err := api.Post("/posts", payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusCreated,
			resp.Status,
		)
	}

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	id := res["id"].(string)

	if id == "" {
		t.Fatal("Invalid ID")
	}

	if res["username"].(string) != payload["username"] {
		t.Fatal("Invalid Username")
	}

	if res["body"].(string) != payload["body"] {
		t.Fatal("Invalid Body")
	}

	if res["created_at"].(string) == "0001-01-01T00:00:00Z" {
		t.Fatal("Invalid CreatedAt")
	}

	return id
}

func readSuccessfully(id string, t *testing.T) {
	t.Log("*** Read Post")
	username, body := happyPathData()
	api := NewApiClient()

	resp, err := api.Get("/posts/" + id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusOK,
			resp.Status,
		)
	}

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	if res["id"].(string) != id {
		t.Fatal("Invalid ID")
	}

	if res["username"].(string) != username {
		t.Fatal("Invalid Username")
	}

	if res["body"].(string) != body {
		t.Fatal("Invalid Body")
	}
}

func updateSuccessfully(id string, t *testing.T) {
	t.Log("*** Update Post")
	api := NewApiClient()

	payload := map[string]string{
		"username": "new_user",
		"body":     "Other body",
	}

	resp, err := api.Put("/posts/"+id, payload)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusCreated,
			resp.Status,
		)
	}

	res, err := api.ParseBody(resp)
	if err != nil {
		t.Fatal(err.Error())
	}

	if res["id"].(string) == id {
		t.Fatal("Invalid ID")
	}

	if res["username"].(string) != payload["username"] {
		t.Fatal("Invalid Username")
	}

	if res["body"].(string) != payload["body"] {
		t.Fatal("Invalid Body")
	}
}

func deleteSuccessfully(id string, t *testing.T) {
	t.Log("*** Delete Post")
	api := NewApiClient()

	resp, err := api.Delete("/posts/" + id)
	if err != nil {
		t.Fatal(err.Error())
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf(
			"Invalid Status Code. Expected Status \"%d\" and received \"%s\"",
			http.StatusNoContent,
			resp.Status,
		)
	}
}

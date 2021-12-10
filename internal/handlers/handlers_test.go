package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// // all GET tests
	// {"home", "/", "GET", []postData{}, http.StatusOK},
	// {"about", "/about", "GET", []postData{}, http.StatusOK},
	// {"generals-quarters", "/gq", "GET", []postData{}, http.StatusOK},
	// {"majors-suite", "/ms", "GET", []postData{}, http.StatusOK},
	// {"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	// {"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// {"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	// {"contact", "/contact", "GET", []postData{}, http.StatusOK},
	// // all POST tests
	// {"POST-search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start_date", value: "30/11/2021"},
	// 	{key: "end_date", value: "03/12/2021"},
	// }, http.StatusOK},
	// {"POST-search-availability-modal", "/search-availability-modal", "POST", []postData{
	// 	{key: "start_date", value: "30/11/2021"},
	// 	{key: "end_date", value: "03/12/2021"},
	// }, http.StatusOK},
	// {"POST-make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Peter"},
	// 	{key: "last_name", value: "Barrett"},
	// 	{key: "email", value: "peter@barrett.com"},
	// 	{key: "phone", value: "01508 000000"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := testServer.Client().Get(testServer.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for handler %s, status code %d was expected but %d was received", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := testServer.Client().PostForm(testServer.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for handler %s, status code %d was expected but %d was received", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Test Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}
}

// getCtx creates a context for use in TestRepository_Reservation() request
func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

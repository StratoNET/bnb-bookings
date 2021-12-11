package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

// type postData struct {
// 	key   string
// 	value string
// }

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	// all GET tests
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"generals-quarters", "/gq", "GET", http.StatusOK},
	{"majors-suite", "/ms", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},

	// // all POST tests & use session
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

// TestHandlers tests GET handlers only
func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, e := range theTests {
		resp, err := testServer.Client().Get(testServer.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for handler %s, status code %d was expected but %d was received", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	// necessary to create a reservation with room information as would typically be pulled from session at beginning of Reservation()
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
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}

	// test case where RESERVATION IS NOT IN THE SESSION (= reset everything)
	// reinitialise request & context
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// session has now been created but without a reservation, reinitialise 'fake' response writer
	rr = httptest.NewRecorder()

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with NON-EXISTENT ROOM
	// reinitialise request & context
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// session has now been created but without a reservation, reinitialise 'fake' response writer
	rr = httptest.NewRecorder()
	// invoke case of room ID > number of rooms available
	reservation.RoomID = 999
	// put reservation into session
	session.Put(ctx, "reservation", reservation)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	// test EVERYTHING IS CORRECT

	// necessary to create a reservation with room information as would typically be pulled from session at beginning of PostReservation()
	// in fact, start & end dates would also be available at beginning of PostReservation() and could be incorporated in reservation model
	reservation := models.Reservation{
		RoomID: 2,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}
	// create a request body for reservation data to be posted
	postedData := url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	// postedData.Add("room_id", "1") ... already known
	postedData.Add("first_name", "Joe")
	postedData.Add("last_name", "Soap")
	postedData.Add("email", "joe@soap.bar")
	postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// test for MISSING POST BODY
	// initialise request & context
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// =====================================================================================================================================

// getCtx creates a context for use in TestRepository_Reservation() request
func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

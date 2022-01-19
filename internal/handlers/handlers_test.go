package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/database"
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
	{"not-found", "/non/existent/route", "GET", http.StatusNotFound},
	{"login", "/login", "GET", http.StatusOK},
	{"logout", "/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "GET", http.StatusOK},
	{"fail to get new reservations", "/admin/reservations-new/", "GET", http.StatusNotFound},
	{"all reservations", "/admin/reservations-all", "GET", http.StatusOK},
	{"reservations calendar", "/admin/reservations-cal", "GET", http.StatusOK},
	{"reservations calendar with params", "/admin/reservations-cal?y=2022&m=12", "GET", http.StatusOK},
	{"show reservation page", "/admin/reservations/new/1/page", "GET", http.StatusOK},
	{"reservation id missing", "/admin/reservations/all/one/page", "GET", http.StatusOK},
	{"non-existent reservation", "/admin/reservations/all/0/page", "GET", http.StatusOK},

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

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
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

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
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

	// test case where RESERVATION IS NOT IN THE SESSION (= reset everything)
	// reinitialise request & context
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// session has now been created but without a reservation, reinitialise 'fake' response writer
	rr = httptest.NewRecorder()

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
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

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure of FORM VALIDATION
	// create a request body for reservation data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	// postedData.Add("room_id", "1") ... already known
	postedData.Add("first_name", "J") // will fail validation
	postedData.Add("last_name", "Soap")
	postedData.Add("email", "joe@soap.bar")
	postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}

	// test for failure to INSERT ROOM RESERVATION
	// create a request body for reservation data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	// postedData.Add("room_id", "1") ... already known
	postedData.Add("first_name", "Joe")
	postedData.Add("last_name", "Soap")
	postedData.Add("email", "joe@soap.bar")
	postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	reservation.RoomID = 99 //will fail to insert reservation, room number does not exist

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// test for failure to INSERT ROOM RESTRICTION
	// create a request body for reservation data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	// postedData.Add("room_id", "1") ... already known
	postedData.Add("first_name", "Joe")
	postedData.Add("last_name", "Soap")
	postedData.Add("email", "joe@soap.bar")
	postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	reservation.RoomID = 9999 //will fail to insert reservation, room number does not exist

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_PostAvailabilityModal(t *testing.T) {
	// 1. NO ROOMS AVAILABLE
	// create a request body for reservation data to be posted
	postedData := url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	postedData.Add("room_id", "1")
	// postedData.Add("first_name", "Joe")
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ := http.NewRequest("POST", "/search-availability-modal", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailabilityModal)

	handler.ServeHTTP(rr, req)

	var jr jsonResponse

	err := json.Unmarshal([]byte(rr.Body.Bytes()), &jr)

	if err != nil {
		t.Error("PostAvailabilityModal handler failed to parse JSON", err)
	}

	// 2. CANNOT PARSE FORM (no form body)
	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability-modal", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailabilityModal)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &jr)

	if err != nil && !jr.Ok {
		t.Error("PostAvailabilityModal handler failed to parse JSON", err)
	}

	// 3. CANNOT CONNECT TO DATABASE
	// create a request body for reservation data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	postedData.Add("room_id", "1")
	// postedData.Add("first_name", "Joe")
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")
	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability-modal", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailabilityModal)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &jr)

	if err != nil && !jr.Ok {
		t.Error("PostAvailabilityModal handler failed to parse JSON & did NOT connect to database", err)
	}

	// 4. ROOMS AVAILABLE
	// create a request body for reservation data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "02/01/2099")
	postedData.Add("end_date", "03/01/2099")
	postedData.Add("room_id", "1")
	// postedData.Add("first_name", "Joe")
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")
	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability-modal", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailabilityModal)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.Bytes()), &jr)

	if err != nil && !jr.Ok {
		t.Error("PostAvailabilityModal did NOT find availability when there is !", err)
	}
}

func TestNewRepository(t *testing.T) {
	var db database.DB
	testRepository := NewRepository(&app, &db)

	if reflect.TypeOf(testRepository).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepository(): got %s, wanted *Repository", reflect.TypeOf(testRepository).String())
	}
}

func TestRepository_PostAvailability(t *testing.T) {
	// 1. ROOMS ARE NOT AVAILABLE
	// create a request body for availability data to be posted
	postedData := url.Values{}
	postedData.Add("start_date", "01/01/2099")
	postedData.Add("end_date", "02/01/2099")
	// postedData.Add("room_id", "1")
	// postedData.Add("first_name", "Joe")
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler (NO ROOMS AVAILABLE) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 2. test for MISSING POST BODY
	// reinitialise request & context
	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, reinitialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler (MISSING POST BODY) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 3. test for ROOMS AVAILABLE
	// create a request body for availability data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2089")
	postedData.Add("end_date", "02/01/2089")
	// postedData.Add("room_id", "1")
	// postedData.Add("first_name", "J") // will fail validation
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("PostAvailability handler (ROOMS AVAILABLE) returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}

	// 4. test for INVALID START DATE
	// create a request body for availability data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "02/01/2089")
	// postedData.Add("room_id", "1")
	// postedData.Add("first_name", "J") // will fail validation
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler (INVALID START DATE) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 5. test for INVALID END DATE
	// create a request body for availability data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2089")
	postedData.Add("end_date", "invalid")
	// postedData.Add("room_id", "1")
	// postedData.Add("first_name", "J") // will fail validation
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler (INVALID END DATE) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 6. test for DATABASE QUERY FAILURE
	// create a request body for availability data to be posted
	postedData = url.Values{}
	postedData.Add("start_date", "01/01/2100")
	postedData.Add("end_date", "03/01/2100")
	// postedData.Add("room_id", "1")
	// postedData.Add("first_name", "J") // will fail validation
	// postedData.Add("last_name", "Soap")
	// postedData.Add("email", "joe@soap.bar")
	// postedData.Add("phone", "01234 567890")

	// initialise request & context, incorporate & encode data in post body
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// session has now been incorporated, initialise 'fake' response writer
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostAvailability)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler (DATABASE QUERY FAILURE) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	// necessary to create a reservation with room & dates information as would typically be within session before displaying summary
	layout := "02/01/2006"
	sd, _ := time.Parse(layout, "01/01/2080")
	ed, _ := time.Parse(layout, "09/01/2080")
	reservation := models.Reservation{
		RoomID:    2,
		StartDate: sd,
		EndDate:   ed,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationSummary handler returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_ChooseRoom(t *testing.T) {
	// 1. Normal operation, reservation within session, URL parameter viable
	// necessary to create a reservation with room & dates information as would typically be within session before displaying summary
	layout := "02/01/2006"
	sd, _ := time.Parse(layout, "01/01/2080")
	ed, _ := time.Parse(layout, "09/01/2080")
	reservation := models.Reservation{
		RoomID:    2,
		StartDate: sd,
		EndDate:   ed,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/2", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// set RequestURI on the request after WithContext() & get id
	req.RequestURI = "/choose-room/2"

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler (OPERATING CORRECTLY) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 2. Bad URL
	req, _ = http.NewRequest("GET", "/choose-room/two", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set RequestURI on the request after WithContext() & get id
	req.RequestURI = "/choose-room/two"

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr = httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ChooseRoom)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler (BAD URL) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 3. Reservation missing from session
	req, _ = http.NewRequest("GET", "/choose-room/2", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// set RequestURI on the request after WithContext() & get id
	req.RequestURI = "/choose-room/2"

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler (MISSING RESERVATION) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_ReserveRoom(t *testing.T) {
	// necessary to create a reservation with room & dates information as would typically be within session before reserving room
	layout := "02/01/2006"
	sd, _ := time.Parse(layout, "01/01/2080")
	ed, _ := time.Parse(layout, "09/01/2080")
	reservation := models.Reservation{
		RoomID:    2,
		StartDate: sd,
		EndDate:   ed,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}

	// 1. successful room reservation / database operation
	req, _ := http.NewRequest("GET", "/reserve-room?sd=01/01/2099&ed=02/01/2099&id=2", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReserveRoom)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ReserveRoom handler (OPERATING CORRECTLY) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}

	// 2. cannot get room name from database
	req, _ = http.NewRequest("GET", "/reserve-room?sd=01/01/2099&ed=02/01/2099&id=99", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReserveRoom)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ReserveRoom handler (ROOM NAME MISSING) returned code: %d, expected code: %d", rr.Code, http.StatusSeeOther)
	}
}

func TestRepository_AdminReservationsCalendar(t *testing.T) {
	// necessary to create a reservation with room & dates information as would typically be within session before displaying summary
	// layout := "02/01/2006"
	// sd, _ := time.Parse(layout, "01/01/2080")
	// ed, _ := time.Parse(layout, "09/01/2080")
	// reservation := models.Reservation{
	// 	RoomID:    2,
	// 	StartDate: sd,
	// 	EndDate:   ed,
	// 	Room: models.Room{
	// 		ID:       2,
	// 		RoomName: "Major's Suite",
	// 	},
	// }

	req, _ := http.NewRequest("GET", "/reservations-cal", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	blocks := 1
	reservations := 1
	now := time.Now()
	bm := make(map[string]int)
	rm := make(map[string]int)

	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	for i := 0; i < 2; i++ {

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2-01-2006")] = 0
			bm[d.Format("2-01-2006")] = 0
		}

		if blocks > 0 {
			bm[firstOfMonth.Format("2-01-2006")] = blocks
		}

		if reservations > 0 {
			rm[lastOfMonth.Format("2-01-2006")] = reservations
		}

		session.Put(ctx, fmt.Sprintf("block_map_%d", i), bm)

	}

	// following code emulates a request/response lifecycle, therefor do NOT need to call any routes
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.AdminReservationsCalendar)

	// serve request using 'fake' response writer
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("AdminReservationsCalendar handler returned code: %d, expected code: %d", rr.Code, http.StatusOK)
	}

}

// ================================================= USING TABLE TESTS FROM HERE ON =======================================================

var loginTests = []struct {
	name, email, expectedHTML, expectedLocation string
	expectedStatusCode                          int
}{
	{
		"valid-credentials",
		"peter@barrett.com",
		"",
		"/",
		http.StatusSeeOther,
	},
	{
		"invalid-credentials",
		"fred@bloggs.com",
		"",
		"/login",
		http.StatusSeeOther,
	},
	{
		"invalid-data-supplied",
		// in this case an invalid email
		"useless#rubbish.com",
		// this HTML extract will always be returned from login page in response body
		`action="/login"`,
		// no location is involved i.e. no redirect, just a straightforward render (200)
		"",
		http.StatusOK,
	},
}

func TestRepository_PostLogin(t *testing.T) {
	// range through all tests
	for _, v := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", v.email)
		postedData.Add("password", "password")

		// initialise request & context, incorporate & encode data in post body
		req, _ := http.NewRequest("POST", "/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// initialise 'fake' response writer
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostLogin)

		handler.ServeHTTP(rr, req)

		if rr.Code != v.expectedStatusCode {
			t.Errorf("PostLogin handler (%s) returned code: %d, expected code: %d", v.name, rr.Code, v.expectedStatusCode)
		}

		// test expected location
		if v.expectedLocation != "" {
			// get actual URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != v.expectedLocation {
				t.Errorf("PostLogin handler (%s) returned location: %s, expected location: %s", v.name, actualLoc.String(), v.expectedLocation)
			}
		}

		// test expected HTML values
		if v.expectedHTML != "" {
			// read response body into string to get HTML
			html := rr.Body.String()
			if !strings.Contains(html, v.expectedHTML) {
				t.Errorf("PostLogin handler (%s) returned HTML: %s, expected HTML: %s", v.name, html, v.expectedHTML)
			}
		}
	}
}

// ========================================================================================================================================

var adminPostReservationTests = []struct {
	name, url, expectedLocation, expectedHTML string
	expectedStatusCode                        int
	postedData                                url.Values
}{
	{
		name:               "valid-data-(NEW)",
		url:                "/admin/reservations/new/1/page",
		expectedLocation:   "/admin/reservations-new",
		expectedHTML:       "",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"first_name": {"Johnny"},
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
		},
	},
	{
		name:               "valid-data-(ALL)",
		url:                "/admin/reservations/all/1/page",
		expectedLocation:   "/admin/reservations-all",
		expectedHTML:       "",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"first_name": {"Johnny"},
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
		},
	},
	{
		name:               "valid-data-(CAL)",
		url:                "/admin/reservations/cal/1/page",
		expectedLocation:   "/admin/reservations-cal?y=2022&m=12",
		expectedHTML:       "",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"first_name": {"Johnny"},
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
			"y":          {"2022"},
			"m":          {"12"},
		},
	},
	{
		name: "invalid-data-(ALL)",
		url:  "/admin/reservations/all/1/page",
		// no location is involved i.e. no redirect, just a straightforward render (200)
		expectedLocation:   "",
		expectedHTML:       "",
		expectedStatusCode: http.StatusOK,
		postedData: url.Values{
			"first_name": {"J"}, // minimum 2 letters required
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
		},
	},
	{
		name:               "cannot get reservation by id",
		url:                "/admin/reservations/all/0/page",
		expectedLocation:   "/admin/reservations-all",
		expectedHTML:       "",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"first_name": {"Johnny"},
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
		},
	},
	{
		name:               "cannot update reservation",
		url:                "/admin/reservations/all/1/page",
		expectedLocation:   "/admin/reservations-all",
		expectedHTML:       "",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"first_name": {"Bonzo"},
			"last_name":  {"Begood"},
			"email":      {"johnny@begood.com"},
			"phone":      {"0123456789"},
		},
	},
}

func TestRepository_AdminPostReservation(t *testing.T) {
	// necessary to create a reservation with room & dates information as would typically be within session before getting it by id & updating
	layout := "02/01/2006"
	sd, _ := time.Parse(layout, "01/01/2080")
	ed, _ := time.Parse(layout, "09/01/2080")
	reservation := models.Reservation{
		RoomID:    2,
		StartDate: sd,
		EndDate:   ed,
		Room: models.Room{
			ID:       2,
			RoomName: "Major's Suite",
		},
	}

	// range through all tests
	for _, v := range adminPostReservationTests {
		// initialise request & context, incorporate & encode data in post body
		var req *http.Request
		if v.postedData != nil {
			req, _ = http.NewRequest("POST", "/login", strings.NewReader(v.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/login", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		req.RequestURI = v.url

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// initialise 'fake' response writer
		rr := httptest.NewRecorder()

		session.Put(ctx, "reservation", reservation)

		handler := http.HandlerFunc(Repo.AdminPostReservation)

		handler.ServeHTTP(rr, req)

		if rr.Code != v.expectedStatusCode {
			t.Errorf("AdminPostReservation handler (%s) returned code: %d, expected code: %d", v.name, rr.Code, v.expectedStatusCode)
		}

		// test expected location
		if v.expectedLocation != "" {
			// get actual URL from test
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != v.expectedLocation {
				t.Errorf("AdminPostReservation handler (%s) returned location: %s, expected location: %s", v.name, actualLoc.String(), v.expectedLocation)
			}
		}

		// test expected HTML values
		if v.expectedHTML != "" {
			// read response body into string to get HTML
			html := rr.Body.String()
			if !strings.Contains(html, v.expectedHTML) {
				t.Errorf("AdminPostReservation handler (%s) returned HTML: %s, expected HTML: %s", v.name, html, v.expectedHTML)
			}
		}
	}
}

// ========================================================================================================================================

var adminPostReservationsCalendarTests = []struct {
	name, expectedLocation, expectedHTML     string
	expectedStatusCode, blocks, reservations int
	postedData                               url.Values
}{
	{
		name:               "cal",
		expectedStatusCode: http.StatusSeeOther,
		postedData: url.Values{
			"year":  {time.Now().Format("2006")},
			"month": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2-01-2006")): {"1"},
		},
	},
	{
		name:               "cal-blocks",
		postedData:         url.Values{},
		expectedStatusCode: http.StatusSeeOther,
		blocks:             1,
	},
	{
		name:               "cal-res",
		postedData:         url.Values{},
		expectedStatusCode: http.StatusSeeOther,
		reservations:       1,
	},
}

func TestRepository_AdminPostReservationsCalendar(t *testing.T) {
	for _, v := range adminPostReservationsCalendarTests {
		var req *http.Request
		if v.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-cal", strings.NewReader(v.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-cal", nil)
		}
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2-01-2006")] = 0
			bm[d.Format("2-01-2006")] = 0
		}

		if v.blocks > 0 {
			bm[firstOfMonth.Format("2-01-2006")] = v.blocks
		}

		if v.reservations > 0 {
			rm[lastOfMonth.Format("2-01-2006")] = v.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminPostReservationsCalendar)
		handler.ServeHTTP(rr, req)

		if rr.Code != v.expectedStatusCode {
			t.Errorf("AdminPostReservationsCalendar handler (%s) returned code: %d, expected code: %d", v.name, rr.Code, v.expectedStatusCode)
		}
	}
}

// ========================================================================================================================================

var adminReservationProcessTests = []struct {
	name, queryParams, expectedLocation string
	expectedStatusCode                  int
}{
	{
		name:               "reservation-processed",
		queryParams:        "",
		expectedLocation:   "",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "reservation-processed-cal",
		queryParams:        "?y=2022&m=10",
		expectedLocation:   "",
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestRepository_AdminReservationProcess(t *testing.T) {
	for _, v := range adminReservationProcessTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/reservation-processed/cal/1/page%s", v.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AdminReservationProcess)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusSeeOther {
			t.Errorf("AdminReservationProcess handler (%s) returned code: %d, expected code: %d", v.name, rr.Code, v.expectedStatusCode)
		}
	}
}

// ========================================================================================================================================

var adminReservationDeleteTests = []struct {
	name, queryParams, expectedLocation string
	expectedStatusCode                  int
}{
	{
		name:               "reservation-deleted",
		queryParams:        "",
		expectedLocation:   "",
		expectedStatusCode: http.StatusSeeOther,
	},
	{
		name:               "reservation-deleted-cal",
		queryParams:        "?y=2022&m=10",
		expectedLocation:   "",
		expectedStatusCode: http.StatusSeeOther,
	},
}

func TestRepository_AdminReservationDelete(t *testing.T) {
	for _, v := range adminReservationDeleteTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/reservation-deleted/cal/1/page%s", v.queryParams), nil)
		ctx := getCtx(req)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AdminReservationDelete)
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusSeeOther {
			t.Errorf("AdminReservationDelete handler (%s) returned code: %d, expected code: %d", v.name, rr.Code, v.expectedStatusCode)
		}
	}
}

// ========================================================================================================================================

// getCtx creates a context for use in TestRepository_Reservation() request
func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}

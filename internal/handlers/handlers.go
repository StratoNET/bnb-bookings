package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/database"
	"github.com/StratoNET/bnb-bookings/internal/helpers"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/StratoNET/bnb-bookings/internal/render"
	"github.com/StratoNET/bnb-bookings/internal/repository"
	"github.com/StratoNET/bnb-bookings/internal/repository/dbrepository"
	forms "github.com/StratoNET/bnb-bookings/internal/validation"
	"github.com/go-chi/chi/v5"
)

// Repo repository used by handlers
var Repo *Repository

// Repository, which incorporates a database repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepository
}

// NewRepository creates a new repository, which incorporates a database repository
func NewRepository(a *config.AppConfig, db *database.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepository.NewMariaDBRepository(db.SQL, a),
	}
}

// NewHandlers sets repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Index is the handler for the home page
func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// GQ is the handler for the General's Quarters (gq) page
func (m *Repository) GQ(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "gq.page.tmpl", &models.TemplateData{})
}

// MS is the handler for the Major's Suite (ms) page
func (m *Repository) MS(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "ms.page.tmpl", &models.TemplateData{})
}

// Availability is the handler for the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles request for search-availability page data
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	// parse dates as appropriate, format required is dd/mm/yyyy -- (Go format reminder is 01/02 03:04:05PM '06 -0700)
	layout := "02/01/2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "Sorry, no availability for the specific requested period)")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	// by this point there must be availability, utilise data object available throughout from templatedata
	data := make(map[string]interface{})
	data["rooms"] = rooms

	// instantiate a reservation with only the information known so far from search availability (i.e. the dates, room is still unknown)
	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	// store this within session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// using predefined Data object from templatedata struct, pass in data
	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// PostAvailabilityModal handles request for modal availability search & returns JSON response
func (m *Repository) PostAvailabilityModal(w http.ResponseWriter, r *http.Request) {
	rsp := jsonResponse{
		Ok:      false,
		Message: "Available !",
	}
	out, err := json.MarshalIndent(rsp, "", "    ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation renders the make-reservation page & displays associated form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// update the reservation data further, currently stored with only start/end dates & also room id (since added by ChooseRoom() )
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation dates & room number from session"))
		return
	}

	// get room name & populate Room, which is a member of Reservation model (only RoomName is required)
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation.Room.RoomName = room.RoomName

	// format start/end dates as strings (instead of time.Time) & place into a StringMap (templatedata) for display in make-reservation page
	sd := reservation.StartDate.Format("02/01/2006")
	ed := reservation.EndDate.Format("02/01/2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	// put updated (with Room.RoomName) reservation back into session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// provide access to template data's (initially empty) Form object first time this page is rendered
		Form:      forms.NewForm(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation is the handler for posting the make-reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// update the reservation data further, currently stored with only start/end dates, room id & room name
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation dates, room number & room name from session"))
		return
	}

	// initially ensure form data is parsed correctly
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.NewForm(r.PostForm)

	// perform all necessary validations

	// form.HasField("first_name", r)
	form.RequiredFields("first_name", "last_name", "email")
	form.MinLength("first_name", 2)
	form.IsEmail("email")

	if !form.ValidForm() {
		// format start/end dates as strings (instead of time.Time) & place into a StringMap (templatedata) for re-display in make-reservation page
		sd := reservation.StartDate.Format("02/01/2006")
		ed := reservation.EndDate.Format("02/01/2006")
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	// after all validation procedures are complete, insert reservation record into database
	lastReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// populate and instance of the RoomRestriction object, get ReservationID from LastInsertId() after executing InsertReservation()
	room_restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		ReservationID: lastReservationID,
		RestrictionID: 1,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
	}

	err = m.DB.InsertRoomRestriction(room_restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//for a valid reservation form, put into session & redirect to summary page
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary is the handler for displaying reservation details
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get reservation from session which requires type assertion/casting, this sets ok true (or false on failure)
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "There are no reservation details available to display")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//reaching this point implies 'reservation' was successfully retrieved, therefore can now be removed from session
	m.App.Session.Remove(r.Context(), "reservation")

	// format start/end dates as strings (instead of time.Time) & place into a StringMap (templatedata) for display in reservation-summary page
	sd := reservation.StartDate.Format("Monday 02 January 2006")
	ed := reservation.EndDate.Format("Monday 02 January 2006")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	// create data object and populate with reservation data
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// get id from room link clicked
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// update the reservation data (currently stored with only start/end dates in session), with the required room id
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}
	// add roomID to reservation model, put it back into the session & redirect to make-reservation page
	reservation.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

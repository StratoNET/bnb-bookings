package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// NewTestRepository creates a new testing repository, which does NOT incorporate a db repository because db access NOT needed for unit tests
func NewTestRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepository.NewTestingDBRepository(a),
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
	// initially ensure form data is parsed correctly
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0003: cannot parse search availability, all rooms, form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	// parse dates as appropriate, format required is dd/mm/yyyy -- (Go format reminder is 01/02 03:04:05PM '06 -0700)
	layout := "02/01/2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "unable to parse START DATE")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "unable to parse END DATE")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "unable to get AVAILABILITY FOR ALL ROOMS")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "Sorry, no availability for the specific requested period")
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
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// PostAvailabilityModal handles request for modal availability search & returns JSON response
func (m *Repository) PostAvailabilityModal(w http.ResponseWriter, r *http.Request) {
	// need to parse form request body, good practice & allows for testing
	err := r.ParseForm()
	if err != nil {
		// i.e. cannot parse form, so return an appropriate JSON response
		resp := jsonResponse{
			Ok:      false,
			Message: "Internal Server Error",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")

	// parse dates as appropriate, format required is dd/mm/yyyy -- (Go format reminder is 01/02 03:04:05PM '06 -0700)
	layout := "02/01/2006"
	startDate, _ := time.Parse(layout, start)
	endDate, _ := time.Parse(layout, end)
	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesAndRoomID(startDate, endDate, roomID)
	if err != nil {
		// return an appropriate JSON response
		resp := jsonResponse{
			Ok:      false,
			Message: "#0011: error connecting to database during dates & room number search",
		}
		out, _ := json.MarshalIndent(resp, "", "    ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		Ok:        available,
		Message:   "",
		RoomID:    strconv.Itoa(roomID),
		StartDate: start,
		EndDate:   end,
	}

	// ignore error check at this point because all aspects of JSON response have already been handled
	out, _ := json.MarshalIndent(resp, "", "    ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Reservation renders the make-reservation page & displays associated form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// update the reservation data further, currently stored with only start/end dates & also room id (since added by ChooseRoom() )
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "#0001: cannot get reservation dates & room number from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// get room name & populate Room, which is a member of Reservation model (only RoomName is required)
	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0002: cannot get room number from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "#0004: cannot get reservation dates, room number & room name from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// initially ensure form data is parsed correctly
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0005: cannot parse reservation form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.NewForm(r.PostForm)

	// perform all necessary validations

	form.RequiredFields("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 2)
	form.MinLength("last_name", 2)
	form.MinLength("phone", 6)
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

		m.App.Session.Put(r.Context(), "error", "#0006: invalid form details submitted")

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
		m.App.Session.Put(r.Context(), "error", "#0007: cannot insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
		m.App.Session.Put(r.Context(), "error", "#0008: cannot insert 'room restriction', corresponding to current reservation, into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send email notification to guest
	htmlMsg := fmt.Sprintf(`
		<h3 class="text-center">Eden House: Reservation Confirmation</h3>
		<p>Dear %s&nbsp;%s</p>
		<p>This is to confirm your reservation from %s to %s in the %s, we look forward to seeing you then.</p>
	`, reservation.FirstName, reservation.LastName, reservation.StartDate.Format("Monday 02 January 2006"),
		reservation.EndDate.Format("Monday 02 January 2006"), reservation.Room.RoomName)
	msg := models.MailData{
		To:       reservation.Email,
		From:     "reservations@edenhouse.com",
		Subject:  "Eden House: Reservation Confirmation",
		Content:  htmlMsg,
		Template: "basic.html",
	}
	m.App.MailChannel <- msg

	// send email notification to owner / admin
	htmlMsg = fmt.Sprintf(`
		<h3>Eden House: Reservation Notification</h3>
		<p>A reservation has been made by %s %s covering %s to %s for the %s.</p>
	`, reservation.FirstName, reservation.LastName, reservation.StartDate.Format("02/01/2006"), reservation.EndDate.Format("02/01/2006"),
		reservation.Room.RoomName)
	msg = models.MailData{
		To:      "reservations@edenhouse.com",
		From:    "reservations@edenhouse.com",
		Subject: "Eden House: Reservation Notification",
		Content: htmlMsg,
	}
	m.App.MailChannel <- msg

	//for a valid reservation form, put into session & redirect to summary page
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary is the handler for displaying reservation details
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get reservation from session which requires type assertion/casting, this sets ok true (or false on failure)
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
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

// ChooseRoom puts room ID into session after clicking on available rooms in clickable list
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// get id from room link clicked, get elements exploded on '/' & convert 3rd element into string
	elements := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(elements[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0009: missing URL parameter (id)")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// update the reservation data (currently stored with only start/end dates in session), with the required room id
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "#0010: unable to retrieve reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

// ReserveRoom takes URL parameters from modal dialog, puts them into session & redirects to make-reservation
func (m *Repository) ReserveRoom(w http.ResponseWriter, r *http.Request) {
	// retrieve parameters from url get request
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	sd := r.URL.Query().Get("sd")
	ed := r.URL.Query().Get("ed")

	// parse dates as appropriate, format required is dd/mm/yyyy -- (Go format reminder is 01/02 03:04:05PM '06 -0700)
	layout := "02/01/2006"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	// instantiate a reservation
	var reservation models.Reservation

	// get room name & populate Room, which is a member of Reservation model (only RoomName is required)
	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0012: cannot get room name from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// populate reservation with currently known details
	reservation.RoomID = roomID
	reservation.Room.RoomName = room.RoomName
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	// put all details back into the session & redirect to make-reservation page
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// Login displays login page & gets the administrator's login form
func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.NewForm(nil),
	})
}

//PostLogin handles administrator login
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// prevent session fixation attack
	_ = m.App.Session.RenewToken(r.Context())
	// parse form
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	// validation...
	form := forms.NewForm(r.PostForm)
	form.RequiredFields("email", "password")
	form.IsEmail("email")
	if !form.ValidForm() {
		stringMap := make(map[string]string)
		stringMap["email"] = email
		stringMap["password"] = password
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form:      form,
			StringMap: stringMap,
		})
		return
	}
	// attempt administrator authentication
	id, _, err := m.DB.AuthenticateAdministrator(email, password)
	if err != nil {
		m.App.ErrorLog.Println(err)
		m.App.Session.Put(r.Context(), "error", "Please check - invalid login credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// store the returned administrator's id in session
	m.App.Session.Put(r.Context(), "admin_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs an administrator out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// AdminDashboard gets the reservations management dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminReservationsNew gets any new (unprocessed) reservations
func (m *Repository) AdminReservationsNew(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetNewReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0013: cannot get new reservations from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-reservations-new.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminReservationsAll gets all reservations
func (m *Repository) AdminReservationsAll(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0014: cannot get all reservations from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-reservations-all.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminReservationsCalendar gets the reservations displayed in calendar block format
func (m *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// initially assume that request URL has no added parameters i.e. no year/month specified, which is equal to now...
	now := time.Now()
	// ... if URL parameter are included
	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		// allows updating of 'now' to incorporate parameter values for next/prev(ious) dates
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	}
	// allows next/prev to be created based on +/- one month from 'now'
	next := now.AddDate(0, 1, 0)
	prev := now.AddDate(0, -1, 0)
	// format for display in reservations calendar page, remembering that a month can roll from 12 to 1 (or vice-versa) if year changes
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	prevMonth := prev.Format("01")
	prevMonthYear := prev.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["prev_month"] = prevMonth
	stringMap["prev_month_year"] = prevMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")
	stringMap["show_month"] = now.Format("January")

	// need to get first & last days of each month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)
	// store number of days in month in intMap which will vary based on current result from preceding code block
	dimIntMap := make(map[string]int)
	dimIntMap["days_in_month"] = lastDayOfMonth.Day()

	// get information for all rooms
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// pass all rooms data to calendar template
	data := make(map[string]interface{})
	data["rooms"] = rooms

	// create reservation & owner blocked maps for each day of current month, for each room
	for _, rm := range rooms {
		reservationsMap := make(map[string]int)
		blockedMap := make(map[string]int)

		for d := firstDayOfMonth; !d.After(lastDayOfMonth); d = d.AddDate(0, 0, 1) {
			// initialise each day in each map as 0
			reservationsMap[d.Format("2-01-2006")] = 0
			blockedMap[d.Format("2-01-2006")] = 0
		}
		// get all room restrictions for each room
		restrictions, err := m.DB.GetRoomRestrictionsByDate(rm.ID, firstDayOfMonth, lastDayOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, rs := range restrictions {
			if rs.ReservationID == 0 {
				// owner blocked entries where query has substituted null with 0 for reservation_id
				for d := rs.StartDate; !d.After(rs.EndDate); d = d.AddDate(0, 0, 1) {
					blockedMap[d.Format("2-01-2006")] = rs.ID
				}
			} else {
				// reservations
				for d := rs.StartDate; !d.After(rs.EndDate); d = d.AddDate(0, 0, 1) {
					reservationsMap[d.Format("2-01-2006")] = int(rs.ReservationID)
				}
			}
		}
		data[fmt.Sprintf("reservations_map_%d", rm.ID)] = reservationsMap
		data[fmt.Sprintf("blocked_map_%d", rm.ID)] = blockedMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("blocked_map_%d", rm.ID), blockedMap)

	}

	render.Template(w, r, "admin-reservations-cal.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		IntMap:    dimIntMap,
		Data:      data,
	})
}

// AdminPostReservationsCalendar handles posting of edited data submitted from the reservations calendar
func (m *Repository) AdminPostReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

	month, _ := strconv.Atoi(r.Form.Get("month"))
	year, _ := strconv.Atoi(r.Form.Get("year"))
	blocks, _ := strconv.Atoi(r.Form.Get("blocks"))

	// process owner room blocks

	// get information for all rooms
	rooms, err := m.DB.GetAllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get access to forms.go HasField()
	form := forms.NewForm(r.PostForm)

	// 1. handle owner blocks to be REMOVED for current month, for each room
	for _, rm := range rooms {
		// get blocked maps from session for given rooms / month, which contain current blocked information at point calendar was displayed.
		// loop through entire map for each room. If there is an entry in a map which does NOT exist in posted data AND if its restriction_id
		// is > 0, then it is a room block which needs to be removed. This is because unchecking a checkbox means that it will never be submitted
		// within post data, which is normal behaviour. Consequently if it was previously checked in the map, it will have a value > 0 and need
		// to be removed.
		currentBlockedMap := m.App.Session.Get(r.Context(), fmt.Sprintf("blocked_map_%d", rm.ID)).(map[string]int)

		for dated, rsID := range currentBlockedMap {
			// ok will be false if value is NOT in the map
			if value, ok := currentBlockedMap[dated]; ok {
				// only concerned with values > 0 AND which are NOT in form post data (i.e. unchecked), the remainder amount to placeholders for days without /// blocks
				if value > 0 {
					if !form.HasField(fmt.Sprintf("remove_blocked_%d_%s", rm.ID, dated)) {
						// delete the blocked restriction by its id
						err := m.DB.DeleteRoomBlock(rsID)
						if err != nil {
							m.App.ErrorLog.Println(err)
						}
					}
				}
			}
		}
	}

	// 2. handle owner blocks to be ADDED for current month, for each room
	for blockedName := range r.PostForm {
		if strings.HasPrefix(blockedName, "add_blocked") {
			elements := strings.Split(blockedName, "_")
			roomID, _ := strconv.Atoi(elements[2])
			startDate, _ := time.Parse("2-01-2006", elements[3])
			endDate := startDate.AddDate(0, 0, (blocks - 1))
			// insert a new owner block
			err := m.DB.InsertRoomBlock(roomID, startDate, endDate)
			if err != nil {
				m.App.ErrorLog.Println(err)
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Room(s) blocked... changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-cal?y=%d&m=%d", year, month), http.StatusSeeOther)
}

// AdminReservationProcess marks the associated reservation as processed
func (m *Repository) AdminReservationProcess(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.ErrorLog.Println(" missing URL parameter (id)")
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.UpdateReservationProcessed(id, 1)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0019: cannot update requested reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusTemporaryRedirect)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Reservation (id=%d) has been marked as processed", id))

	if (year == "") && (month == "") {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-cal?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminReservationDelete deletes the associated reservation from the database
func (m *Repository) AdminReservationDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		m.App.ErrorLog.Println(" missing URL parameter (id)")
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.DeleteReservation(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0020: cannot delete requested reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusTemporaryRedirect)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Reservation (id=%d) has been deleted", id))

	if (year == "") && (month == "") {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-cal?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminReservation gets a single reservation by id & displays it in form layout
func (m *Repository) AdminReservation(w http.ResponseWriter, r *http.Request) {
	// get id from reservation link clicked, get elements exploded on '/' & convert 4th element into string
	elements := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(elements[4])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0015: missing URL parameter (id)")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}
	stringMap := make(map[string]string)
	src := elements[3]
	stringMap["src"] = src
	stringMap["year"] = r.URL.Query().Get("y")
	stringMap["month"] = r.URL.Query().Get("m")

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0016: cannot get requested reservation from database")
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "admin-reservation.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
		Form:      forms.NewForm(nil),
	})
}

// AdminPostReservation updates a single reservation by id
func (m *Repository) AdminPostReservation(w http.ResponseWriter, r *http.Request) {
	// initially ensure form data is parsed correctly
	err := r.ParseForm()
	if err != nil {
		m.App.ErrorLog.Println("unable to parse update reservation form")
		return
	}
	// get id from reservation link clicked, get elements exploded on '/' & convert 4th element into string
	elements := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(elements[4])
	if err != nil {
		m.App.ErrorLog.Println(" missing URL parameter (id)")
		return
	}
	stringMap := make(map[string]string)
	src := elements[3]
	stringMap["src"] = src

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0017: cannot get requested reservation from database")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusTemporaryRedirect)
		return
	}
	// overwrite any/all of 4 possible updatable fields
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.NewForm(r.PostForm)

	// perform all necessary validations

	// form.HasField("first_name", r)
	form.RequiredFields("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 2)
	form.MinLength("last_name", 2)
	form.MinLength("phone", 6)
	form.IsEmail("email")

	if !form.ValidForm() {
		// format start/end dates as strings (instead of time.Time) & place into a StringMap (templatedata) for re-display in make-reservation page
		sd := reservation.StartDate.Format("02/01/2006")
		ed := reservation.EndDate.Format("02/01/2006")
		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed
		// also must persist 'src' to maintain URL integrity if validation fails and function is subsequently recalled
		stringMap["src"] = src

		data := make(map[string]interface{})
		data["reservation"] = reservation

		m.App.Session.Put(r.Context(), "error", "#0021: invalid form details submitted")

		render.Template(w, r, "admin-reservation.page.tmpl", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	// after all validation procedures are complete, update reservation record in database
	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "#0018: cannot update reservation")
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusTemporaryRedirect)
		return
	}

	year := r.Form.Get("y")
	month := r.Form.Get("m")

	m.App.Session.Put(r.Context(), "flash", fmt.Sprintf("Reservation (id=%d) has been updated", reservation.ID))

	if (year == "") && (month == "") {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-cal?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

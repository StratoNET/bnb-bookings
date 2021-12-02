package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/helpers"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/StratoNET/bnb-bookings/internal/render"
	forms "github.com/StratoNET/bnb-bookings/internal/validation"
)

// Repo repository used by handlers
var Repo *Repository

// Repository
type Repository struct {
	App *config.AppConfig
}

// NewRepository creates a new repository
func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Index is the handler for the home page
func (m *Repository) Index(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "index.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

// GQ is the handler for the General's Quarters (gq) page
func (m *Repository) GQ(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "gq.page.tmpl", &models.TemplateData{})
}

// MS is the handler for the Major's Suite (ms) page
func (m *Repository) MS(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "ms.page.tmpl", &models.TemplateData{})
}

// Availability is the handler for the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability handles request for search-availability page data
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	w.Write([]byte(fmt.Sprintf("Start date is %s and end date is %s.", start, end)))
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
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// provide access to template data's (initially empty) Form object first time this page is rendered
		Form: forms.NewForm(nil),
		Data: data,
	})
}

// PostReservation is the handler for posting the make-reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// initially ensure form data is parsed correctly
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	// populate an instance of reservation object with data user has entered, even if 'bad' data
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.NewForm(r.PostForm)

	// perform all necessary validations

	// form.HasField("first_name", r)
	form.RequiredFields("first_name", "last_name", "email")
	form.MinLength("first_name", 2)
	form.IsEmail("email")

	if !form.ValidForm() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
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

	// create data object and populate with reservation data
	data := make(map[string]interface{})
	data["reservation"] = reservation

	// using predefined Data object from templatedata struct, pass in data
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

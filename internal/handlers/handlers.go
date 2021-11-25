package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/StratoNET/bnb-bookings/internal/config"
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
		log.Println(err)
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
		log.Println(err)
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
	form.MinLength("first_name", 2, r)
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
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

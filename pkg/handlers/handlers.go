package handlers

import (
	"net/http"

	"github.com/StratoNET/bnb-bookings/pkg/config"
	"github.com/StratoNET/bnb-bookings/pkg/models"
	"github.com/StratoNET/bnb-bookings/pkg/render"
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
	remoteIP := r.RemoteAddr
	// put remote IP address into session data via App Repository instance
	m.App.Session.Put(r.Context(), "remote_IP", remoteIP)

	render.RenderTemplate(w, "index.page.tmpl", &models.TemplateData{})
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic to produce passable data
	stringMap := make(map[string]string)

	// get remote IP, via key, stored in session (during Home handler request),
	// this could be "" if this handler request is made first
	remoteIP := m.App.Session.GetString(r.Context(), "remote_IP")

	// add to stringMap
	stringMap["remote_IP"] = remoteIP

	// pass data to template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// GQ is the handler for the General's Quarters (gq) page
func (m *Repository) GQ(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "gq.page.tmpl", &models.TemplateData{})
}

// MS is the handler for the Major's Suite (ms) page
func (m *Repository) MS(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "ms.page.tmpl", &models.TemplateData{})
}

// Availability is the handler for the search-availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

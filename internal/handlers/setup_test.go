package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/StratoNET/bnb-bookings/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	// session needs to be informed for storage of complex types, in this case...
	gob.Register(models.Reservation{})

	// set development / production mode
	app.ProductionMode = false

	// create InfoLog & ErrorLog, making them available throughout application via config
	infoLog := log.New(os.Stdout, "\033[36;1mINFO\033[0;0m\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "\033[31;1mERROR\033[0;0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// invoke session management via 'scs' package
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.ProductionMode
	// apply session parameters throughout application
	app.Session = session

	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create application template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true

	repo := NewTestRepository(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	os.Exit(m.Run())
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()

	// invoke middleware requirements prior to routing
	//===========================
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	//===========================

	mux.Get("/", Repo.Index)
	mux.Get("/about", Repo.About)
	mux.Get("/gq", Repo.GQ)
	mux.Get("/ms", Repo.MS)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-modal", Repo.PostAvailabilityModal)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/contact", Repo.Contact)

	// creat fileserver for static content
	staticFileServer := http.FileServer(http.Dir("./static/"))
	// handle any file within static sub-folders
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))

	return mux
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.ProductionMode,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads & saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts

	}

	return myCache, nil
}

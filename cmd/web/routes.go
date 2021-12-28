package main

import (
	"net/http"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	// invoke middleware requirements prior to routing
	//===========================
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	//===========================

	mux.Get("/", handlers.Repo.Index)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/gq", handlers.Repo.GQ)
	mux.Get("/ms", handlers.Repo.MS)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/reserve-room", handlers.Repo.ReserveRoom)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-modal", handlers.Repo.PostAvailabilityModal)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/login", handlers.Repo.Login)
	mux.Get("/logout", handlers.Repo.Logout)
	mux.Post("/login", handlers.Repo.PostLogin)

	// provides access to any protected routes prepended with "/admin" & uses middleware Auth() to determine authentication
	mux.Route("/admin", func(mux chi.Router) {
		if app.ProductionMode {
			mux.Use(Auth)
		}
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new", handlers.Repo.AdminReservationsNew)
		mux.Get("/reservations-all", handlers.Repo.AdminReservationsAll)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)
		// these routes can be reached via either 'all' or 'new' reservations administration pages
		mux.Get("/reservation-processed/{src}/{id}", handlers.Repo.AdminReservationProcess)
		mux.Get("/reservation-deleted/{src}/{id}", handlers.Repo.AdminReservationDelete)
		mux.Get("/reservations/{src}/{id}", handlers.Repo.AdminReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostReservation)
	})

	// creat fileserver for static content
	staticFileServer := http.FileServer(http.Dir("./static/"))
	// handle any file within static sub-folders
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))

	return mux
}

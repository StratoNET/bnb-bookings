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

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-modal", handlers.Repo.PostAvailabilityModal)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	mux.Get("/contact", handlers.Repo.Contact)

	// creat fileserver for static content
	staticFileServer := http.FileServer(http.Dir("./static/"))
	// handle any file within static sub-folders
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))

	return mux
}

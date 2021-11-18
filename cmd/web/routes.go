package main

import (
	"net/http"

	"github.com/StratoNET/bnb-bookings/pkg/config"
	"github.com/StratoNET/bnb-bookings/pkg/handlers"
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

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// creat fileserver for static content
	staticFileServer := http.FileServer(http.Dir("./static/"))
	// handle any file within static sub-folders
	mux.Handle("/static/*", http.StripPrefix("/static", staticFileServer))

	return mux
}

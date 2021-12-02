package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/StratoNET/bnb-bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for handling helper functions
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.ErrorLog.Println("Client error: status", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	// print error & stack trace...
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	// then supply feedback to user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

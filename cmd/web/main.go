package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/database"
	"github.com/StratoNET/bnb-bookings/internal/handlers"
	"github.com/StratoNET/bnb-bookings/internal/helpers"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/StratoNET/bnb-bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main application function
func main() {

	db, err := run_main()
	if err != nil {
		log.Fatal(err)
	}

	// only allow the database connection to close after main() executes fully
	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

// initial main() code function for testing
func run_main() (*database.DB, error) {

	// session needs to be informed for storage of complex types, in this case...
	gob.Register(models.Administrator{})
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.RestrictionCategory{})

	// set development / production mode
	app.ProductionMode = false

	// create InfoLog & ErrorLog, making them available throughout application via config
	infoLog = log.New(os.Stdout, "\033[36;1mINFO\033[0;0m\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "\033[31;1mERROR\033[0;0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// invoke session management via 'scs' package
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.ProductionMode
	// apply session parameters throughout application
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := database.ConnectSQL("root:@tcp(localhost:3306)/bnb-bookings")
	if err != nil {
		log.Fatal("Cannot connect to database ! ... terminating...")
	}
	log.Println("Connected to database OK")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create application template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepository(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}

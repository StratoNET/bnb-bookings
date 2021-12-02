package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// session needs to be informed for storage of complex types, in this case...
	gob.Register(models.Reservation{})

	// set development / production mode
	testApp.ProductionMode = false

	// create InfoLog & ErrorLog, making them available throughout application via config
	infoLog := log.New(os.Stdout, "\033[36;1mINFO\033[0;0m\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "\033[31;1mERROR\033[0;0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// invoke session management via 'scs' package
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.ProductionMode
	// apply session parameters throughout application
	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

// invoke an http.ResponseWriter providing a (1)header, (2)write & (3)write header for testing
type myWriter struct{}

// 1
func (mw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

// 2
func (mw *myWriter) Write(b []byte) (int, error) {
	l := len(b)
	return l, nil
}

// 3
func (mw *myWriter) WriteHeader(i int) {

}

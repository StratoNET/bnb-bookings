package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/StratoNET/bnb-bookings/internal/config"
	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/justinas/nosurf"
)

// functions / FuncMap makes utility functions available throughout application templates
var functions = template.FuncMap{
	// "add":         Add,
	"dateUK":      DateUK,
	"iterateDays": IterateDays,
}

var app *config.AppConfig

var pathToTemplates = "./templates"

// Add returns sum of 2 integers (currently unused)
// func Add(x, y int) int {
// 	return x + y
// }

// DateTimeUK returns a date & time formatted for UK
func DateUK(t time.Time) string {
	return t.Format("02/01/2006")
}

// IterateDays returns a []int, starting at 1 up to value of count for every 'days_in_month' stored in calendar's IntMap
func IterateDays(count int) []int {
	var i int
	var days []int
	for i = 1; i <= count; i++ {
		days = append(days, i)
	}
	return days
}

// NewRenderer sets config for template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// add any default data, required on every page template, at this point
	td.CSRFToken = nosurf.Token(r)
	td.RemoteIP = r.RemoteAddr
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	// determine if session user is authenticated by checking for their id
	if app.Session.Exists(r.Context(), "admin_id") {
		td.IsAuthenticated = true
	}
	return td
}

// Template renders templates using html/template (modified to return any potential error)
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from app config (PRODUCTION mode)
		tc = app.TemplateCache
	} else {
		// (DEVELOPMENT mode)
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		return errors.New("unable to retrieve template from template cache")
	}

	buf := new(bytes.Buffer)

	// incorporate any default template data before executing/writing template
	td = AddDefaultData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("error rendering template in browser:", err)
		return err
	}

	return nil

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
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

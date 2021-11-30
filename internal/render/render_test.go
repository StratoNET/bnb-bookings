package render

import (
	"net/http"
	"testing"

	"github.com/StratoNET/bnb-bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "saviour of the universe")

	result := AddDefaultData(&td, r)
	if result.Flash != "saviour of the universe" {
		t.Error(" flash value of \"saviour of the universe\" not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	// call this method to get an http.Request object
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	//	app.UseCache = true

	err = RenderTemplate(&ww, r, "index.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error rendering template in browser")
	}

	err = RenderTemplate(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("has rendered a template that does not exist !")
	}

}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

// getSession provides necessary session data for TestAddDefaultData()
func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/any-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil

}

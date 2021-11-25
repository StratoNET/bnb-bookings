package models

import (
	forms "github.com/StratoNET/bnb-bookings/internal/validation"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	RemoteIP  string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}
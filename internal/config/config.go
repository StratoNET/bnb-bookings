package config

import (
	"html/template"
	"log"

	"github.com/StratoNET/bnb-bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

// AppConfig holds application configuration
type AppConfig struct {
	UseCache       bool
	TemplateCache  map[string]*template.Template
	Session        *scs.SessionManager
	InfoLog        *log.Logger
	ErrorLog       *log.Logger
	ProductionMode bool
	MailChannel    chan models.MailData
}

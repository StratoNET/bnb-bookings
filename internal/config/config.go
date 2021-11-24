package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds application configuration
type AppConfig struct {
	UseCache       bool
	TemplateCache  map[string]*template.Template
	Session        *scs.SessionManager
	InfoLog        *log.Logger
	ProductionMode bool
}

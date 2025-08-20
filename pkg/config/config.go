package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// The config package is imported by other parts of the application
// It does not import anything from other parts of the application

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}

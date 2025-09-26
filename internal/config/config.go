package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

// The config package is imported by other parts of the application
// It does not import anything from other parts of the application

// DateLayout is the format we expect dates to be sent in as
const DateLayout = "2006-01-02" // 01/02 03:04:05PM '06 -0700

// AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}

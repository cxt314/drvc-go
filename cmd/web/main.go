package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cxt314/drvc-go/pkg/config"
	"github.com/cxt314/drvc-go/pkg/handlers"
	"github.com/cxt314/drvc-go/pkg/render"
)

const portNumber = ":8080"

// main is hte main application function
func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewTemplates(&app)

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Starting application on port %s\n", portNumber)
	_ = http.ListenAndServe(portNumber, nil)
}

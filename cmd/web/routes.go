package main

import (
	"net/http"

	"github.com/cxt314/drvc-go/internal/config"
	"github.com/cxt314/drvc-go/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// middleware
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/vehicles", handlers.Repo.VehicleList)
	mux.Get("/new-vehicle", handlers.Repo.VehicleCreate)
	mux.Post("/new-vehicle", handlers.Repo.VehicleCreatePost)
	mux.Get("/vehicles/{id}", handlers.Repo.VehicleEdit)
	mux.Post("/vehicles/{id}", handlers.Repo.VehicleEditPost)

	// sample reservation routes
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// create a fileserver for serving static files
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

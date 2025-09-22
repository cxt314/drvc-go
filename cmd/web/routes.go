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

	// vehicles handlers
	mux.Get("/vehicles", handlers.Repo.VehicleList)
	mux.Get("/new-vehicle", handlers.Repo.VehicleCreate)
	mux.Post("/new-vehicle", handlers.Repo.VehicleCreatePost)
	mux.Get("/vehicles/{id}", handlers.Repo.VehicleEdit)
	mux.Post("/vehicles/{id}", handlers.Repo.VehicleEditPost)
	//mux.Get("/vehicles/{id}/delete", handlers.Repo.VehicleDelete)
	mux.Get("/vehicles/{id}/deactivate", handlers.Repo.VehicleDeactivate)

	// members handlers
	mux.Get("/members", handlers.Repo.MemberList)
	mux.Get("/new-member", handlers.Repo.MemberCreate)
	mux.Post("/new-member", handlers.Repo.MemberCreatePost)
	mux.Get("/members/{id}", handlers.Repo.MemberEdit)
	mux.Post("/members/{id}", handlers.Repo.MemberEditPost)
	//mux.Get("/members/{id}/delete", handlers.Repo.MemberDelete)
	mux.Get("/members/{id}/deactivate", handlers.Repo.MemberDeactivate)

	// mileage logs handlers
	mux.Get("/mileage-logs", handlers.Repo.MileageLogList)
	mux.Get("/mileage-logs/list/{id}", handlers.Repo.MileageLogListByVehicle)
	mux.Get("/new-mileage-log", handlers.Repo.MileageLogCreate)
	mux.Post("/new-mileage-log", handlers.Repo.MileageLogCreatePost)
	mux.Get("/mileage-logs/{id}", handlers.Repo.MileageLogEdit)
	mux.Post("/mileage-logs/{id}", handlers.Repo.MileageLogEditPost)
	mux.Get("/mileage-logs/{id}/delete", handlers.Repo.MileageLogDelete)
	mux.Get("/mileage-logs/{id}/edit-trips", handlers.Repo.TripsEdit)
	mux.Post("/mileage-logs/{id}/edit-trips", handlers.Repo.TripsEditPost)
	mux.Post("/mileage-logs/{id}/add-trip", handlers.Repo.AddTripPost)

	// htmx routes
	mux.Get("/remove-item", handlers.Repo.RemoveItem)
	mux.Get("/members/add-alias", handlers.Repo.AddAlias)

	// sample reservation routes
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// create a fileserver for serving static files
	fileServer := http.FileServer(http.Dir("./static/"))

	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

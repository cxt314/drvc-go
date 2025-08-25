package handlers

import (
	"net/http"
	"strconv"

	"github.com/cxt314/drvc-go/internal/config"
	"github.com/cxt314/drvc-go/internal/driver"
	"github.com/cxt314/drvc-go/internal/forms"
	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
	"github.com/cxt314/drvc-go/internal/repository"
	"github.com/cxt314/drvc-go/internal/repository/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (m *Repository) VehicleList(w http.ResponseWriter, r *http.Request) {
	vehicles, err := m.DB.AllVehicles()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["vehicles"] = vehicles

	render.Template(w, r, "vehicle-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) VehicleCreate(w http.ResponseWriter, r *http.Request) {
	var emptyModel models.Vehicle
	data := make(map[string]interface{})
	data["vehicle"] = emptyModel

	render.Template(w, r, "edit-vehicle.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) VehicleCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	v := models.Vehicle{
		Name:         r.Form.Get("name"),
		Make:         r.Form.Get("make"),
		Model:        r.Form.Get("model"),
		FuelType:     r.Form.Get("fuel_type"),
		Vin:          r.Form.Get("vin"),
		LicensePlate: r.Form.Get("license_plate"),
	}
	v.Year, err = strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		m.App.InfoLog.Println("Could not convert year to int,", r.Form.Get("year"))
	}

	//v.PurchaseDate
	v.PurchasePrice = models.StrToUSD(r.Form.Get("purchase_price"))

	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["vehicle"] = v

		render.Template(w, r, "edit-vehicle.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertVehicle(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/vehicles", http.StatusSeeOther)
}

// sample form filling

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservation is a sample function for how to handle form parsing & validation
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

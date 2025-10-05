package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/cxt314/drvc-go/internal/forms"
	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

// VehicleList displays a list of all vehicles
func (m *Repository) VehicleList(w http.ResponseWriter, r *http.Request) {
	vehicles, err := m.DB.GetVehicleByActive(true)
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

// VehicleCreate displays the page to create a new vehicle
func (m *Repository) VehicleCreate(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	data["fuelTypes"] = models.FuelTypes
	data["billingTypes"] = models.BillingTypes

	render.Template(w, r, "edit-vehicle.page.tmpl", &models.TemplateData{
		Data: data,
		Form: forms.New(nil),
	})
}

// VehicleCreatePost processes the POST request for creating a new vehicle
func (m *Repository) VehicleCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse received form values into vehicle object
	v := models.Vehicle{}
	err := helpers.ParseFormToVehicle(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["vehicle"] = v
		data["fuelTypes"] = models.FuelTypes
		data["billingTypes"] = models.BillingTypes

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

// VehicleEdit shows the edit form for a vehicle by id
func (m *Repository) VehicleEdit(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get vehicle from database
	v, err := m.DB.GetVehicleByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["vehicle"] = v
	data["fuelTypes"] = models.FuelTypes
	data["billingTypes"] = models.BillingTypes

	render.Template(w, r, "edit-vehicle.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// VehicleEditPost processes the POST request for editing a vehicle by id
func (m *Repository) VehicleEditPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get vehicle from database
	v, err := m.DB.GetVehicleByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form into fetched vehicle
	err = helpers.ParseFormToVehicle(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["vehicle"] = v
		data["fuelTypes"] = models.FuelTypes
		data["billingTypes"] = models.BillingTypes

		render.Template(w, r, "edit-vehicle.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateVehicle(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Updated vehicle successfully")
	http.Redirect(w, r, fmt.Sprintf("/vehicles/%d", id), http.StatusSeeOther)
}

// VehicleDeactivate updates a vehicle's active status by id
func (m *Repository) VehicleDeactivate(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// update vehicle set is_active to false
	err = m.DB.UpdateVehicleActiveByID(id, false)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// redirect to vehicles list after making vehicle inactive
	http.Redirect(w, r, "/vehicles", http.StatusSeeOther)
}

/*// VehicleDelete deletes the vehicle with the given id
func (m *Repository) VehicleDelete(w http.ResponseWriter, r *http.Request) {
	// exploded := strings.Split(r.RequestURI, "/")
	// id, err := strconv.Atoi(exploded[2])
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	render.Template(w, r, "vehicle-list.page.tmpl", &models.TemplateData{})
}*/

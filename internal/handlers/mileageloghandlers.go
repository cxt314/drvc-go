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

func (m *Repository) MileageLogList(w http.ResponseWriter, r *http.Request) {
	// get vehicles
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["vehicles"] = vehicles

	render.Template(w, r, "mileage-log-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) MileageLogListByVehicle(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})

	// get vehicles
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["vehicles"] = vehicles

	// get active vehicle
	selectedVehicle, err := m.DB.GetVehicleByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["selected-vehicle"] = selectedVehicle

	// get mileage logs for vehicle id
	logs, err := m.DB.GetMileageLogsByVehicleID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["mileage-logs"] = logs

	render.Template(w, r, "mileage-log-list.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

func (m *Repository) MileageLogCreate(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "edit-mileage-log.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) MileageLogCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse received form values into mileage log object
	v := models.MileageLog{}
	err := helpers.ParseFormToMileageLog(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["mileage-log"] = v

		render.Template(w, r, "edit-mileage-log.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	insertedID, err := m.DB.InsertMileageLog(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mileage-logs/list/%d", insertedID), http.StatusSeeOther)
}

func (m *Repository) MileageLogEdit(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) MileageLogEditPost(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) MileageLogDelete(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) TripsEdit(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) TripsEditPost(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

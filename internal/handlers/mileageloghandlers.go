package handlers

import (
	"fmt"
	"html/template"
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
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["vehicles"] = vehicles

	render.Template(w, r, "edit-mileage-log.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
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

		vehicles, err := m.DB.GetVehicleByActive(true)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data["vehicles"] = vehicles

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
	data := make(map[string]interface{})
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get vehicles for dropdown list
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["vehicles"] = vehicles

	// get mileage log from database
	v, err := m.DB.GetMileageLogByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["mileage-log"] = v

	render.Template(w, r, "edit-mileage-log.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) MileageLogEditPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get mileage log from database
	v, err := m.DB.GetMileageLogByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form into fetched mileage log
	err = helpers.ParseFormToMileageLog(r, &v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	// do form validation checks

	if !form.Valid() {
		data := make(map[string]interface{})
		data["mileage-log"] = v

		vehicles, err := m.DB.GetVehicleByActive(true)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		data["vehicles"] = vehicles

		render.Template(w, r, "edit-mileage-log.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.UpdateMileageLog(v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "flash", "Updated mileage log successfully")
	http.Redirect(w, r, fmt.Sprintf("/mileage-logs/%d", id), http.StatusSeeOther)
}

func (m *Repository) MileageLogDelete(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func createRiderOptionsTomSelect(members []models.Member) template.JS {
	optString := ``

	for _, member := range members {
		aliasNames := []string{}
		for _, a := range member.Aliases {
			aliasNames = append(aliasNames, a.Name)
		}
		aliasString := strings.Join(aliasNames, ",")

		memberString := fmt.Sprintf("{id: %d, name: '%s', aliases: '%s'},",
			member.ID, member.Name, aliasString)
		optString = optString + memberString
	}
	return template.JS(optString)
}

func (m *Repository) TripsEdit(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get mileage log from database
	data := make(map[string]interface{})
	v, err := m.DB.GetMileageLogByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["mileage-log"] = v

	// get Members from database & send as TomSelect compatible for rider selection
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["member-options"] = createRiderOptionsTomSelect(members)

	// calculate last odometer value from trips & mileage log start odometer
	intmap := make(map[string]int)
	intmap["last-odometer-value"] = calcLastOdometerValue(v)

	render.Template(w, r, "edit-mileage-log-trips.page.tmpl", &models.TemplateData{
		Form:   forms.New(nil),
		Data:   data,
		IntMap: intmap,
	})

}

func (m *Repository) TripsEditPost(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get mileage log from database
	v, err := m.DB.GetMileageLogByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	t := models.Trip{}
	// parse form into trip
	err = helpers.ParseFormToTrip(r, &t, v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	_, err = m.DB.InsertTrip(t)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/mileage-logs/%d/edit-trips", id), http.StatusSeeOther)
}

// calcLastOdometerValue takes a mileage log and calculates
// the last odometer value from trips entered
func calcLastOdometerValue(log models.MileageLog) int {
	if len(log.Trips) == 0 {
		// if there are no trips yet, return StartOdometer from the mileage log
		return log.StartOdometer
	}
	return log.Trips[len(log.Trips)-1].EndMileage
}

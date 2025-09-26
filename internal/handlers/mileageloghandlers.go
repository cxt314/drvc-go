package handlers

import (
	"bytes"
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

func (m *Repository) getTripEditTemplateData(mileageLogId int) (*models.TemplateData, error) {
	td := models.TemplateData{}

	// get mileage log from database
	data := make(map[string]interface{})
	v, err := m.DB.GetMileageLogByID(mileageLogId)
	if err != nil {
		return &td, err
	}

	data["mileage-log"] = v

	// get Members from database & send as TomSelect compatible for rider selection
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		return &td, err
	}

	data["member-options"] = createRiderOptionsTomSelect(members)

	td.Data = data

	// calculate last odometer value from trips & mileage log start odometer
	intmap := make(map[string]int)
	intmap["last-odometer-value"] = calcLastOdometerValue(v)

	td.IntMap = intmap
	return &td, nil
}

func (m *Repository) TripsEdit(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td, err := m.getTripEditTemplateData(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td.Form = forms.New(nil)

	render.Template(w, r, "edit-mileage-log-trips.page.tmpl", td)
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

// AddTripPost is the HTMX route that inserts a new trip
// On successful insert, returns a new trip form & the table row of the new trip
// If unsuccessful, returns the trip form with data & errors
func (m *Repository) AddTripPost(w http.ResponseWriter, r *http.Request) {
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

	buf := new(bytes.Buffer)
	form := forms.New(r.PostForm)
	// do form validation checks
	form.Required("trip-day", "start-mileage", "end-mileage", "end-mileage-input", "riders")

	// if there were errors, only generate the partial form w/ errors
	if !form.Valid() {
		td, err := m.getTripEditTemplateData(id)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		td.Form = form
		render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripForm", td)
		buf.WriteTo(w)
		return
	}

	w.Header().Set("Content-Type", "text/html")

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

	td, err := m.getTripEditTemplateData(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td.Form = forms.New(nil)
	// create HTMX response
	render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripForm", td)
	render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripTableSwap", td)

	//fmt.Println(buf.String())

	buf.WriteTo(w)
}

// calcLastOdometerValue takes a mileage log and calculates
// the last odometer value from trips entered
func calcLastOdometerValue(log models.MileageLog) int {
	if len(log.Trips) == 0 {
		// if there are no trips yet, return StartOdometer from the mileage log
		return log.StartOdometer
	}
	return log.Trips[0].EndMileage
}

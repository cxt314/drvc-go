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

	data["billing-rates"] = models.BillingRates

	data["ld-days"] = models.LongDistanceDays

	td.Data = data

	// calculate last odometer value from trips & mileage log start odometer
	intmap := make(map[string]int)
	intmap["last-odometer-value"] = calcLastOdometerValue(v)

	td.IntMap = intmap
	return &td, nil
}

func (m *Repository) getBillingTemplateData(mileageLogId int) (*models.TemplateData, error) {
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

	data["members"] = members

	td.Data = data

	// total cost of all trips
	data["total-trip-cost"] = m.calcTotalTripCost(v)

	// breakdown of cost by member
	data["member-billings"] = m.calcPerMemberBillings(v, members)

	return &td, nil
}

func (m *Repository) calcTotalTripCost(log models.MileageLog) models.USD {
	totalCost := models.ToUSD(0.0)

	for _, v := range log.Trips {
		totalCost = totalCost.AddUSD(v.Cost())
		//fmt.Printf("Trip Cost: %s Total Cost: %s", v.Cost(), totalCost)
	}

	return totalCost
}

func (m *Repository) calcPerMemberBillings(log models.MileageLog, members []models.Member) []models.MemberMileageLogBilling {
	var memberBillings []models.MemberMileageLogBilling

	tripMap := make(map[int]models.USD)
	ldMap := make(map[int]models.USD)

	for _, v := range log.Trips {
		numRiders := len(v.Riders)
		tripShare := v.Cost().Divide(float64(numRiders))

		for _, r := range v.Riders {
			if v.LongDistanceDays > 0 {
				// long distance trip, add to ldMap
				ldMap[r.ID] = ldMap[r.ID].AddUSD(tripShare)

			} else {
				// regular trip, add to tripMap
				tripMap[r.ID] = tripMap[r.ID].AddUSD(tripShare)
			}
		}
	}

	for _, v := range members {
		memberBilling := models.MemberMileageLogBilling{
			Member: v,
			RegularTripsCost: tripMap[v.ID],
			LongDistanceTripsCost: ldMap[v.ID],
		}

		memberBillings = append(memberBillings, memberBilling)
		
	}

	//fmt.Println(tripMap)
	//fmt.Println(ldMap)
	//fmt.Println(memberBillings)

	return memberBillings
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
	//fmt.Println(r.PostForm)

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

// EditTrip is the HTMX route that returns an edit trip form for a trip id
func (m *Repository) EditTrip(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get trip by id
	t, err := m.DB.GetTripByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td, err := m.getTripEditTemplateData(t.MileageLog.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td.Data["trip"] = t
	td.Form = forms.New(nil)

	buf := new(bytes.Buffer)

	// create HTMX response
	render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripEditForm", td)

	//fmt.Println(buf.String())

	buf.WriteTo(w)

}

// EditTripPost is the HTMX route that updates a trip
func (m *Repository) EditTripPost(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	// get trip by id
	t, err := m.DB.GetTripByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get later trips in case we need to update Start and End Mileages
	originalEndMileage := t.EndMileage
	laterTrips, err := m.DB.GetLaterTrips(t)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	form := forms.New(r.PostForm)

	// do form validation checks
	form.Required("trip-day", "start-mileage", "end-mileage", "riders")
	// check for valid end mileage if there are any later trips
	if len(laterTrips) > 0 {
		form.IsValidEndMileage("end-mileage", t.StartMileage, originalEndMileage, laterTrips[0].EndMileage)
	}

	// if there were errors, re-generate the partial form w/ errors
	if !form.Valid() {
		td, err := m.getTripEditTemplateData(t.MileageLog.ID)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		td.Data["trip"] = t
		td.Form = form
		render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripEditTableSwapError", td)
		buf.WriteTo(w)
		return
	}

	// get mileage log
	v, err := m.DB.GetMileageLogByID(t.MileageLog.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// parse form into trip
	err = helpers.ParseFormToTrip(r, &t, v)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// update later trips if end mileage changed
	if len(laterTrips) > 0 {
		err = m.updateFutureTripMileages(t, laterTrips, originalEndMileage)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}

	// update trip
	err = m.DB.UpdateTripByID(t)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td, err := m.getTripEditTemplateData(t.MileageLog.ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td.Form = forms.New(nil)

	// create HTMX response
	render.PartialHTMX(buf, r, "edit-mileage-log-trips.page.tmpl", "tripEditTableSwap", td)

	buf.WriteTo(w)

}

// updateFutureTripMileages updates subsequent start & end mileages when a trip is updated
// TODO: what if we need to reduce end mileage?
func (m *Repository) updateFutureTripMileages(updated models.Trip, laterTrips []models.Trip, originalEndMileage int) error {
	//fmt.Println("In updateFutureTripMileages")
	diff := updated.EndMileage - originalEndMileage
	//fmt.Printf("Mileage diff: %d\n", diff)

	if diff == 1000 {
		//fmt.Println("1000 mile roll-over")
		// 1000 roll-over: increase start and end mileages of all future trips by 1000
		for _, t := range laterTrips {
			t.StartMileage += 1000
			t.EndMileage += 1000

			err := m.DB.UpdateTripByID(t)
			//fmt.Printf("Trip ID: %d Mileage Log ID: %d StartMileage: %d EndMileage %d\n", t.ID, t.MileageLog.ID, t.StartMileage, t.EndMileage)
			if err != nil {
				return err
			}
		}

	} else if diff != 0 && diff < 1000 && diff > -1000 {
		//fmt.Println("Mileage changed by less than 1000")
		// only update the next trip's starting mileage. next trip's END mileage should be greater than updated.EndMileage
		nextTrip := laterTrips[0]
		nextTrip.StartMileage = updated.EndMileage

		if nextTrip.StartMileage > nextTrip.EndMileage {
			return fmt.Errorf("start mileage (%d) cannot be greater than end mileage (%d)", nextTrip.StartMileage, nextTrip.EndMileage)
		}

		//fmt.Printf("Trip ID: %d Mileage Log ID: %d StartMileage: %d EndMileage %d\n", nextTrip.ID, nextTrip.MileageLog.ID, nextTrip.StartMileage, nextTrip.EndMileage)
		err := m.DB.UpdateTripByID(nextTrip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Repository) MileageLogBilling(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td, err := m.getBillingTemplateData(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	render.Template(w, r, "mileage-log-billing.page.tmpl", td)
}

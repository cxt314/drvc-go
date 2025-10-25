package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

// // MemberList displays a list of all members
// func (m *Repository) MemberList(w http.ResponseWriter, r *http.Request) {
// 	members, err := m.DB.GetMemberByActive(true)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}
// 	data := make(map[string]interface{})
// 	data["members"] = members

// 	render.Template(w, r, "member-list.page.tmpl", &models.TemplateData{
// 		Data: data,
// 	})
// }

// // MemberCreate displays the page to create a new member
// func (m *Repository) MemberCreate(w http.ResponseWriter, r *http.Request) {
// 	render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
// 		Form: forms.New(nil),
// 	})
// }

// // MemberCreatePost processes the POST request for creating a new member
// func (m *Repository) MemberCreatePost(w http.ResponseWriter, r *http.Request) {
// 	// parse received form values into member object
// 	v := models.Member{}
// 	err := helpers.ParseFormToMember(r, &v)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	form := forms.New(r.PostForm)
// 	// do form validation checks

// 	if !form.Valid() {
// 		data := make(map[string]interface{})
// 		data["member"] = v

// 		render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
// 			Form: form,
// 			Data: data,
// 		})
// 		return
// 	}

// 	err = m.DB.InsertMember(v)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	http.Redirect(w, r, "/members", http.StatusSeeOther)
// }

// // MemberEdit shows the edit form for a member by id
// func (m *Repository) MemberEdit(w http.ResponseWriter, r *http.Request) {
// 	exploded := strings.Split(r.RequestURI, "/")
// 	id, err := strconv.Atoi(exploded[2])
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	// get member from database
// 	v, err := m.DB.GetMemberByID(id)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	data := make(map[string]interface{})
// 	data["member"] = v

// 	render.Template(w, r, "edit-member.page.tmpl", &models.TemplateData{
// 		Form: forms.New(nil),
// 		Data: data,
// 	})
// }

// BillingSummaryYearMonth gives a summary billing for each member and each vehicle by year and month
func (m *Repository) BillingSummaryYearMonth(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	year, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	month, err := strconv.Atoi(exploded[3])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	td, err := m.getSummaryBillingTemplateData(year, month)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	render.Template(w, r, "billing-summary.page.tmpl", td)

}

func (m *Repository) getMileageLogBilling(log models.MileageLog, members []models.Member) (models.MileageLogBilling, error) {
	var logBilling models.MileageLogBilling

	logBilling.Log = log

	// total cost of all trips
	logBilling.TotalTripCost = m.calcTotalTripCost(log)

	// breakdown of cost by member
	memberBillings := m.calcPerMemberBillings(log, members)
	logBilling.MemberBills = memberBillings

	// total cost of member billings (for checksum)
	logBilling.TotalMemberBillings = m.calcTotalMemberBillingsCost(memberBillings)

	return logBilling, nil
}

func (m *Repository) getSummaryBillingDisplay(vehicleBills map[string]models.MileageLogBilling, members []models.Member) [][]string {
	var displayArray [][]string


	// create & append header row
	var header []string
	header = append(header, "Member")

	displayArray = append(displayArray, header)

	return displayArray
}

func (m *Repository) getSummaryBillingTemplateData(year int, month int) (*models.TemplateData, error) {
	td := models.TemplateData{}

	// get mileage logs for given year and month from database
	logs, err := m.DB.GetMileageLogsByYearMonth(year, month)
	if err != nil {
		return &td, err
	}

	// get Members from database
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		return &td, err
	}

	mileageLogBills := make(map[string]models.MileageLogBilling)

	for _, v := range logs {
		mileageLogBill, err := m.getMileageLogBilling(v, members)
		if err != nil {
			return &td, err
		}

		mileageLogBills[v.Vehicle.Name] = mileageLogBill
	}

	// get active vehicles
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		return &td, err
	}

	data := make(map[string]interface{})
	data["vehicles"] = vehicles
	data["mileage-log-bills"] = mileageLogBills

	intmap := make(map[string]int)
	intmap["year"] = year
	intmap["month"] = month

	td.IntMap = intmap
	td.Data = data

	return &td, nil
}

func (m *Repository) getMileageLogBillingTemplateData(mileageLogId int) (*models.TemplateData, error) {
	td := models.TemplateData{}

	// get mileage log from database
	data := make(map[string]interface{})
	v, err := m.DB.GetMileageLogByID(mileageLogId)
	if err != nil {
		return &td, err
	}

	data["mileage-log"] = v

	// get Members from database
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		return &td, err
	}

	data["members"] = members

	// total cost of all trips
	data["total-trip-cost"] = m.calcTotalTripCost(v)

	// breakdown of cost by member
	memberBillings := m.calcPerMemberBillings(v, members)
	data["member-billings"] = memberBillings

	// total cost of member billings (for checksum)
	data["total-member-billings"] = m.calcTotalMemberBillingsCost(memberBillings)

	td.Data = data

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

func (m *Repository) calcPerMemberBillings(log models.MileageLog, members []models.Member) map[int]models.MemberMileageLogBilling {
	memberBillings := make(map[int]models.MemberMileageLogBilling)

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
			Member:                v,
			RegularTripsCost:      tripMap[v.ID],
			LongDistanceTripsCost: ldMap[v.ID],
		}

		memberBillings[v.ID] = memberBilling

	}

	return memberBillings
}

func (m *Repository) calcTotalMemberBillingsCost(billings map[int]models.MemberMileageLogBilling) models.USD {
	var totalUSD models.USD

	for _, v := range billings {
		totalUSD = totalUSD.AddUSD(v.RegularTripsCost)
		totalUSD = totalUSD.AddUSD(v.LongDistanceTripsCost)
	}

	return totalUSD
}

// // MemberDeactivate updates a member's active status by id
// func (m *Repository) MemberDeactivate(w http.ResponseWriter, r *http.Request) {
// 	exploded := strings.Split(r.RequestURI, "/")
// 	id, err := strconv.Atoi(exploded[2])
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	// update member set is_active to false
// 	err = m.DB.UpdateMemberActiveByID(id, false)
// 	if err != nil {
// 		helpers.ServerError(w, err)
// 		return
// 	}

// 	// redirect to members list after making member inactive
// 	http.Redirect(w, r, "/members", http.StatusSeeOther)
// }

// func (m *Repository) AddAlias(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	html := `
// 		<div class="col-3 mb-2 alias-group">
// 			<div class="input-group">
// 				<input class="form-control" type="text" name="aliases">
// 				<span class="input-group-text">
// 					<button type="button" class="btn-close" aria-label="Close"
// 						hx-get="/remove-item" hx-trigger="click" hx-target="closest .alias-group" hx-swap="outerHTML" hx-confirm="Delete alias?"></button>
// 				</span>
// 			</div>
// 		</div>
// 		`

// 	w.Write([]byte(html))
// }

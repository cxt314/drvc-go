package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cxt314/drvc-go/internal/helpers"
	"github.com/cxt314/drvc-go/internal/models"
	"github.com/cxt314/drvc-go/internal/render"
)

// BillingIndex redirects to Billing Summary for the previous year/month
func (m *Repository) BillingIndex(w http.ResponseWriter, r *http.Request) {
	year := time.Now().Year()
	month := time.Now().Month() - 1

	http.Redirect(w, r, fmt.Sprintf("/billings/%04d/%02d", year, month), http.StatusSeeOther)
}

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

// getSummaryBillingDisplay creates a slice of maps for a 2D display of all member billings per vehicle
func (m *Repository) getSummaryBillingDisplay(vehicleBills map[string]models.MileageLogBilling, members []models.Member, vehicles []models.Vehicle) ([]map[string]string, []string) {
	var displayArray []map[string]string
	var keyOrder []string

	//fmt.Println(vehicleBills)
	// create & append header row
	keyOrder = append(keyOrder, "Member")
	for _, v := range vehicles {
		keyOrder = append(keyOrder, v.Name)
		keyOrder = append(keyOrder, v.Name+" LD")
	}
	keyOrder = append(keyOrder, "Total")

	for _, i := range members {
		row := make(map[string]string)
		row["Member"] = i.Name

		memberTotal := models.ToUSD(0.0)

		for _, v := range vehicles {
			row[v.Name] = vehicleBills[v.Name].MemberBills[i.ID].RegularTripsCost.String()
			memberTotal = memberTotal.AddUSD(vehicleBills[v.Name].MemberBills[i.ID].RegularTripsCost)
			row[v.Name+" LD"] = vehicleBills[v.Name].MemberBills[i.ID].LongDistanceTripsCost.String()
			memberTotal = memberTotal.AddUSD(vehicleBills[v.Name].MemberBills[i.ID].LongDistanceTripsCost)
		}

		row["Total"] = memberTotal.String()
		displayArray = append(displayArray, row)
	}

	//fmt.Println(displayArray)
	return displayArray, keyOrder
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

	//fmt.Println(logs)
	for _, v := range logs {
		mileageLogBill, err := m.getMileageLogBilling(v, members)
		if err != nil {
			return &td, err
		}
		//fmt.Println(mileageLogBill.MemberBills)
		mileageLogBills[v.Vehicle.Name] = mileageLogBill
	}

	// get active vehicles
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		return &td, err
	}

	billDisplay, keyOrder := m.getSummaryBillingDisplay(mileageLogBills, members, vehicles)

	data := make(map[string]interface{})
	data["vehicles"] = vehicles
	data["mileage-log-bills"] = mileageLogBills
	data["bill-display"] = billDisplay
	data["key-order"] = keyOrder

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

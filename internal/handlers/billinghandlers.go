package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cxt314/drvc-go/internal/config"
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

func (m *Repository) BillingSummaryPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
	}

	// parse year & month str to int
	year, err := strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		helpers.ServerError(w, err)
	}

	// parse year str to int
	month, err := strconv.Atoi(r.Form.Get("month"))
	if err != nil {
		helpers.ServerError(w, err)
	}

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
	// create column headers
	keyOrder = append(keyOrder, "Member")
	for _, v := range vehicles {
		keyOrder = append(keyOrder, v.Name)
		keyOrder = append(keyOrder, v.Name+" LD")
	}
	keyOrder = append(keyOrder, "Total")

	// create & append row for each member
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

		// only append to display if member's bill is > 0
		if memberTotal != 0 {
			displayArray = append(displayArray, row)
		}
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

	// calculate member bills for each log
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

// calcPerMemberBillings takes a mileage log and a list of members
// and returns a map of member id to a models.MemberMileageLogBilling
// that splits LD trips from regular trips for that member
func (m *Repository) calcPerMemberBillings(log models.MileageLog, members []models.Member) map[int]models.MemberMileageLogBilling {
	memberBillings := make(map[int]models.MemberMileageLogBilling)

	tripMap := make(map[int]models.USD)
	ldMap := make(map[int]models.USD)
	tripMapFloat := make(map[int]float64)
	ldMapFloat := make(map[int]float64)

	for _, v := range log.Trips {
		numRiders := len(v.Riders)
		// when summing up trip costs, keep things as a float to avoid rounding issues
		tripShare := v.Cost().Float64() / float64(numRiders)

		for _, r := range v.Riders {
			if v.LongDistanceDays > 0 {
				// long distance trip, add to ldMap
				ldMapFloat[r.ID] = ldMapFloat[r.ID] + tripShare
			} else {
				// regular trip, add to tripMap
				tripMapFloat[r.ID] = tripMapFloat[r.ID] + tripShare
			}
		}
	}

	// convert trip cost floats to USD
	for k, v := range tripMapFloat {
		tripMap[k] = models.ToUSD(v)
	}
	for k, v := range ldMapFloat {
		ldMap[k] = models.ToUSD(v)
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

func (m *Repository) BillingCreateMileageLogs(w http.ResponseWriter, r *http.Request) {
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

	// get all active vehicles
	vehicles, err := m.DB.GetVehicleByActive(true)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// create new mileage log for each active vehicle for given year/month
	for _, v := range vehicles {
		err = m.createMileageLogStub(v, year, month)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/billings/%d/%d", year, month), http.StatusSeeOther)
}

// BillingCSV generates a csv download of mileage logs for all vehicles in a given billing
func (m *Repository) BillingCSV(w http.ResponseWriter, r *http.Request) {
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

	logs, err := m.DB.GetMileageLogsByYearMonth(year, month)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Set headers so browser will download the file
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=mileage-logs-%04d%02d.csv", year, month))
	w.Header().Set("Transfer-Encoding", "chunked")

	// Create a CSV writer using our HTTP response writer as our io.Writer
	wr := csv.NewWriter(w)

	for _, l := range logs {
		logCSV := convertMileageLogToCSVRaw(l)

		if err := wr.WriteAll(logCSV); err != nil {
			helpers.ServerError(w, err)
			return
		}
	}

	// Flush the writer and check for any errors
	wr.Flush()
	if err := wr.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
		helpers.ServerError(w, err)
		return
	}

}

// QBOBulkInvoicesCSV generates a csv for download that uses Quickbooks Online's bulk import
// invoice via csv tool to quickly transfer a monthly billing to Quickbooks invoices
// TODO: how do we add gas mileage that needs to be tacked on? this is happening manually atm
func (m *Repository) QBOBulkInvoicesCSV(w http.ResponseWriter, r *http.Request) {
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

	logs, err := m.DB.GetMileageLogsByYearMonth(year, month)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Set headers so browser will download the file
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=qbo-bulk-invoices-%04d%02d.csv", year, month))
	w.Header().Set("Transfer-Encoding", "chunked")

	// Create a CSV writer using our HTTP response writer as our io.Writer
	wr := csv.NewWriter(w)

	// add header row
	headerRow := m.getQBOInvoicesHeaderRow()
	wr.Write(headerRow)

	// get invoice lines for each mileage log
	for _, l := range logs {
		logCSV := m.convertMileageLogToQBOInvoiceLineRaw(l)

		if err := wr.WriteAll(logCSV); err != nil {
			helpers.ServerError(w, err)
			return
		}
	}

	// Flush the writer and check for any errors
	wr.Flush()
	if err := wr.Error(); err != nil {
		fmt.Println("Error flushing CSV writer:", err)
		helpers.ServerError(w, err)
		return
	}

}

// getQBOInvoicesHeaderRow returns a csvSlice of the headers for QBO's bulk invoices
func (m *Repository) getQBOInvoicesHeaderRow() []string {

	headerRow := []string{"*InvoiceNo", "*Customer", "*InvoiceDate", "*DueDate", "Terms", "Location",
		"Memo", "Item(Product/Service)", "ItemDescription", "ItemQuantity", "ItemRate",
		"*ItemAmount", "Class", "Shipping address", "Ship via", "Shipping date",
		"Tracking no.", "Shipping Charge", "Service Date"}

	return headerRow
}

// convertMileageLogToQBOInvoiceLineRaw convers a mileage log to a csv string
// containing QBO invoice line items for every member with a non-zero billing
// for that log
func (m *Repository) convertMileageLogToQBOInvoiceLineRaw(log models.MileageLog) [][]string {
	csvSlice := [][]string{{}}

	// get active members
	members, err := m.DB.GetMemberByActive(true)
	if err != nil {
		return csvSlice
	}

	// get per member billings for the log
	memberBillings := m.calcPerMemberBillings(log, members)

	// create new lines for each member in the per member billing where either RegularTripsCost or LongDistanceTripsCost are non-zero
	for k, v := range memberBillings {
		// use QBOName for customer name, unless QBOName is empty
		customer := v.Member.QBOName
		if customer == "" {
			customer = v.Member.Name
		}

		// skip the line if the customer is DRVC
		if customer == "DRVC" {
			continue
		}

		// calculate invoice date & due date
		nextMonthFirstDay := time.Date(log.Year, time.Month(log.Month+1), 1, 0, 0, 0, 0, time.UTC)
		invoiceDate := nextMonthFirstDay.AddDate(0, 0, -1)
		dueDate := invoiceDate.AddDate(0, 0, 15)

		// check if amount owed is > 0 for trips
		if v.RegularTripsCost > 0 {
			tripRow := []string{
				strconv.Itoa(k),                       // *InvoiceNo - use key (member id)
				customer,                              // *Customer - use QBOName if not empty, else member name
				invoiceDate.Format(config.DateLayout), // *InvoiceDate - last date of Mileage Log's month
				dueDate.Format(config.DateLayout),     // *DueDate - 15 days from InvoiceDate
				"Net 15",                              // Terms
				"",                                    // Location - blank
				"",                                    // Memo - blank
				"Mileage Fee",                         // Item(Product/Service)
				log.Vehicle.Name,                      // ItemDescription - name of vehicle
				"1",                                   // ItemQuantity
				v.RegularTripsCost.String(),           // ItemRate
				v.RegularTripsCost.String(),           // *ItemAmount
				log.Vehicle.QBOClass,                  // Class
				"",                                    // Shipping address
				"",                                    // Ship via - blank
				"",                                    // Shipping date - blank
				"",                                    // Tracking no. - blank
				"",                                    // Shipping Charge - blank
				"",                                    // Service Date - blank
			}

			csvSlice = append(csvSlice, tripRow)
		}

		// check if amount owed is > = for ld
		if v.LongDistanceTripsCost > 0 {
			ldRow := []string{
				strconv.Itoa(k),                       // *InvoiceNo - use key (member id)
				customer,                              // *Customer - use QBOName if not empty, else member name
				invoiceDate.Format(config.DateLayout), // *InvoiceDate - last date of Mileage Log's month
				dueDate.Format(config.DateLayout),     // *DueDate - 15 days from InvoiceDate
				"Net 15",                              // Terms
				"",                                    // Location - blank
				"",                                    // Memo - blank
				"Long Distance",                       // Item(Product/Service)
				log.Vehicle.Name,                      // ItemDescription - name of vehicle
				"1",                                   // ItemQuantity
				v.LongDistanceTripsCost.String(),      // ItemRate
				v.LongDistanceTripsCost.String(),      // *ItemAmount
				log.Vehicle.QBOClass,                  // Class
				"",                                    // Shipping address
				"",                                    // Ship via - blank
				"",                                    // Shipping date - blank
				"",                                    // Tracking no. - blank
				"",                                    // Shipping Charge - blank
				"",                                    // Service Date - blank
			}

			csvSlice = append(csvSlice, ldRow)
		}
	}

	return csvSlice
}

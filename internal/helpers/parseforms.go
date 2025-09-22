package helpers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// date_layout is the format we expect dates to be sent in as
const date_layout = "2006-01-02" // 01/02 03:04:05PM '06 -0700

func ParseFormToVehicle(r *http.Request, v *models.Vehicle) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// parse string fields
	v.Name = r.Form.Get("name")
	v.Make = r.Form.Get("make")
	v.Model = r.Form.Get("model")
	v.FuelType = r.Form.Get("fuel_type")
	v.Vin = r.Form.Get("vin")
	v.LicensePlate = r.Form.Get("license_plate")

	// parse year str to int
	v.Year, err = strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		return err
	}

	// parse purchase date string to *time.Time
	tempPD, err := time.Parse(date_layout, r.Form.Get("purchase_date"))
	if err != nil {
		return err
	}
	v.PurchaseDate = &tempPD
	v.PurchasePrice = models.StrToUSD(r.Form.Get("purchase_price"))

	// parse sale date string to *time.Time
	tempSD, err := time.Parse(date_layout, r.Form.Get("sale_date"))
	if err != nil {
		return err
	}
	v.SaleDate = &tempSD
	v.SalePrice = models.StrToUSD(r.Form.Get("sale_price"))

	return nil
}

func ParseFormToMember(r *http.Request, v *models.Member) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// parse string fields
	v.Name = r.Form.Get("name")
	v.Email = r.Form.Get("email")

	// parse member aliases
	// clear out old aliases
	v.Aliases = []models.MemberAlias{}

	formAliases := r.Form["aliases"]
	for _, a := range formAliases {
		if a != "" {
			v.Aliases = append(v.Aliases, models.MemberAlias{
				Name: a,
			})
		}
	}
	//log.Println(v)

	return nil
}

func ParseFormToMileageLog(r *http.Request, v *models.MileageLog) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// parse string fields
	v.Name = r.Form.Get("name")

	// parse number fields
	v.Year, err = strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		return err
	}

	v.Month, err = strconv.Atoi(r.Form.Get("month"))
	if err != nil {
		return err
	}

	v.Vehicle.ID, err = strconv.Atoi(r.Form.Get("vehicle"))
	if err != nil {
		return err
	}

	v.StartOdometer, err = strconv.Atoi(r.Form.Get("start_odometer"))
	if err != nil {
		return err
	}

	v.EndOdometer, err = strconv.Atoi(r.Form.Get("end_odometer"))
	if err != nil {
		return err
	}

	return nil
}

func ParseFormToTrip(r *http.Request, v *models.Trip, log models.MileageLog) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	fmt.Println(r.Form)

	v.MileageLog = log

	// parse number fields
	v.StartMileage, err = strconv.Atoi(r.Form.Get("start-mileage"))
	if err != nil {
		return err
	}

	// TODO: switch back to "end-mileage" once js is being called
	v.EndMileage, err = strconv.Atoi(r.Form.Get("end-mileage-input"))
	if err != nil {
		return err
	}

	// build TripDate from trip-day form input & mileage log year/month
	tripDay, err := strconv.Atoi(r.Form.Get("trip-day"))
	if err != nil {
		return err
	}

	v.TripDate, err = time.Parse(date_layout, fmt.Sprintf("%d-%02d-%02d", log.Year, log.Month, tripDay))
	if err != nil {
		return err
	}

	// process riders
	riders := r.Form["riders"]
	//fmt.Println(riders)
	for _, riderID := range riders {
		newRider := models.Member{}
		newRider.ID, err = strconv.Atoi(riderID)
		if err != nil {
			return err
		}

		v.Riders = append(v.Riders, newRider)
	}

	// parse string fields
	v.Destination = r.Form.Get("destination")
	v.Purpose = r.Form.Get("purpose")

	//fmt.Println(v)
	return nil
}

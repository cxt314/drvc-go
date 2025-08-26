package helpers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/cxt314/drvc-go/internal/models"
)

// date_layout is the format we expect dates to be sent in as
const date_layout = "2006-01-02" // 01/02 03:04:05PM '06 -0700

func ParseFormToVehicle(r *http.Request, v *models.Vehicle) error {
	//v := models.Vehicle{}
	err := r.ParseForm()
	if err != nil {
		return err
	}

	v.Name = r.Form.Get("name")
	v.Make = r.Form.Get("make")
	v.Model = r.Form.Get("model")
	v.FuelType = r.Form.Get("fuel_type")
	v.Vin = r.Form.Get("vin")
	v.LicensePlate = r.Form.Get("license_plate")

	v.Year, err = strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		return err
	}

	pd := r.Form.Get("purchase_date")
	v.PurchaseDate, err = time.Parse(date_layout, pd)
	if err != nil {
		return err
	}
	v.PurchasePrice = models.StrToUSD(r.Form.Get("purchase_price"))

	return nil
}

package models

import "time"

// FuelTypes contains the list of possible fuel types a vehicle can be
var FuelTypes = map[string]string{
	"HY": "Hybrid",
	"EC": "Electric",
	"GS": "Gasoline",
	"DI": "Diesel",
}

// BillingTypes contains the list of allowed billing methods
var BillingTypes = [...]string{"Basic", "Truck"}

// BillingRates contains the list of allowed billing rates
var BillingRates = [...]string{"Primary", "Secondary"}

// LongDistanceDays contains the list of allowed days for a long distance trip
var LongDistanceDays = [...]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

// Vehicles is the vehicle model
type Vehicle struct {
	ID               int
	Name             string
	Year             int
	Make             string
	Model            string
	FuelType         string
	PurchasePrice    USD
	PurchaseDate     *time.Time
	Vin              string
	LicensePlate     string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Active           bool
	SalePrice        USD
	SaleDate         *time.Time
	BillingType      string
	BasePerMile      USD
	SecondaryPerMile USD
	MinimumFee       USD
}

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

// Vehicles is the vehicle model
type Vehicle struct {
	ID            int
	Name          string
	Year          int
	Make          string
	Model         string
	FuelType      string
	PurchasePrice USD
	PurchaseDate  *time.Time
	Vin           string
	LicensePlate  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Active        bool
	SalePrice     USD
	SaleDate      *time.Time
	BillingType   string
}

type BillingMethod interface{
	Name() string
	
}
package models

import "time"

// FuelTypes contains the list of possible fuel types a vehicle can be
var FuelTypes = [...]string{"Hybrid", "Electric", "Gasoline", "Diesel"}

// Vehicles is the vehicle model
type Vehicles struct {
	ID            int
	Name          string
	Year          int
	Make          string
	Model         string
	FuelType      string
	PurchasePrice USD
	Vin           string
	LicensePlate  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

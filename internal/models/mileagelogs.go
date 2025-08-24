package models

import "time"

// MileageLogs describes the mileage log model
// Each MileageLog is made up of multiple Trips
type MileageLogs struct {
	ID            int
	Vehicle       Vehicles
	Name          string
	Year          int
	Month         int
	StartOdometer int
	EndOdometer   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Trips describes the trip model
// Each trip can have multiple DRVC Members as riders
type Trips struct {
	ID             int
	MileageLog     MileageLogs
	TripDate       time.Time
	StartMileage   int
	EndMileage     int
	IsLongDistance bool
	Destination    string
	Purpose        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Riders describes the riders model
// This model associates a DRVC member with a trip
type Riders struct {
	ID        int
	Trip      Trips
	Member    Members
	CreatedAt time.Time
	UpdatedAt time.Time
}

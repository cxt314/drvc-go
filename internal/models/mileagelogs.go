package models

import "time"

// MileageLog describes the mileage log model
// Each MileageLog is made up of multiple Trips
type MileageLog struct {
	ID            int
	Vehicle       Vehicle
	Name          string
	Year          int
	Month         int
	StartOdometer int
	EndOdometer   int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Trip describes the trip model
// Each trip can have multiple DRVC Members as riders
type Trip struct {
	ID             int
	MileageLog     MileageLog
	TripDate       time.Time
	StartMileage   int
	EndMileage     int
	IsLongDistance bool
	Destination    string
	Purpose        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Rider describes the rider model
// This model associates a DRVC member with a trip
type Rider struct {
	ID        int
	Trip      Trip
	Member    Member
	CreatedAt time.Time
	UpdatedAt time.Time
}

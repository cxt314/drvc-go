package models

import (
	"time"
)

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
	Trips         []Trip
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
	Riders         []Member
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Distance returns the calcaulted value of the distance of the trip
func (t Trip) Distance() float64 {
	return float64(t.EndMileage - t.StartMileage)
}

// BillingMethod returns the billing method for a trip based on the vehicle. If the trip was Long Distance then long distance billing is used
func (t Trip) BillingMethod() BillingMethod {
	// if this is a long distance trip, we use long distance billing regardless of vehicle
	if t.IsLongDistance {
		return &LongDistanceBilling{
			BillingName:   "Long Distance",
			SingleDayRate: ToUSD(85.0),
			MultiDayRate:  ToUSD(50.0),
		}
	}

	if t.MileageLog.Vehicle.BillingType == "Basic" {
		return &SimplePerMileBilling{
			BillingName: "Simple Per Mile",
			BasePerMile: t.MileageLog.Vehicle.BasePerMile,
		}
	}

	if t.MileageLog.Vehicle.BillingType == "Truck" {
		return &TruckBilling{
			BillingName:      "Truck",
			BasePerMile:      t.MileageLog.Vehicle.BasePerMile,
			SecondaryPerMile: t.MileageLog.Vehicle.SecondaryPerMile,
			MinimumFee:       t.MileageLog.Vehicle.MinimumFee,
		}
	}

	return nil
}

func (t Trip) Cost() USD {
	if t.IsLongDistance {
		// TODO: implement long distance trip days tracking
		return t.BillingMethod().TripCost(0.0, false)
	}

	if t.BillingMethod() != nil {
		return t.BillingMethod().TripCost(t.Distance(), false)
	}

	return ToUSD(0.0)
}

// Rider describes the rider model
// This model describes the riders table which represents
// the M2M relationship between a DRVC member and a trip
type Rider struct {
	ID        int
	Trip      Trip
	Member    Member
	CreatedAt time.Time
	UpdatedAt time.Time
}

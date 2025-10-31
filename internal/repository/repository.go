package repository

import "github.com/cxt314/drvc-go/internal/models"

type DatabaseRepo interface {
	AllUsers() ([]models.User, error)
	UpdateUser(v models.User) error
	GetUserByID(id int) (models.User, error) 
	Authenticate(email, testPassword string) (int, string, error)

	InsertVehicle(v models.Vehicle) error
	AllVehicles() ([]models.Vehicle, error)
	GetVehicleByActive(active bool) ([]models.Vehicle, error)
	GetVehicleByID(id int) (models.Vehicle, error)
	UpdateVehicle(v models.Vehicle) error
	UpdateVehicleActiveByID(id int, active bool) error
	DeleteVehicle(id int) error

	InsertMember(v models.Member) error
	AllMembers() ([]models.Member, error)
	GetMemberByActive(active bool) ([]models.Member, error)
	GetMemberByID(id int) (models.Member, error)
	UpdateMember(v models.Member) error
	UpdateMemberActiveByID(id int, active bool) error
	DeleteMember(id int) error

	InsertMileageLog(v models.MileageLog) (int, error)
	AllMileageLogs() ([]models.MileageLog, error)
	GetMileageLogsByVehicleID(vehicle_id int) ([]models.MileageLog, error)
	GetMileageLogByID(id int) (models.MileageLog, error)
	UpdateMileageLog(v models.MileageLog) error
	DeleteMileageLog(id int) error
	InsertTrip(v models.Trip) (int, error)
	GetTripByID(id int) (models.Trip, error)
	UpdateTripByID(v models.Trip) error 
	GetLaterTrips(v models.Trip) ([]models.Trip, error)
	GetMileageLogsByYearMonth(year int, month int) ([]models.MileageLog, error)
}

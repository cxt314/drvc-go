package repository

import "github.com/cxt314/drvc-go/internal/models"

type DatabaseRepo interface {
	AllUsers() ([]models.User, error)
	UpdateUser(v models.User) error

	InsertVehicle(v models.Vehicle) error
	AllVehicles() ([]models.Vehicle, error)
	GetVehicleByActive(active bool) ([]models.Vehicle, error)
	GetVehicleByID(id int) (models.Vehicle, error)
	UpdateVehicle(v models.Vehicle) error
	DeleteVehicle(id int) error
}

package repository

import "github.com/cxt314/drvc-go/internal/models"

type DatabaseRepo interface {
	AllUsers() bool

	InsertVehicle(v models.Vehicle) error
	AllVehicles() ([]models.Vehicle, error)
	GetVehicleByID(id int) (models.Vehicle, error)
}

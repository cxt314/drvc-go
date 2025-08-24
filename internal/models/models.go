package models

import (
	"time"
)

// Reservation is a sample model from the learning course
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Users is the user model
type Users struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

package models

import "time"

// Member is the DRVC Member model.
// Email is not required to be unique for members
type Member struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemberAlias is the member alias model.
// Each Member can have multiple aliases by which they are referred
type MemberAlias struct {
	ID        int
	Member    Member
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

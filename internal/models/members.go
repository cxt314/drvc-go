package models

import "time"

// Members is the DRVC Members model.
// Email is not required to be unique for members
type Members struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MemberAliases is the member alias model.
// Each Member can have multiple aliases by which they are referred
type MemberAliases struct {
	ID        int
	Member    Members
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

package models

import "time"

// SurrogateModel is a base model which defines surrogate key field `ID`.
type SurrogateModel struct {
	ID uint `gorm:"primary_key"`
}

// AuditModel is a base model which defines `CreatedAt`, `Creator`, `UpdatedAt` and `Updator` fields
type AuditModel struct {
	CreatedAt time.Time
	Creator   string
	UpdatedAt time.Time
	Updater   string
}

package entity

import "github.com/google/uuid"

type Tag struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Location    string    `gorm:"type:varchar(100);not null" json:"location"`
	Class       string    `gorm:"type:varchar(20);not null" json:"class"`
	ReportCount int       `gorm:"not null" json:"report_count"`

	Reports []Report `gorm:"foreignKey:TagID" json:"reports,omitempty"`
	Timestamp
}

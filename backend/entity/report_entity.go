package entity

import "github.com/google/uuid"

type ReportStatus string

const (
	StatusUnverified ReportStatus = "unverified"
	StatusVerified   ReportStatus = "verified"
	StatusRejected   ReportStatus = "rejected"
	StatusHandled    ReportStatus = "handled"
	StatusCompleted  ReportStatus = "completed"
)

type Report struct {
	ID             string       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Text           string       `gorm:"type:text" json:"text"`
	Image          string       `gorm:"type:varchar(500)" json:"image"`
	Status         ReportStatus `gorm:"type:varchar(50);not null;default:'unverified'" json:"status"`
	PredConfidence *int         `gorm:"" json:"pred_confidence"`
	Upvotes        int          `gorm:"default:0" json:"upvotes"`
	ShareCount     int          `gorm:"default:0" json:"share_count"`
	Location       string       `gorm:"type:varchar(255)" json:"location"`

	UserID string `gorm:"type:char(32);not null"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	TagID uuid.UUID `gorm:"type:uuid"`
	Tag   Tag       `gorm:"foreignKey:TagID"`

	Timestamp
}

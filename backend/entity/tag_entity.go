package entity

type Tag struct {
	ID          string `gorm:"type:char(32);primaryKey"`
	Location    string `gorm:"type:varchar(100);not null"`
	Class       string `gorm:"type:varchar(20);not null"`
	ReportCount int    `gorm:"not null"`

	Reports []Report `gorm:"foreignKey:TagID"`
	Timestamp
}

package models

type Invoice struct {
	ID        uint    `gorm:"primaryKey"`
	Hours     float64 `gorm:"not null"`
	Total     float64 `gorm:"not null"`
	StudentID uint    `gorm:"not null"` // Foreign key
	Date      string  `gorm:"size:10;not null"`
}

package models

import "fmt"

type Invoice struct {
	ID        uint    `gorm:"primaryKey" form:"invoice_id"`
	Hours     float64 `gorm:"not null" form:"hours" binding:"required"`
	Total     float64 `gorm:"not null" form:"total" binding:"required"`
	StudentID uint    `gorm:"not null" form:"student_id" binding:"required"` // Foreign key
	Date      string  `gorm:"size:10;not null" form:"date"`
}

func (i Invoice) String() string {
	return fmt.Sprintf("<Invoice %d>", i.ID)
}

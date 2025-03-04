package models

type Student struct {
	ID           uint      `gorm:"primaryKey" form:"student_id"`
	Name         string    `gorm:"size:100;not null" form:"name" binding:"required"`
	Surname      string    `gorm:"size:100;not null" form:"surname" binding:"required"`
	Parent       string    `gorm:"size:100;not null" form:"parent" binding:"required"`
	Email        string    `gorm:"size:120;not null" form:"email" binding:"required,email"`
	Address      string    `gorm:"size:200;not null" form:"address" binding:"required"`
	PhoneNumber  string    `gorm:"size:30;not null" form:"phone_number" binding:"required" json:"phone_number"`
	PricePerHour int       `gorm:"not null" form:"price_per_hour" binding:"required" json:"price_per_hour"`
	InvoiceCount int       `gorm:"default:1;not null"`
	Invoices     []Invoice `gorm:"foreignKey:StudentID"`
}

func (s Student) String() string {
	return "<Student " + s.Name + " " + s.Surname + ">"
}

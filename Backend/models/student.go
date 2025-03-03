package models

type Student struct {
	ID           uint      `gorm:"primaryKey" json:"student_id"`
	Name         string    `gorm:"size:100;not null"`
	Surname      string    `gorm:"size:100;not null"`
	Parent       string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:120;not null"`
	Address      string    `gorm:"size:200;not null"`
	PhoneNumber  string    `gorm:"size:30;not null" json:"phone_number"`
	PricePerHour int       `gorm:"not null" json:"price_per_hour"`
	InvoiceCount int       `gorm:"default:1;not null"`
	Invoices     []Invoice `gorm:"foreignKey:StudentID"`
}

func (s Student) String() string {
	return "<Student " + s.Name + " " + s.Surname + ">"
}

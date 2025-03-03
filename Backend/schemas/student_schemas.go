package schemas

type StudentCreate struct {
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Parent       string `json:"parent"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	PricePerHour int    `json:"price_per_hour"`
}

type StudentDelete struct {
	ID int `json:"student_id" binding:"required"`
}

type StudentUpdate struct {
	ID           uint   `json:"student_id" binding:"required"`
	Name         string `json:"name,omitempty"`
	Surname      string `json:"surname,omitempty"`
	Parent       string `json:"parent,omitempty"`
	Email        string `json:"email,omitempty"`
	Address      string `json:"address,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	PricePerHour int    `json:"price_per_hour,omitempty"`
}

type StudentGet struct {
	ID           int    `json:"student_id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Parent       string `json:"parent"`
	Email        string `json:"email"`
	Address      string `json:"address"`
	PhoneNumber  string `json:"phone_number"`
	PricePerHour int    `json:"price_per_hour"`
}

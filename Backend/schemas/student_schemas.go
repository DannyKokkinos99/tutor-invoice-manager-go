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
	ID int `json:"student_id" binding:"required" form:"student_id"`
}

type StudentUpdate struct {
	ID           int    `form:"student_id"`
	Name         string `form:"name,omitempty"`
	Surname      string `form:"surname,omitempty"`
	Parent       string `form:"parent,omitempty"`
	Email        string `form:"email,omitempty"`
	Address      string `form:"address,omitempty"`
	PhoneNumber  string `form:"phone_number,omitempty"`
	PricePerHour int    `form:"price_per_hour,omitempty"`
}

type StudentGet struct {
	ID           int    `json:"student_id" form:"student_id"`
	Name         string `json:"name" form:"name"`
	Surname      string `json:"surname" form:"surname"`
	Parent       string `json:"parent" form:"parent"`
	Email        string `json:"email" form:"email"`
	Address      string `json:"address" form:"address"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
	PricePerHour int    `json:"price_per_hour" form:"price_per_hour"`
}

package student

import (
	"gopkg.in/guregu/null.v3"
)

//Student struct
type Student struct {
	ID         int         `json:"id" gorm:"primary key"`
	SeatNumber null.String `json:"seat_number,omitempty"`
	Gpa        float32     `json:"gpa"`
	Hours      int         `json:"hours"`
	UserID     int         `gorm:"unique_index:idx_student_user"`
}

//StudentRequest struct
type StudentRequest struct {
	ID        int     `json:"id"`
	Gpa       float32 `json:"gpa" validate:"omitempty"`
	Hours     int     `json:"hours" validate:"omitempty"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	UserName  string  `json:"user_name" validate:"required,email"`
	Password  string  `json:"password" validate:"required,gte=6,lte=13" gorm:"type:text"`
	Phone     string  `json:"phone"`
}

//StudentUpdate struct API
type StudentUpdate struct {
	Gpa       float32 `json:"gpa" validate:"required"`
	Hours     int     `json:"hours" validate:"required"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	UserName  string  `json:"user_name" validate:"required"`
	Telephone string  `json:"phone"`
}

// StudentResponse for student api
type StudentResponse struct {
	StudentID    int     `json:"id"`
	StudentGPA   float32 `json:"gpa"`
	StudentHours int     `json:"hours"`
	SeatNumber   string  `json:"seat_number"`
	User         struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"user_name"`
		Telephone string `json:"phone"`
	}
}

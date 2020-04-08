package enroll

// Enrollment struct
type Enrollment struct {
	ID        int `json:"id" gorm:"primary key"`
	StudentID int `json:"student_id" validate:"required"`
	CourseID  int `json:"course_id" validate:"required"`
}

package course

import (
	"gopkg.in/guregu/null.v3"
)

// Course struct table
type Course struct {
	ID          int      `json:"id" gorm:"primary key"`
	Name        string   `json:"name" validate:"required" gorm:"unique_index:idx_course_name_sid"`
	ProfessorID null.Int `json:"professor_id,omitempty"`
	SemesterID  null.Int `json:"semester_id,omitempty" gorm:"unique_index:idx_course_name_sid"`
}

// CourseSerializer struct
type CourseSerializer struct {
	ID            int    `json:"course_id"`
	CourseName    string `json:"course_name"`
	ProfessorName string `json:"professor_name"`
}

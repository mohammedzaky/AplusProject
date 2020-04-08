package studentExam

import (
	"gopkg.in/guregu/null.v3"
)

// StudentExam struct
type StudentExam struct {
	ID            int        `json:"id" gorm:"primary key"`
	StudentID     int        `json:"student_id" validate:"required"`
	ExamID        int        `json:"exam_id" validate:"required"`
	StudentDegree null.Float `json:"student_degree" validate:"required"`
}

// ExamResult struct
type ExamResult []struct {
	StudentSeatNumber string    `json:"seat_number"`
	StudentName       string    `json:"student_name"`
	ExamDegrees       []Degrees `json:"exams"`
}

//Degrees struct for each exam show student degree
type Degrees struct {
	ExamName      string     `json:"exam_name"`
	StudentDegree null.Float `json:"student_degree"`
}

//ExamResults show all students degrees for all exams in specefic course
type ExamResults []struct {
	ExamName       string           `json:"exam_name"`
	StudentDegrees []StudentDegrees `json:"students_Degrees"`
}

// StudentDegrees for more that one exam
type StudentDegrees struct {
	StudentSeatNumber string     `json:"seat_number"`
	StudentName       string     `json:"student_name"`
	StudentDegree     null.Float `json:"student_degree"`
}

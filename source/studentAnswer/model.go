package studentAnswer

import (
	"gopkg.in/guregu/null.v3"
)

//StudentAnswer struct
type StudentAnswer struct {
	ID            int      `json:"id" gorm:"primary key"`
	QuestionID    int      `json:"question_id" validate:"required"`
	StudentID     int      `json:"student_id" validate:"required"`
	StudentChoice null.Int `json:"student_choice"`
}

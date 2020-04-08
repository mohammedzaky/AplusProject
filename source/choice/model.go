package choice

//Choice struct
type Choice struct {
	ID         int    `json:"id" gorm:"primary key"`
	Name       string `json:"name" validate:"required"`
	QuestionID int    `json:"question_id" validate:"required"`
	IsCorrect  bool   `json:"is_correct" validate:"omitempty"`
}

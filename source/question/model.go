package question

//Question struct
type Question struct {
	ID         int     `json:"id" gorm:"primary key"`
	Name       string  `json:"name" validate:"required" gorm:"type:text;unique"`
	Degree     float32 `json:"degree" validate:"required"`
	ChoiceType string  `json:"choice_type" validate:"omitempty"` //if this field is has two options 1- multi choices 2- True/False
	ExamID     int     `json:"exam_id" validate:"required"`      //FK
}

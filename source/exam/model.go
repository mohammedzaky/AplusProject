package exam

import (
	choiceObject "gitlab.com/mohamedzaky/aplusProject/source/choice"
)

// Exam struct
type Exam struct {
	ID       int    `json:"id" gorm:"primary key"`
	Name     string `json:"name" validate:"required" gorm:"type:text"`
	Degree   int    `json:"degree" validate:"required"`
	Duration int    `json:"duration" validate:"required"`
	IsEnable bool   `json:"is_enable" validate:"omitempty"` //if this field is true then the exam will be shown else will be hidden
	CourseID int    `json:"course_id" validate:"required"`
}

//QuestionResponse struct for professor
type QuestionResponse []struct {
	ID         int                   `json:"id"`
	Name       string                `json:"name"`
	Degree     float32               `json:"degree"`
	ChoiceType string                `json:"choice_type"`
	ExamID     int                   `json:"exam_id"`
	Choices    []choiceObject.Choice `json:"choices"`
}

//QuestionRequest struct for professor
type QuestionRequest []struct {
	Name       string                `json:"name"`
	Degree     float32               `json:"degree"`
	ChoiceType string                `json:"choice_type"`
	ExamID     int                   `json:"exam_id"`
	Choices    []choiceObject.Choice `json:"choices"`
}

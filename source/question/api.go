package question

import (
	"encoding/json"
	"net/http"

	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewQuestions insert new question in db
func NewQuestions(c echo.Context) error {
	db := connectHandler.ConnectDB()
	question := new(Question)

	c.Bind(question)

	vaildate = validator.New()

	err := vaildate.Struct(question)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errCreate := db.Create(&question)

	if errCreate.RowsAffected == 0 {
		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)

	}
	return c.JSON(http.StatusCreated, question)
}

// ShowAllQuestions get all questions from db
func ShowAllQuestions(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var questions []Question
	db.Find(&questions)
	return c.JSON(http.StatusOK, questions)
}

// ShowQuestion get specefic question from db
func ShowQuestion(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var question Question
	id := c.Param("id")
	ObjectNotFoundError := db.Where("id=?", id).Find(&question)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, question)
}

// DeleteQuestion delete specefic question from db
func DeleteQuestion(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var question Question
	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&question).Delete(&question)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, question)
}

// UpdateQuestion update specefic question from db
func UpdateQuestion(c echo.Context) error {
	db := connectHandler.ConnectDB()
	question := new(Question)

	c.Bind(question)

	vaildate = validator.New()

	paramID := c.Param("id")

	err := vaildate.Struct(question)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"name":        question.Name,
		"degree":      question.Degree,
		"choice_type": question.ChoiceType,
		"exam_id":     question.ExamID}

	ObjectNotFoundError := db.Model(&question).Where("id= ?", paramID).Updates(attrMap)

	if ObjectNotFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectNotFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, question)
}

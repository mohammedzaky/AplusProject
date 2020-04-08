package choice

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewChoices insert new choice in db
func NewChoices(c echo.Context) error {
	db := connectHandler.ConnectDB()
	choice := new(Choice)

	c.Bind(choice)

	vaildate = validator.New()

	err := vaildate.Struct(choice)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errCreate := db.Create(&choice)
	if errCreate != nil && errCreate.RowsAffected == 0 {
		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusCreated, choice)
}

// ShowAllChoices get all choices from db
func ShowAllChoices(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var choices []Choice
	db.Find(&choices)
	return c.JSON(http.StatusOK, choices)
}

// ShowChoice get specefic chocie from db
func ShowChoice(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var choice Choice
	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&choice)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, choice)
}

// DeleteChoices delete specefic choice from db
func DeleteChoices(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var choice Choice
	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&choice).Delete(&choice)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, choice)
}

// UpdateChoices update specefic chocie from db
func UpdateChoices(c echo.Context) error {
	db := connectHandler.ConnectDB()
	choice := new(Choice)

	c.Bind(choice)

	vaildate = validator.New()

	paramID := c.Param("id")

	err := vaildate.Struct(choice)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)

	}

	attrMap := map[string]interface{}{
		"name":        choice.Name,
		"is_correct":  choice.IsCorrect,
		"question_id": choice.QuestionID,
	}

	ObjectNotFoundError := db.Model(&choice).Where("id= ?", paramID).Updates(attrMap)

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
	return c.JSON(http.StatusOK, choice)
}

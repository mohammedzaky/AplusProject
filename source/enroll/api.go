package enroll

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

// NewEnrolls insert new stundet id and course id  in db
func NewEnrolls(c echo.Context) error {
	db := connectHandler.ConnectDB()
	enrollStudent := new(Enrollment)

	c.Bind(enrollStudent)

	vaildate = validator.New()

	err := vaildate.Struct(enrollStudent)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errCreate := db.Create(&enrollStudent)

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
	return c.JSON(http.StatusCreated, enrollStudent)
}

// ShowAllEnrolls show all enrolled students
func ShowAllEnrolls(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var enrollStudent []Enrollment
	db.Find(&enrollStudent)
	return c.JSON(http.StatusOK, enrollStudent)
}

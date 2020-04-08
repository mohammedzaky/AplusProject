package studentAnswer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewStudentAnswer insert new studentAnswer in db
func NewStudentAnswer(c echo.Context) error {
	db := connectHandler.ConnectDB()

	var studentAnswers []StudentAnswer

	//Get The Request body of studentAnswers
	body, errRequest := ioutil.ReadAll(c.Request().Body)
	if errRequest != nil {
		result := map[string]string{
			"message": errRequest.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errMarashal := json.Unmarshal([]byte(body), &studentAnswers)
	if errMarashal != nil {
		result := map[string]string{
			"message": errMarashal.Error(),
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	vaildate = validator.New()
	for i := 0; i < len(studentAnswers); i++ {

		errValidate := vaildate.Struct(studentAnswers[i])

		// return validation error from Front End if there is an error in  bind
		if errValidate != nil {
			result := map[string]string{
				"message": errValidate.Error(),
			}
			return c.JSON(http.StatusBadRequest, result)
		}

		errCreate := db.Create(&studentAnswers[i])

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
	}

	return c.JSON(http.StatusCreated, studentAnswers)

}

// ShowAllStudentAnswers show all students answers
func ShowAllStudentAnswers(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var studentAnswer []StudentAnswer
	db.Find(&studentAnswer)
	return c.JSON(http.StatusOK, studentAnswer)
}

// GetAllStudentAnswer method get all answers of student for speceifc exam by studentID
func GetAllStudentAnswer(c echo.Context) error {
	db := connectHandler.ConnectDB()
	studentID := c.Param("studentID")
	var studentAnswer []StudentAnswer

	examID := c.Param("examID")

	ObjectFoundError := db.Raw("SELECT * FROM student_answers WHERE  student_id = ? and question_id in (select id from questions where exam_id = ?) ", studentID, examID).Scan(&studentAnswer)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, studentAnswer)

}

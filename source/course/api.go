package course

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	enrollObject "gitlab.com/mohamedzaky/aplusProject/source/enroll"
	examObject "gitlab.com/mohamedzaky/aplusProject/source/exam"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

//ShowAllCourses get all courses
func ShowAllCourses(c echo.Context) error {
	db := connectHandler.ConnectDB()

	professorValue := c.QueryParam("professor")
	semesterValue := c.QueryParam("semester")
	var courses []Course

	switch {

	case professorValue == "null" && semesterValue == "null":
		db.Where(map[string]interface{}{"professor_id": nil, "semester_id": nil}).Find(&courses)
		break

	case professorValue != "null" && semesterValue == "null":
		db.Where(map[string]interface{}{"professor_id": professorValue, "semester_id": nil}).Find(&courses)
		break

	case professorValue == "null" && semesterValue != "null":
		db.Where(map[string]interface{}{"professor_id": nil, "semester_id": semesterValue}).Find(&courses)
		break

	case professorValue == "" && semesterValue == "":
		db.Find(&courses)
		break

	// in Default this case professorValue != "null" && semesterValue != "null":
	default:
		db.Where(map[string]interface{}{"professor_id": professorValue, "semester_id": semesterValue}).Find(&courses)
		break
	}

	return c.JSON(http.StatusOK, courses)
}

//NewCourses insert new course
func NewCourses(c echo.Context) error {

	db := connectHandler.ConnectDB()

	course := new(Course)

	c.Bind(course)

	vaildate = validator.New()

	err := vaildate.Struct(course)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	errCreate := db.Create(&course)

	//if there no semester or prof id in db
	//or Duplicate row in courses (course_name,semester_id)
	if errCreate.RowsAffected == 0 {
		//make serialization to get specefic entity for error back
		//from db from struct ErrorModel
		//get the error in json format from db by marshaling it and unmarshaling

		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)

	}
	return c.JSON(http.StatusCreated, course)
}

//ShowCourse get specefcic course
func ShowCourse(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var course Course
	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&course)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, course)
}

//UpdateCourse update specefic course
func UpdateCourse(c echo.Context) error {

	db := connectHandler.ConnectDB()
	course := new(Course)

	c.Bind(course)

	vaildate = validator.New()

	paramID := c.Param("id")

	err := vaildate.Struct(course)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)

	}

	//Check first if there is the id valid for this Course or not
	//If true get the id of this Course
	//If not return specefic message

	checkCourse := new(Course)

	ObjectNotFoundError := db.Where("id=?", paramID).Find(&checkCourse)

	if ObjectNotFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"name":         course.Name,
		"professor_id": course.ProfessorID.Int64,
		"semester_id":  course.SemesterID.Int64,
	}

	ObjectFoundError := db.Model(&course).Where("id= ?", paramID).Updates(attrMap)

	if ObjectFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, course)
}

//DeleteCourse delete specefic course
func DeleteCourse(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var course Course
	id := c.Param("id")
	ObjectFoundError := db.Where("id=?", id).Find(&course).Delete(&course)
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, course)
}

// GetAllExams get all exams of specefic course
func GetAllExams(c echo.Context) error {
	db := connectHandler.ConnectDB()
	courseID := c.Param("id")

	var exams []examObject.Exam

	ObjectFoundError := db.Where("course_id=?", courseID).Find(&exams)

	// Return Empty array
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, exams)

}

//ShowAllStudentIds show all students id of specefic course_id
func ShowAllStudentIds(c echo.Context) error {
	db := connectHandler.ConnectDB()

	courseID := c.Param("id")

	var enrollStudent []enrollObject.Enrollment

	ObjectFoundError := db.Where("course_id=?", courseID).Find(&enrollStudent)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, enrollStudent)

}

// GetExamCourse get exam info for specefic course_ID
func GetExamCourse(c echo.Context) error {
	db := connectHandler.ConnectDB()
	courseID := c.Param("id")
	var exam examObject.Exam
	ObjectFoundError := db.Where("course_id=?", courseID).Find(&exam)
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, exam)

}

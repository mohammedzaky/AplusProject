package semester

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	courseObject "gitlab.com/mohamedzaky/aplusProject/source/course"
	professorObject "gitlab.com/mohamedzaky/aplusProject/source/professor"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// ShowAllSemesters to show all semesters from db
func ShowAllSemesters(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var semesters []Semester
	db.Find(&semesters)
	return c.JSON(http.StatusOK, semesters)
}

// ShowSemester to show semester from db
func ShowSemester(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var semester Semester
	id := c.Param("id")

	//check here if there is id for this semester if not return error
	ObjectFoundError := db.Where("id=?", id).Find(&semester)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, semester)
}

// NewSemesters to insert a new semester in db
func NewSemesters(c echo.Context) error {

	db := connectHandler.ConnectDB()

	semester := new(Semester)

	c.Bind(semester)

	vaildate = validator.New()

	err := vaildate.Struct(semester)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)

	}

	//check here if err happen in create a row in db
	//Example unique index error
	errCreate := db.Create(&semester)

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

	return c.JSON(http.StatusCreated, semester)
}

// UpdateSemester to update a semester in db
func UpdateSemester(c echo.Context) error {

	db := connectHandler.ConnectDB()

	semester := new(Semester)

	c.Bind(semester)

	vaildate = validator.New()

	paramID := c.Param("id")

	err := vaildate.Struct(semester)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)

	}
	//Check first if there is the id valid for this Semester or not
	//If not return specefic message
	checksemester := new(Semester)
	ObjectFoundError := db.Where("id=?", paramID).Find(checksemester)

	if ObjectFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"name": semester.Name,
		"year": semester.Year,
	}

	ObjectFoundError = db.Model(&semester).Where("id= ?", paramID).Updates(attrMap)

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
	return c.JSON(http.StatusOK, semester)
}

// DeleteSemester to delete a semester from db
func DeleteSemester(c echo.Context) error {

	db := connectHandler.ConnectDB()
	var semester Semester

	id := c.Param("id")

	ObjectFoundError := db.Where("id=?", id).Find(&semester).Delete(&semester)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, semester)
}

// GetAllCourses Show all courses in specefic semester
func GetAllCourses(c echo.Context) error {

	db := connectHandler.ConnectDB()

	id := c.Param("id")

	var courses []courseObject.Course

	// Execute Query and Get all course object in courses
	ObjectFoundError := db.Where("semester_id=?", id).Find(&courses)

	// Return Empty array
	if ObjectFoundError.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, courses)
	}

	size := len(courses)

	//Get professor id from table professors
	var professorIds []int64

	//get all professorIds for each student in slice userIds
	for i := 0; i < size; i++ {
		professorIds = append(professorIds, courses[i].ProfessorID.Int64)
	}

	//get all userid from professor table based on professorids slice

	var userIds []int
	professors := make([]professorObject.Professor, len(professorIds))

	for i := 0; i < size; i++ {

		ObjectFoundError := db.Where("id=?", professorIds[i]).Find(&professors[i])

		if ObjectFoundError.RowsAffected == 0 {
			result := map[string]string{
				"message": "No Record Found in DB for this id",
			}
			return c.JSON(http.StatusNotFound, result)
		}

		userIds = append(userIds, professors[i].UserID)
	}

	//get all professor names from user table based on userid slice
	users := make([]userObject.User, len(userIds))

	fullname := make([]FullName, len(userIds))

	for i := 0; i < size; i++ {

		ObjectFoundError := db.Where("id=?", userIds[i]).Find(&users[i])

		if ObjectFoundError.RowsAffected == 0 {
			result := map[string]string{
				"message": "No Record Found in DB for this id",
			}
			return c.JSON(http.StatusNotFound, result)
		}
		fullname[i].firstName = users[i].FirstName
		fullname[i].lastName = users[i].LastName
	}

	// map of courseResponse object of size courses
	courseResponse := make([]courseObject.CourseSerializer, size)

	for i := 0; i < size; i++ {

		courseResponse[i].ID = courses[i].ID
		courseResponse[i].CourseName = courses[i].Name
		courseResponse[i].ProfessorName = fullname[i].firstName + " " + fullname[i].lastName

	}
	return c.JSON(http.StatusOK, courseResponse)
}

func newCourseInSemester(c echo.Context) error {
	db := connectHandler.ConnectDB()
	course := new(courseObject.Course)

	c.Bind(course)

	vaildate = validator.New()

	err := vaildate.Struct(course)
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)

	}

	attrMap := map[string]interface{}{
		"name":         course.Name,
		"professor_id": course.ProfessorID.Int64,
		"semester_id":  course.SemesterID.Int64,
	}

	//Get first id from course name
	checkCourse := new(courseObject.Course)

	errFind := db.Where("name = ?", course.Name).Find(&checkCourse)
	if errFind.RowsAffected == 0 {
		return c.JSON(http.StatusBadRequest, errFind)
	}

	ObjectFoundError := db.Model(&course).Where("id = ?", checkCourse.ID).Updates(attrMap)

	if ObjectFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)
		//get detail and message from the error back
		result := map[string]string{
			"message1": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	return c.JSON(http.StatusOK, course)
}

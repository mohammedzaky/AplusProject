package student

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	enrollObject "gitlab.com/mohamedzaky/aplusProject/source/enroll"
	Hash "gitlab.com/mohamedzaky/aplusProject/source/hashHelper"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewStudent insert a new student to db
func NewStudent(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(userObject.User)
	student := new(Student)

	//create a studentRequest for request API
	studentRequest := new(StudentRequest)

	// create studentResponse for response in return API
	studentResponse := new(StudentResponse)

	c.Bind(studentRequest)

	vaildate = validator.New()

	err := vaildate.Struct(studentRequest)

	// return validation error from Front End if there is an error in bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//psHash to hasing the password and insert it into DB

	password := studentRequest.Password
	passHash, _ := Hash.DoHash(password)

	//Insert the value in user object
	user.FirstName = studentRequest.FirstName
	user.LastName = studentRequest.LastName
	user.UserName = studentRequest.UserName

	//Put the hased password in user object
	user.Password = passHash

	user.Phone = studentRequest.Phone

	//Create a User Object here
	errCreate := db.Create(&user)

	//Create first a user of hased password
	if errCreate.RowsAffected == 0 {
		//make serialization to get specefic entity for error back from db from struct ErrorModel
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
	student.Gpa = studentRequest.Gpa
	student.Hours = studentRequest.Hours
	//student.SeatNumber = studentRequest.SeatNumber
	student.UserID = user.ID

	//create a student object in db
	errCreate = db.Create(&student)

	//Create a student his/her user_id the id of specefic created user
	if errCreate.RowsAffected == 0 {
		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)
		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//Update seat_number for this student
	//SeatNumber is combination of Entry Year + StudentID

	t := time.Now()

	year := strings.SplitAfter(t.String(), "-")

	yearNoDash := strings.Replace(year[0], "-", "", -1)

	studentID := strconv.Itoa(student.ID)

	student.SeatNumber.String = yearNoDash + studentID

	db.Model(&student).Where("id= ?", student.ID).Update("seat_number", student.SeatNumber.String)

	//insert a student response for new student API
	studentResponse.StudentID = student.ID
	studentResponse.SeatNumber = student.SeatNumber.String
	studentResponse.StudentHours = student.Hours
	studentResponse.StudentGPA = student.Gpa
	studentResponse.User.ID = user.ID
	studentResponse.User.FirstName = user.FirstName
	studentResponse.User.LastName = user.LastName
	studentResponse.User.Username = user.UserName
	studentResponse.User.Telephone = user.Phone

	return c.JSON(http.StatusCreated, studentResponse)
}

// ShowAllStudents get all students from db
func ShowAllStudents(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var student []Student

	db.Find(&student)

	var userIds []int
	//get  all userid for each student in slice userIds
	for i := 0; i < len(student); i++ {
		userIds = append(userIds, student[i].UserID)
	}

	user := make([]userObject.User, len(userIds))

	//get all user info from user table based on user id from student table on student_id
	for i := 0; i < len(userIds); i++ {

		if ObjectNotFoundError := db.Where("id=?", userIds[i]).Find(&user[i]); ObjectNotFoundError.RowsAffected == 0 {

			result := map[string]string{
				"message": "No Record Found in DB for this id",
			}
			return c.JSON(http.StatusNotFound, result)
		}
	}
	//create and get all user info and student info in studentResponse Struct back in API response
	studentResponse := make([]StudentResponse, len(userIds))
	for i := 0; i < len(userIds); i++ {
		studentResponse[i].StudentID = student[i].ID
		studentResponse[i].StudentGPA = student[i].Gpa
		studentResponse[i].StudentHours = student[i].Hours
		studentResponse[i].User.ID = user[i].ID
		studentResponse[i].User.FirstName = user[i].FirstName
		studentResponse[i].User.LastName = user[i].LastName
		studentResponse[i].User.Username = user[i].UserName
		studentResponse[i].User.Telephone = user[i].Phone
	}

	return c.JSON(http.StatusOK, studentResponse)

}

// ShowStudent get specefic student from db
func ShowStudent(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var student Student
	var user userObject.User

	id := c.Param("id")

	//check if there id is valid or not
	ObjectNotFoundError := db.Where("id=?", id).Find(&student)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	usrID := student.UserID

	ObjectNotFoundError = db.Where("id=?", usrID).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//create a response for studentAPI

	var studentResponse StudentResponse
	studentResponse.StudentID = student.ID
	studentResponse.StudentGPA = student.Gpa
	studentResponse.StudentHours = student.Hours
	studentResponse.SeatNumber = student.SeatNumber.String
	studentResponse.User.ID = user.ID
	studentResponse.User.FirstName = user.FirstName
	studentResponse.User.LastName = user.LastName
	studentResponse.User.Username = user.UserName
	studentResponse.User.Telephone = user.Phone

	return c.JSON(http.StatusOK, studentResponse)
}

// UpdateStudent update specefic student from db
func UpdateStudent(c echo.Context) error {
	db := connectHandler.ConnectDB()
	studentUpdate := new(StudentUpdate)
	user := new(userObject.User)
	student := new(Student)

	c.Bind(studentUpdate)

	vaildate = validator.New()

	studentID := c.Param("id")

	err := vaildate.Struct(studentUpdate)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//Check first if there is the id valid for this student or not
	//If true get the user_id of this student
	//If not return specefic message

	checkStudent := new(Student)
	ObjectNotFoundError := db.Where("id=?", studentID).Find(checkStudent)

	if ObjectNotFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	//Get the student info from studentupdate object and insert it into student table
	student.Gpa = studentUpdate.Gpa
	student.Hours = studentUpdate.Hours
	//student.SeatNumber = studentUpdate.SeatNumber

	//First update the student info from student table
	attrStudent := map[string]interface{}{
		"gpa":   student.Gpa,
		"hours": student.Hours,
	}

	ObjectNotFoundError = db.Model(&student).Where("id= ?", studentID).Updates(attrStudent)

	if ObjectNotFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectNotFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//Second Update his/her info from user table
	//First Get the user id from student table to update user row in user table for this specefic id
	db.Where("user_id=?", student.UserID)

	//Get in user info from studentupdate object and insert it into user table
	user.FirstName = studentUpdate.FirstName
	user.LastName = studentUpdate.LastName
	user.UserName = studentUpdate.UserName
	user.Phone = studentUpdate.Telephone

	userID := checkStudent.UserID
	attrUser := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"user_name":  user.UserName,
		"phone":      user.Phone,
	}

	ObjectNotFoundError = db.Model(&user).Where("id= ?", userID).Updates(attrUser)

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

	//create a response for studentAPI

	var studentResponse StudentResponse
	studentResponse.StudentID = checkStudent.ID
	studentResponse.StudentGPA = student.Gpa
	studentResponse.StudentHours = student.Hours
	studentResponse.SeatNumber = student.SeatNumber.String
	student.SeatNumber.Valid = true
	studentResponse.User.ID = userID
	studentResponse.User.FirstName = user.FirstName
	studentResponse.User.LastName = user.LastName
	studentResponse.User.Username = user.UserName
	studentResponse.User.Telephone = user.Phone

	return c.JSON(http.StatusOK, studentResponse)

}

//DeleteStudent delete student row and his/her info
func DeleteStudent(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user userObject.User
	var student Student
	id := c.Param("id")
	//Check first if there is the id valid for this student or not
	//If true get the user_id of this student
	//If not return specefic message

	ObjectNotFoundError := db.Where("id=?", id).Find(&student)
	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}
	user.ID = student.UserID

	//Remove student row first
	ObjectNotFoundError = db.Where("id=?", id).Find(&student).Delete(&student)
	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//Remove user row
	ObjectNotFoundError = db.Where("id=?", user.ID).Find(&user).Delete(&user)
	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, user)

}

//ShowAllCourseIds show all course id of specefic student_id
func ShowAllCourseIds(c echo.Context) error {
	db := connectHandler.ConnectDB()

	studentID := c.Param("id")

	var enrollStudent []enrollObject.Enrollment

	ObjectFoundError := db.Where("student_id=?", studentID).Find(&enrollStudent)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, enrollStudent)

}

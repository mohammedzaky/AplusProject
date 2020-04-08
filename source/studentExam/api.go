package studentExam

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	choiceObject "gitlab.com/mohamedzaky/aplusProject/source/choice"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	courseObject "gitlab.com/mohamedzaky/aplusProject/source/course"
	examObject "gitlab.com/mohamedzaky/aplusProject/source/exam"
	questionObject "gitlab.com/mohamedzaky/aplusProject/source/question"
	studentObject "gitlab.com/mohamedzaky/aplusProject/source/student"
	studentAnswerObject "gitlab.com/mohamedzaky/aplusProject/source/studentAnswer"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

//NewStudentExam method
func NewStudentExam(c echo.Context) error {
	db := connectHandler.ConnectDB()
	studentExam := new(StudentExam)

	c.Bind(studentExam)
	vaildate = validator.New()

	err := vaildate.Struct(studentExam)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	errCreate := db.Create(&studentExam)

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
	return c.JSON(http.StatusCreated, studentExam)

}

//ShowAllStudentExams method
func ShowAllStudentExams(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var studentExam []StudentExam
	db.Find(&studentExam)
	return c.JSON(http.StatusOK, studentExam)
}

// UpdateStudentExam update specefic studentExam from db
func UpdateStudentExam(c echo.Context) error {
	db := connectHandler.ConnectDB()
	studentExam := new(StudentExam)

	c.Bind(studentExam)

	vaildate = validator.New()

	err := vaildate.Struct(studentExam)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"student_id":     studentExam.StudentID,
		"exam_id":        studentExam.ExamID,
		"student_degree": studentExam.StudentDegree,
	}

	ObjectNotFoundError := db.Model(&studentExam).Where("student_id= ? and exam_id = ?", studentExam.StudentID, studentExam.ExamID).Updates(attrMap)

	if ObjectNotFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectNotFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		if dbErr.Detail == "" {
			result := map[string]string{
				"message": "No Record Found in DB for this id",
			}
			return c.JSON(http.StatusBadRequest, result)
		}
		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	return c.JSON(http.StatusOK, studentExam)
}

//NewStudentDegree after he/she examed
//insert & calculate his degree in student_exam table
//compare the student answers with the model anwser (studentExam)
func NewStudentDegree(c echo.Context) error {
	db := connectHandler.ConnectDB()
	studentID := c.Param("studentID")
	examID := c.Param("examID")

	var studentAnswer []studentAnswerObject.StudentAnswer
	var studentExam StudentExam
	var exam examObject.Exam
	var student studentObject.Student

	//check first if there is a valid id of this student or not
	ObjectFoundError := db.Where("id=?", studentID).Find(&student)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//check first if there is a valid id of this exam or not
	ObjectFoundError = db.Where("id=?", examID).Find(&exam)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	//select here the all student answers from speceifc exam
	db.Raw("select * from student_answers where question_id in (select id from questions where exam_id = ?)", 5).Scan(&studentAnswer)

	// studentDegree for speceifc exam
	var studentDegree float32

	var questionResponse questionObject.Question
	var choice choiceObject.Choice

	//Get question degree for each question in studentAnswer array
	//Check first if student answer this question or leave it empty
	for i := 0; i < len(studentAnswer); i++ {
		db.Select("degree").Where("id = ?", studentAnswer[i].QuestionID).Find(&questionResponse)

		if studentAnswer[i].StudentChoice.Int64 == 0 {
			//these mean that the student not answer these question and the value
			//of this studentAnswer is null
			studentDegree -= questionResponse.Degree
		} else {
			//check here if the choice of student is correct
			//(this means that choice_id of correct answer is equal to choice_id of student_answer table)
			//or the choose was wrong choice answer

			db.Where("question_id=? and is_correct = ?", studentAnswer[i].QuestionID, true).Find(&choice)
			if studentAnswer[i].StudentChoice.Int64 == (int64)(choice.ID) {
				studentDegree += questionResponse.Degree

			} else {
				studentDegree -= questionResponse.Degree
			}
		}
		if studentDegree < 0 {
			studentDegree = 0
		}

		// set the values of choice struct to default for the second iteration
		choice.ID = 0
		choice.IsCorrect = false
		choice.Name = ""
		choice.QuestionID = 0
	}

	//get degree for this exam
	//index of 0 here or any index because there all forward to one exam
	db.Where("id = ?", studentAnswer[0].QuestionID).Find(&exam)

	studentExam.ExamID = exam.ID
	studentExam.StudentDegree.Valid = true
	studentExam.StudentDegree.Float64 = (float64)(studentDegree)
	studentId, _ := strconv.Atoi(studentID)
	studentExam.StudentID = studentId

	errCreate := db.Create(&studentExam)

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
	return c.JSON(http.StatusOK, studentExam)

}

//GetStudentDegree after he/she examed
//Get only one object of student degree
//calcuate the result and get it
func GetStudentDegree(c echo.Context) error {
	db := connectHandler.ConnectDB()

	studentID := c.Param("studentID")
	examID := c.Param("examID")

	var studentExam StudentExam

	ObjectNotFoundError := db.Raw("select * from student_exams where student_id = ? and exam_id = ?", studentID, examID).Scan(&studentExam)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, studentExam)
}

//GetStudentDegrees method for only one exam
func GetStudentDegrees(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var studentExam []StudentExam
	var exam examObject.Exam

	examID := c.Param("id")

	objectFoundError := db.Where("exam_id=?", examID).Find(&studentExam)

	if objectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	examResultSlice := make(ExamResult, len(studentExam))
	student := make([]studentObject.Student, len(studentExam))
	user := make([]userObject.User, len(studentExam))

	//Get studnet name for each student id and exam name
	for i := 0; i < len(studentExam); i++ {

		db.Where("id = ?", studentExam[i].StudentID).Find(&student[i])
		db.Where("id = ?", student[i].UserID).Find(&user[i])
		db.Where("id = ?", examID).Find(&exam)

		examResultSlice[i].StudentName = user[i].FirstName + " " + user[i].LastName
		examResultSlice[i].StudentSeatNumber = student[i].SeatNumber.String

		degrees := make([]Degrees, 1)

		for j := 0; j < 1; j++ {
			degrees[j].ExamName = exam.Name
			degrees[j].StudentDegree = studentExam[i].StudentDegree
		}
		examResultSlice[i].ExamDegrees = degrees

	}

	return c.JSON(http.StatusOK, examResultSlice)
}

//ResetStudentDegrees method for only one exam set
func ResetStudentDegrees(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var studentExam []StudentExam
	var exam examObject.Exam

	examID := c.Param("id")

	objectFoundError := db.Where("exam_id=?", examID).Find(&studentExam)

	if objectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	// Set here all student degrees from studentExam table to null
	for i := 0; i < len(studentExam); i++ {

		//valid will set that value to null if it was false
		studentExam[i].StudentDegree.Valid = false
		db.Model(&studentExam).Where("id = ?", studentExam[i].ID).Update("student_degree", studentExam[i].StudentDegree)

	}
	examResultSlice := make(ExamResult, len(studentExam))
	student := make([]studentObject.Student, len(studentExam))
	user := make([]userObject.User, len(studentExam))

	//Get studnet name for each student id and exam name
	for i := 0; i < len(studentExam); i++ {

		db.Where("id = ?", studentExam[i].StudentID).Find(&student[i])
		db.Where("id = ?", student[i].UserID).Find(&user[i])
		db.Where("id = ?", examID).Find(&exam)

		examResultSlice[i].StudentName = user[i].FirstName + " " + user[i].LastName
		examResultSlice[i].StudentSeatNumber = student[i].SeatNumber.String

		degrees := make([]Degrees, 1)

		for j := 0; j < 1; j++ {
			degrees[j].ExamName = exam.Name
			degrees[j].StudentDegree = studentExam[i].StudentDegree
		}
		examResultSlice[i].ExamDegrees = degrees

	}

	return c.JSON(http.StatusOK, examResultSlice)
}

// GetAllStudentDegrees method for all exams in specefic course
func GetAllStudentDegrees(c echo.Context) error {
	db := connectHandler.ConnectDB()

	var studentExam []StudentExam
	var course courseObject.Course

	courseID := c.Param("id")

	courseFoundError := db.Where("id=?", courseID).Find(&course)

	if courseFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	var exams []examObject.Exam

	examFoundError := db.Where("course_id=?", courseID).Find(&exams)

	if examFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id or no exams found for this course",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	examResults := make(ExamResults, len(exams))

	//Geting the length of studentExam
	for i := 0; i < len(exams); i++ {
		//get all students that enter that exam and there degrees
		// in studentExam object
		db.Where("exam_id = ?", exams[i].ID).Find(&studentExam)
	}

	student := make([]studentObject.Student, len(studentExam))
	user := make([]userObject.User, len(studentExam))
	studentDegrees := make([]StudentDegrees, len(studentExam))
	studentDegreesSlice := make([]StudentDegrees, 0)

	for i := 0; i < len(exams); i++ {

		examResults[i].ExamName = exams[i].Name
		db.Where("exam_id = ?", exams[i].ID).Find(&studentExam)

		for j := 0; j < len(studentExam); j++ {

			db.Where("id = ?", studentExam[j].StudentID).Find(&student[j])

			db.Where("id = ?", student[j].UserID).Find(&user[j])

			studentDegrees[j].StudentSeatNumber = student[j].SeatNumber.String

			studentDegrees[j].StudentName = user[j].FirstName + " " + user[j].LastName
			studentDegrees[j].StudentDegree = studentExam[j].StudentDegree

			studentDegreesSlice = append(studentDegreesSlice, studentDegrees[j])

		}

		examResults[i].StudentDegrees = studentDegreesSlice
		// Empty studentDegreesSlice for new exam in another itration

		studentDegreesSlice = make([]StudentDegrees, 0)

	}
	return c.JSON(http.StatusOK, examResults)

}

//ResetAllStudentDegrees for all exams in specefic course
func ResetAllStudentDegrees(c echo.Context) error {
	db := connectHandler.ConnectDB()

	var studentExam []StudentExam
	var course courseObject.Course

	courseID := c.Param("id")

	courseFoundError := db.Where("id=?", courseID).Find(&course)

	if courseFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	var exams []examObject.Exam

	examFoundError := db.Where("course_id=?", courseID).Find(&exams)

	if examFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id or no exams found for this course",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	examResults := make(ExamResults, len(exams))

	//Geting the length of studentExam
	for i := 0; i < len(exams); i++ {
		//get all students that enter that exam and there degrees
		// in studentExam object and Reset all degrees
		db.Where("exam_id = ?", exams[i].ID).Find(&studentExam)
		for j := 0; j < len(studentExam); j++ {

			studentExam[j].StudentDegree.Valid = false
			db.Model(&studentExam).Where("id = ?", studentExam[j].ID).Update("student_degree", studentExam[j].StudentDegree)

		}
	}

	student := make([]studentObject.Student, len(studentExam))
	user := make([]userObject.User, len(studentExam))
	studentDegrees := make([]StudentDegrees, len(studentExam))
	studentDegreesSlice := make([]StudentDegrees, 0)

	for i := 0; i < len(exams); i++ {

		examResults[i].ExamName = exams[i].Name
		db.Where("exam_id = ?", exams[i].ID).Find(&studentExam)

		for j := 0; j < len(studentExam); j++ {

			db.Where("id = ?", studentExam[j].StudentID).Find(&student[j])

			db.Where("id = ?", student[j].UserID).Find(&user[j])

			studentDegrees[j].StudentSeatNumber = student[j].SeatNumber.String

			studentDegrees[j].StudentName = user[j].FirstName + " " + user[j].LastName
			studentDegrees[j].StudentDegree = studentExam[j].StudentDegree

			studentDegreesSlice = append(studentDegreesSlice, studentDegrees[j])

		}

		examResults[i].StudentDegrees = studentDegreesSlice
		// Empty studentDegreesSlice for new exam in another itration

		studentDegreesSlice = make([]StudentDegrees, 0)

	}
	return c.JSON(http.StatusOK, examResults)
}

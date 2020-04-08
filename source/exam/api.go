package exam

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	choiceObject "gitlab.com/mohamedzaky/aplusProject/source/choice"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	questionObject "gitlab.com/mohamedzaky/aplusProject/source/question"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewExams insert new Exam in db
func NewExams(c echo.Context) error {
	db := connectHandler.ConnectDB()
	exam := new(Exam)

	c.Bind(exam)

	vaildate = validator.New()

	err := vaildate.Struct(exam)

	// return validation error from Front End if there is an error in  bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errCreate := db.Create(&exam)

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

	return c.JSON(http.StatusCreated, exam)
}

// ShowAllExams get all exams from db
func ShowAllExams(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var exams []Exam
	db.Find(&exams)
	return c.JSON(http.StatusOK, exams)
}

// ShowExam get specefic exam
func ShowExam(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var exam Exam
	id := c.Param("id")
	ObjectNotFoundError := db.Where("id=?", id).Find(&exam)
	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, exam)
}

// DeleteExam delete specefic exam
func DeleteExam(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var exam Exam
	id := c.Param("id")
	ObjectNotFoundError := db.Where("id=?", id).Find(&exam).Delete(&exam)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, exam)
}

// UpdateExam update specefic exam
func UpdateExam(c echo.Context) error {
	db := connectHandler.ConnectDB()
	exam := new(Exam)

	c.Bind(exam)

	vaildate = validator.New()

	paramID := c.Param("id")

	err := vaildate.Struct(exam)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"name":      exam.Name,
		"duration":  exam.Duration,
		"degree":    exam.Degree,
		"is_enable": exam.IsEnable,
		"course_id": exam.CourseID,
	}

	ObjectNotFoundError := db.Model(&exam).Where("id= ?", paramID).Updates(attrMap)

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
	return c.JSON(http.StatusOK, exam)
}

// GetAllQuestions get all questions of specefic exam
func GetAllQuestions(c echo.Context) error {
	db := connectHandler.ConnectDB()

	examID := c.Param("id")

	var questions []questionObject.Question

	// Execute Query and Get all course object in courses
	ObjectFoundError := db.Where("exam_id=?", examID).Find(&questions)

	// Return Empty array
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, questions)

}

//GetExamModel method
func GetExamModel(c echo.Context) error {
	db := connectHandler.ConnectDB()

	examID := c.Param("id")

	var questions []questionObject.Question

	// Execute Query and Get all questions for each exam
	ObjectFoundError := db.Where("exam_id=?", examID).Find(&questions)

	// Return Empty array if there is no question for this exam or invalid id
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//get all choices info from choice table based on question id
	//get no of choices for each question by countChoices
	var countChoices int
	var count int

	for i := 0; i < len(questions); i++ {
		db.Model(&choiceObject.Choice{}).Where("question_id = ?", questions[i].ID).Count(&count)
		countChoices += count
	}

	//len of choices slice will be count of all choices
	choices := make([]choiceObject.Choice, countChoices)

	questionResponseSlice := make(QuestionResponse, len(questions))

	for i := 0; i < len(questions); i++ {
		questionResponseSlice[i].ID = questions[i].ID
		questionResponseSlice[i].Name = questions[i].Name
		questionResponseSlice[i].Degree = questions[i].Degree
		questionResponseSlice[i].ChoiceType = questions[i].ChoiceType
		questionResponseSlice[i].ExamID = questions[i].ExamID

		//Get all choices for specefic exam
		db.Where("question_id=?", questions[i].ID).Find(&choices)
		questionResponseSlice[i].Choices = choices

	}
	return c.JSON(http.StatusOK, questionResponseSlice)

}

//GetAllQuestionsChoices method for student
//Only romove is_correct field of response method
func GetAllQuestionsChoices(c echo.Context) error {
	db := connectHandler.ConnectDB()

	examID := c.Param("id")

	var questions []questionObject.Question

	// Execute Query and Get all questions for each exam
	ObjectFoundError := db.Where("exam_id=?", examID).Find(&questions)

	// Return Empty array if there is no question for this exam or invalid id
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//get all choices info from choice table based on question id
	//get no of choices for each question by countChoices
	var countChoices int
	var count int

	for i := 0; i < len(questions); i++ {
		db.Model(&choiceObject.Choice{}).Where("question_id = ?", questions[i].ID).Count(&count)
		countChoices += count
	}

	//len of choices slice will be count of all choices
	choices := make([]choiceObject.Choice, countChoices)

	questionResponseSlice := make(QuestionResponse, len(questions))

	for i := 0; i < len(questions); i++ {
		questionResponseSlice[i].ID = questions[i].ID
		questionResponseSlice[i].Name = questions[i].Name
		questionResponseSlice[i].Degree = questions[i].Degree
		questionResponseSlice[i].ChoiceType = questions[i].ChoiceType
		questionResponseSlice[i].ExamID = questions[i].ExamID

		//Get all choices for specefic exam
		//Query to select only id,name,question_id and not getting is_correct value if it true
		db.Select("id, name ,question_id").Where("question_id = ?", questions[i].ID).Find(&choices)
		questionResponseSlice[i].Choices = choices

	}
	return c.JSON(http.StatusOK, questionResponseSlice)

}

//CreateExamModel of Ask And Answers
func CreateExamModel(c echo.Context) error {

	db := connectHandler.ConnectDB()
	var questionRequest QuestionRequest
	var question questionObject.Question
	var checkQuestion questionObject.Question
	var choice choiceObject.Choice
	var checkExam Exam

	examID := c.Param("id")

	//check first if the id is valid or not
	ObjectNotFoundError := db.Where("id=?", examID).Find(&checkExam)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//Get The Request body of studentAnswers
	body, errRequest := ioutil.ReadAll(c.Request().Body)
	if errRequest != nil {
		result := map[string]string{
			"message": errRequest.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	errMarashal := json.Unmarshal([]byte(body), &questionRequest)
	if errMarashal != nil {
		result := map[string]string{
			"message": errMarashal.Error(),
		}
		return c.JSON(http.StatusInternalServerError, result)
	}

	for i := 0; i < len(questionRequest); i++ {

		question.Name = questionRequest[i].Name
		question.Degree = questionRequest[i].Degree
		question.ChoiceType = questionRequest[i].ChoiceType

		//convert from string to int
		examID, _ := strconv.Atoi(examID)

		question.ExamID = examID

		//Create Question in DB
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
		//get the id of the question because the insert of question id is 0 value
		db.Raw("select * from questions where name = ? and degree = ?", question.Name, question.Degree).Scan(&checkQuestion)

		lenChoices := len(questionRequest[i].Choices)

		for j := 0; j < lenChoices; j++ {
			choice.Name = questionRequest[i].Choices[j].Name
			choice.QuestionID = checkQuestion.ID
			choice.IsCorrect = questionRequest[i].Choices[j].IsCorrect

			//Create Choice in DB
			errCreate := db.Create(&choice)

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
			// set the values of choice struct to default for the second iteration
			// without set the value of id of choice = 0 will get error in db for Primary key dublicated
			// of zero value
			choice.ID = 0

		}

		// set the values of question struct to default for the second iteration
		question.ID = 0

	}
	return c.JSON(http.StatusCreated, questionRequest)
}

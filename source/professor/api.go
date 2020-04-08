package professor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	courseObject "gitlab.com/mohamedzaky/aplusProject/source/course"
	Hash "gitlab.com/mohamedzaky/aplusProject/source/hashHelper"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

//NewProfessor Struct
func NewProfessor(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(userObject.User)
	professor := new(Professor)

	// create professorRequest for request API
	professorRequest := new(ProfessorRequest)

	// create professorResponse for response in return API
	professorResponse := new(ProfessorResponse)

	c.Bind(professorRequest)

	vaildate = validator.New()

	err := vaildate.Struct(professorRequest)

	// return validation error from Front End if there is an error in bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//psHash to hasing the password and insert it into DB

	password := professorRequest.Password
	passHash, _ := Hash.DoHash(password)

	//Insert the value in user object
	user.FirstName = professorRequest.FirstName
	user.LastName = professorRequest.LastName
	user.UserName = professorRequest.UserName

	//Put the hased password in user object
	user.Password = passHash

	user.Phone = professorRequest.Phone

	//Create a User Object here
	errCreate := db.Create(&user)

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

	professor.Degree = professorRequest.Degree
	professor.Major = professorRequest.Major
	professor.UserID = user.ID

	//create a professor object in db
	errCreate = db.Create(&professor)

	if errCreate.RowsAffected == 0 {
		//make serialization to get specefic entity for error back from db from struct ErrorModel
		//get the error in json format from db by marshaling it and unmarshaling

		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	// insert a professor response for new user API
	professorResponse.ProfessorID = professor.ID
	professorResponse.Degree = professor.Degree
	professorResponse.Major = professor.Major
	professorResponse.ImageURL = professor.ImageURL
	professorResponse.User.ID = user.ID
	professorResponse.User.FirstName = user.FirstName
	professorResponse.User.LastName = user.LastName
	professorResponse.User.UserName = user.UserName
	professorResponse.User.Phone = user.Phone

	return c.JSON(http.StatusCreated, professorResponse)
}

//ShowProfessor Struct
func ShowProfessor(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user userObject.User
	var professor Professor
	//var userResponse userObject.UserResponse
	var professorResponse ProfessorResponse
	professorID := c.Param("id")

	//check if the professor id is valid or not
	//And get the professor id and user id
	ObjectNotFoundError := db.Where("id=?", professorID).Find(&professor)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	//Get the user info from db based on user id from professor table

	ObjectNotFoundError = db.Where("id=?", professor.UserID).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//create a response for userAPI

	professorResponse.ProfessorID = professor.ID
	professorResponse.Degree = professor.Degree
	professorResponse.Major = professor.Major
	professorResponse.ImageURL = professor.ImageURL
	professorResponse.User.ID = user.ID
	professorResponse.User.FirstName = user.FirstName
	professorResponse.User.LastName = user.LastName
	professorResponse.User.UserName = user.UserName
	professorResponse.User.Phone = user.Phone

	return c.JSON(http.StatusOK, professorResponse)
}

//ShowAllProfessors Struct
func ShowAllProfessors(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var professor []Professor

	db.Find(&professor)

	var userIds []int
	//get all userid for each prfessor in slice userIds
	for i := 0; i < len(professor); i++ {
		userIds = append(userIds, professor[i].UserID)
	}

	user := make([]userObject.User, len(userIds))

	//get all user info from user table based on user id from professor table
	for i := 0; i < len(userIds); i++ {

		ObjectNotFoundError := db.Where("id=?", userIds[i]).Find(&user[i])

		if ObjectNotFoundError.RowsAffected == 0 {

			result := map[string]string{
				"message": "No Record Found in DB for this id1",
			}
			return c.JSON(http.StatusNotFound, result)
		}
	}
	//create and get all user info in userResponse Struct back in API response

	professorResponse := make([]ProfessorResponse, len(userIds))

	for i := 0; i < len(userIds); i++ {
		professorResponse[i].Degree = professor[i].Degree
		professorResponse[i].Major = professor[i].Major
		professorResponse[i].ImageURL = professor[i].ImageURL
		professorResponse[i].ProfessorID = professor[i].ID
		professorResponse[i].User.ID = user[i].ID
		professorResponse[i].User.FirstName = user[i].FirstName
		professorResponse[i].User.LastName = user[i].LastName
		professorResponse[i].User.UserName = user[i].UserName
		professorResponse[i].User.Phone = user[i].Phone
	}
	return c.JSON(http.StatusOK, professorResponse)
}

//DeleteProfessors Struct
func DeleteProfessors(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user userObject.User
	var professor Professor
	id := c.Param("id")
	//Check first if there is the id valid for this professor or not
	//If true get the user_id of this professor
	//If not return specefic message

	ObjectNotFoundError := db.Where("id=?", id).Find(&professor)
	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	user.ID = professor.UserID

	//Remove professor row first
	ObjectNotFoundError = db.Where("id=?", id).Find(&professor).Delete(&professor)

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

//UpdateProfessors Struct
func UpdateProfessors(c echo.Context) error {
	db := connectHandler.ConnectDB()
	professorUpdate := new(ProfessorUpdate)
	professor := new(Professor)
	professorResponse := new(ProfessorResponse)
	user := new(userObject.User)

	c.Bind(professorUpdate)

	vaildate = validator.New()

	professorID := c.Param("id")

	err := vaildate.Struct(professorUpdate)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//Check first if there is the id valid for this Professor or not
	//If true get the user_id of this Professor
	//If not return specefic message

	var checkProfessor Professor

	ObjectNotFoundError := db.Where("id=?", professorID).Find(&checkProfessor)

	if ObjectNotFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	//Get the professor info from professorUpdate object and insert it into professor table
	professor.Degree = professorUpdate.Degree
	professor.Major = professorUpdate.Major

	// update the professor info from professor table

	attrProfessor := map[string]interface{}{
		"degree": professor.Degree,
		"major":  professor.Major,
	}

	ObjectNotFoundError = db.Model(&professor).Where("id= ?", professorID).Updates(attrProfessor)

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
	//First Get the user id from professor table to update user row in user table for this specefic id
	db.Where("user_id=?", professor.UserID)

	//Get in user info from studentupdate object and insert it into user table
	user.FirstName = professorUpdate.FirstName
	user.LastName = professorUpdate.LastName
	user.UserName = professorUpdate.UserName
	user.Phone = professorUpdate.Phone

	userID := checkProfessor.UserID

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

	professorResponse.Degree = professor.Degree
	professorResponse.Major = professor.Major
	professorResponse.ImageURL = professor.ImageURL
	professorResponse.ProfessorID = checkProfessor.ID
	professorResponse.User.ID = checkProfessor.UserID
	professorResponse.User.FirstName = user.FirstName
	professorResponse.User.LastName = user.LastName
	professorResponse.User.UserName = user.UserName
	professorResponse.User.Phone = user.Phone

	return c.JSON(http.StatusOK, professorResponse)

}

//UploadImage Struct
func UploadImage(c echo.Context) error {

	db := connectHandler.ConnectDB()
	var professor Professor
	professorID := c.Param("professorID")

	//Check first if there is the id valid for this professor or not
	//If not return specefic message

	ObjectNotFoundError := db.Where("id=?", professorID).Find(&professor)
	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	// Get photo
	image, err := c.FormFile("photo")

	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return err
	}

	//create a folder
	err = os.Mkdir("temp-images", 0755)
	if err != nil {
		log.Println(err)
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("temp-images", "upload-*")
	if err != nil {
		log.Println(err)
	}

	imageName := strings.Split(tempFile.Name(), "/")
	log.Println("name is", imageName[1])

	defer tempFile.Close()

	//Convert (type *multipart.FileHeader) to  type io.Reader
	img, _ := image.Open()

	// read all of the contents of our uploaded image into a
	// byte array

	fileBytes, err := ioutil.ReadAll(img)
	if err != nil {
		log.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	//Insert/Update imageURL for this professor in his/her record in DB
	ObjectNotFoundError = db.Model(&professor).Where("id= ?", professorID).Update("image_url", "static/"+imageName[1])

	// return image path
	// Resquest will be Get method with this url
	path := map[string]string{
		"path":    tempFile.Name(),
		"request": "static/" + imageName[1],
	}

	//insert path in db for this doctor's row image_url
	return c.JSON(http.StatusCreated, path)

}

//GetAllCourses get all courses for specefic professor
func GetAllCourses(c echo.Context) error {
	db := connectHandler.ConnectDB()

	id := c.Param("id")

	var courses []courseObject.Course

	// Execute Query and Get all course object in courses
	ObjectFoundError := db.Where("professor_id=?", id).Find(&courses)

	// Return Empty array
	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, courses)

}

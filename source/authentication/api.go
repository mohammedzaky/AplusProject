package authentication

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	adminObject "gitlab.com/mohamedzaky/aplusProject/source/admin"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	Hash "gitlab.com/mohamedzaky/aplusProject/source/hashHelper"
	professorObject "gitlab.com/mohamedzaky/aplusProject/source/professor"
	studentObject "gitlab.com/mohamedzaky/aplusProject/source/student"
	tokenObject "gitlab.com/mohamedzaky/aplusProject/source/token"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

//login function
func login(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(userObject.User)
	loginAuth := new(LoginAuthentication)
	admin := new(adminObject.Admin)
	professor := new(professorObject.Professor)
	student := new(studentObject.Student)
	c.Bind(loginAuth)

	vaildate = validator.New()

	err := vaildate.Struct(loginAuth)

	// return validation error from Front End if there is an error in bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//check in database if username and password are valid  or not
	// check username and password against DB after hashing the password
	username := loginAuth.UserName
	password := loginAuth.Password

	passHash, _ := Hash.DoHash(password)
	ObjectFoundError := db.Where("user_name = ?", username).Find(&user, "password = ?", passHash)

	if ObjectFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "Your username or password were wrong",
		}
		return c.JSON(http.StatusUnauthorized, result)
	}

	//create type claim to define the user if he is admin , professor or student
	//get first user id from user table
	//then compare it with user_id from professor,admin,student tables
	//userID is the id from student,professor or admin tables
	var userType string
	var userID int

	ObjectFoundError = db.Where("user_id=?", user.ID).Find(&admin)

	if ObjectFoundError.RowsAffected != 0 {
		userType = "Admin"
		userID = admin.ID
	}

	ObjectFoundError = db.Where("user_id=?", user.ID).Find(&professor)

	if ObjectFoundError.RowsAffected != 0 {
		userType = "Professor"
		userID = professor.ID
	}

	ObjectFoundError = db.Where("user_id=?", user.ID).Find(&student)

	if ObjectFoundError.RowsAffected != 0 {
		userType = "Student"
		userID = student.ID
	}

	// create jwt token
	token, err := tokenObject.CreateToken(userID, userType)
	if err != nil {
		result := map[string]string{
			"message": "something went wrong",
		}
		return c.JSON(http.StatusInternalServerError, result)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "You were logged in!",
		"token":   token,
	})
}

//logout function
func logout(c echo.Context) error {
	request := c.Request()
	authToken := request.Header.Get("Authorization")
	result := strings.SplitAfter(authToken, " ")

	//get token value only without Bearer word
	tokenValue := result[1]
	errDelete := tokenObject.DeleteToken(tokenValue)

	message := make(map[string]string)

	if errDelete != nil {

		message = map[string]string{
			"message": errDelete.Error(),
		}
		return c.JSON(http.StatusBadRequest, message)

	}
	message = map[string]string{
		"message": "logout successfuly",
	}
	return c.JSON(http.StatusOK, message)
}

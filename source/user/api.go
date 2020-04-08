package user

import (
	"encoding/json"
	"net/http"

	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	Hash "gitlab.com/mohamedzaky/aplusProject/source/hashHelper"

	"github.com/go-playground/validator"

	"github.com/labstack/echo"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

//NewUsers create new user in DB
func NewUsers(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(User)

	c.Bind(&user)

	vaildate = validator.New()

	err := vaildate.Struct(user)

	// return validation error from Front End if there is an error in bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//psHash to hasing the password and insert it into DB
	password := user.Password
	passHash, _ := Hash.DoHash(password)
	user.Password = passHash

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

	// create a user response for new user API
	var userResponse UserResponse
	userResponse.ID = user.ID
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.UserName = user.UserName
	userResponse.Phone = user.Phone

	return c.JSON(http.StatusCreated, userResponse)
}

//ShowAllUsers show all users from DB
func ShowAllUsers(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var users []User

	db.Find(&users)

	userResponse := make([]UserResponse, len(users))

	for i := 0; i < len(users); i++ {
		userResponse[i].ID = users[i].ID
		userResponse[i].FirstName = users[i].FirstName
		userResponse[i].LastName = users[i].LastName
		userResponse[i].UserName = users[i].UserName
		userResponse[i].Phone = users[i].Phone

	}

	return c.JSON(http.StatusOK, userResponse)
}

//ShowUser show specefic users from DB
func ShowUser(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user User
	var userResponse UserResponse

	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//create a response for userAPI
	userResponse.ID = user.ID
	userResponse.FirstName = user.FirstName
	userResponse.LastName = user.LastName
	userResponse.UserName = user.UserName
	userResponse.Phone = user.Phone

	return c.JSON(http.StatusOK, userResponse)
}

//DeleteUsers delete specefic users from DB
func DeleteUsers(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user User
	id := c.Param("id")
	ObjectNotFoundError := db.Where("id=?", id).Find(&user).Delete(&user)
	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusNoContent, user)
}

//UpdateUsers update specefic users from DB
func UpdateUsers(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(User)
	userUpdate := new(UserUpdate)

	c.Bind(userUpdate)

	vaildate = validator.New()

	userID := c.Param("id")

	err := vaildate.Struct(userUpdate)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	ObjectNotFoundError := db.Where("id=?", userID).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	attrMap := map[string]interface{}{
		"first_name": userUpdate.FirstName,
		"last_name":  userUpdate.LastName,
		"user_name":  userUpdate.UserName,
		"phone":      userUpdate.Phone}

	ObjectNotFoundError = db.Model(&user).Where("id= ?", userID).Updates(attrMap)

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

	return c.JSON(http.StatusOK, userUpdate)
}

// ChangePassword update student password from db
func ChangePassword(c echo.Context) error {
	db := connectHandler.ConnectDB()
	updatePassword := new(UpdatePassword)
	user := new(User)

	c.Bind(updatePassword)

	userID := c.Param("id")

	vaildate = validator.New()

	err := vaildate.Struct(updatePassword)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//check if there is a valid user in db for this id

	ObjectNotFoundError := db.Where("id=?", userID).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	//check if New and Confirm and not the same

	if updatePassword.NewPassword != updatePassword.ConfirmPassword {
		result := map[string]string{
			"message": "Error New Password and Confirm are mismatch !",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	if updatePassword.OldPassword == updatePassword.NewPassword {
		result := map[string]string{
			"message": "Error Old Password and New Password are The Same !",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//Hashing the New password
	newPassword := updatePassword.NewPassword
	passNewHash, _ := Hash.DoHash(newPassword)

	//Hasing the old password from user
	oldPassword := updatePassword.OldPassword
	passOldHash, _ := Hash.DoHash(oldPassword)

	//check if the Entered old password is true or not and compare
	//it from Return Hased password in db

	//return the hash password from db
	returnHash := user.Password

	if returnHash != passOldHash {
		result := map[string]string{
			"message": "old password is not correct",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	//Accecpt New Hash password and Update it

	attrUser := map[string]interface{}{
		"password": passNewHash,
	}

	ObjectUpdateError := db.Model(&user).Where("id= ?", userID).Updates(attrUser)

	if ObjectUpdateError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectUpdateError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)
		//get detail and message from the error back
		result := map[string]string{
			"message": dbErr.Detail,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	result := map[string]string{
		"message": "Password Updated !",
	}
	return c.JSON(http.StatusOK, result)

}

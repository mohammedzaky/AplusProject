package admin

import (
	"encoding/json"
	"net/http"

	Hash "gitlab.com/mohamedzaky/aplusProject/source/hashHelper"

	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
	userObject "gitlab.com/mohamedzaky/aplusProject/source/user"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
)

// create global object of vaildate to check for user input fields
// use a single instance of Validate, it caches struct info
var vaildate *validator.Validate

// NewAdmins insert new admin in db
func NewAdmins(c echo.Context) error {
	db := connectHandler.ConnectDB()
	user := new(userObject.User)
	admin := new(Admin)

	// create professorRequest for request API
	adminRequest := new(AdminRequest)

	// create professorResponse for response in return API
	adminResponse := new(AdminResponse)

	c.Bind(adminRequest)

	vaildate = validator.New()

	err := vaildate.Struct(adminRequest)

	// return validation error from Front End if there is an error in bind
	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//psHash to hasing the password and insert it into DB

	password := adminRequest.Password
	passHash, _ := Hash.DoHash(password)

	//Insert the value in user object
	user.FirstName = adminRequest.FirstName
	user.LastName = adminRequest.LastName
	user.UserName = adminRequest.UserName

	//Put the hased password in user object
	user.Password = passHash

	user.Phone = adminRequest.Phone

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
	admin.Position = adminRequest.Position
	admin.UserID = user.ID

	//create a admin object in db
	errCreate = db.Create(&admin)

	//Create a admin his/her user_id the id of specefic created user
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

	// insert a admin response for new user API
	adminResponse.AdminID = admin.ID
	adminResponse.Position = admin.Position
	adminResponse.User.ID = user.ID
	adminResponse.User.FirstName = user.FirstName
	adminResponse.User.LastName = user.LastName
	adminResponse.User.Username = user.UserName
	adminResponse.User.Phone = user.Phone

	return c.JSON(http.StatusCreated, adminResponse)
}

//DeleteAdmins delete admin row and his/her info
func DeleteAdmins(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var user userObject.User
	var admin Admin
	id := c.Param("id")
	//Check first if there is the id valid for this admin or not
	//If true get the user_id of this admin
	//If not return specefic message

	ObjectNotFoundError := db.Where("id=?", id).Find(&admin)
	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}
	user.ID = admin.UserID

	//Remove admin row first
	ObjectNotFoundError = db.Where("id=?", id).Find(&admin).Delete(&admin)
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

// ShowAdmin get specefic admin from db
func ShowAdmin(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var admin Admin
	var user userObject.User

	id := c.Param("id")

	//check if there id is valid or not
	ObjectNotFoundError := db.Where("id=?", id).Find(&admin)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}
	usrID := admin.UserID

	ObjectNotFoundError = db.Where("id=?", usrID).Find(&user)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	//create a response for studentAPI

	var adminResponse AdminResponse
	adminResponse.AdminID = admin.ID
	adminResponse.Position = admin.Position
	adminResponse.User.ID = user.ID
	adminResponse.User.FirstName = user.FirstName
	adminResponse.User.LastName = user.LastName
	adminResponse.User.Username = user.UserName
	adminResponse.User.Phone = user.Phone

	return c.JSON(http.StatusOK, adminResponse)
}

//ShowAllAdmins Struct
func ShowAllAdmins(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var admin []Admin

	db.Find(&admin)

	var userIds []int
	//get all userid for each admin in slice userIds
	for i := 0; i < len(admin); i++ {
		userIds = append(userIds, admin[i].UserID)
	}

	user := make([]userObject.User, len(userIds))

	//get all user info from user table based on user id from admin table
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
	adminResponse := make([]AdminResponse, len(userIds))

	for i := 0; i < len(userIds); i++ {
		adminResponse[i].AdminID = admin[i].ID
		adminResponse[i].Position = admin[i].Position
		adminResponse[i].User.FirstName = user[i].FirstName
		adminResponse[i].User.LastName = user[i].LastName
		adminResponse[i].User.Username = user[i].UserName
		adminResponse[i].User.Phone = user[i].Phone

	}
	return c.JSON(http.StatusOK, adminResponse)
}

//UpdateAdmins Struct
func UpdateAdmins(c echo.Context) error {
	db := connectHandler.ConnectDB()
	adminUpdate := new(AdminUpdate)
	admin := new(Admin)
	user := new(userObject.User)

	c.Bind(adminUpdate)

	vaildate = validator.New()

	id := c.Param("id")

	err := vaildate.Struct(adminUpdate)

	if err != nil {
		result := map[string]string{
			"message": err.Error(),
		}
		return c.JSON(http.StatusBadRequest, result)
	}

	//Check first if there is the id valid for this Admin or not
	//If true get the user_id of this Admin
	//If not return specefic message

	checkAdmin := new(Admin)
	ObjectNotFoundError := db.Where("id=?", id).Find(checkAdmin)

	if ObjectNotFoundError.RowsAffected == 0 {

		//get detail and message from the error back
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	//Get the admin info from adminUpdate object and insert it into admin table
	admin.Position = adminUpdate.Position

	// update the professor info from professor table
	attrAdmin := map[string]interface{}{
		"position": admin.Position,
	}

	ObjectNotFoundError = db.Model(&admin).Where("id= ?", id).Updates(attrAdmin)

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
	//First Get the user id from admin table to update user row in user table for this specefic id
	db.Where("user_id=?", admin.UserID)

	//Get in user info from adminUpdate object and insert it into user table
	user.FirstName = adminUpdate.FirstName
	user.LastName = adminUpdate.LastName
	user.UserName = adminUpdate.UserName
	user.Phone = adminUpdate.Phone

	userID := checkAdmin.UserID

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
			"message": dbErr.Message,
		}
		return c.JSON(http.StatusBadRequest, result)
	}
	var adminResponse AdminResponse
	adminResponse.AdminID = checkAdmin.ID
	adminResponse.Position = admin.Position
	adminResponse.User.ID = userID
	adminResponse.User.FirstName = user.FirstName
	adminResponse.User.LastName = user.LastName
	adminResponse.User.Username = user.UserName
	adminResponse.User.Phone = user.Phone

	return c.JSON(http.StatusOK, adminResponse)

}

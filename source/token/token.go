package token

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	connectHandler "gitlab.com/mohamedzaky/aplusProject/source/connectDB"
)

//Token value struct
type Token struct {
	ID         int    `json:"id"`
	TokenValue string `json:"token_value" gorm:"unique_index:idx_token_value"`
	ExpireAt   int    `json:"expire_at"`
	UserID     int    `json:"user_id"`
}

//JwtClaims name of user, jwt expire At
type JwtClaims struct {
	Name     string `json:"type"`
	UserType int    `json:"id"`
	jwt.StandardClaims
}

//NewTokenDB function generate a new token and insert it in db for specefic user for that token
//in token table
func NewTokenDB(tokenValue string, expiresAt int, userID int) string {
	db := connectHandler.ConnectDB()
	token := new(Token)

	token.TokenValue = tokenValue
	token.UserID = userID
	token.ExpireAt = expiresAt

	errCreate := db.Create(&token)

	if errCreate.RowsAffected == 0 {
		//make serialization to get specefic entity for error back from db from struct ErrorModel
		//get the error in json format from db by marshaling it and unmarshaling

		errMessage, _ := json.Marshal(errCreate.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//return detail or message from the error back
		return dbErr.Message
	}
	return ""
}

//DeleteAllTokens function
func DeleteAllTokens(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var token Token
	id := c.Param("id")

	ObjectNotFoundError := db.Where("id=?", id).Find(&token).Delete(&token)
	if ObjectNotFoundError.RowsAffected == 0 {

		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}

		return c.JSON(http.StatusBadRequest, result)
	}

	return c.JSON(http.StatusNoContent, token)

}

//ShowToken function
func ShowToken(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var token Token

	id := c.Param("id")

	//check if there id is valid or not
	ObjectNotFoundError := db.Where("id=?", id).Find(&token)

	if ObjectNotFoundError.RowsAffected == 0 {
		result := map[string]string{
			"message": "No Record Found in DB for this id",
		}
		return c.JSON(http.StatusNotFound, result)
	}

	return c.JSON(http.StatusOK, token)
}

// ShowAllTokens function
func ShowAllTokens(c echo.Context) error {
	db := connectHandler.ConnectDB()
	var token []Token
	db.Find(&token)
	return c.JSON(http.StatusOK, token)
}

//GetToken function get row info for that user
//In token table
func GetToken(tokenValue string) (*Token, string) {
	db := connectHandler.ConnectDB()
	token := new(Token)

	var result string
	//check if there id is valid or not
	ObjectNotFoundError := db.Where("token_value = ?", tokenValue).Find(&token)

	if ObjectNotFoundError.RowsAffected == 0 {
		token = nil
		result = "message: No Record Found in DB for this id"
		return token, result
	}
	result = ""
	return token, result
}

//UpdateToken function
func UpdateToken(tokenRaw string, expiresAt int, userID int) (*Token, string) {

	db := connectHandler.ConnectDB()

	token := new(Token)

	checkToken := new(Token)

	var result string
	//check if there id is valid or not
	ObjectNotFoundError := db.Where("token_value=?", tokenRaw).Find(&checkToken)

	if ObjectNotFoundError.RowsAffected == 0 {
		token = nil
		result = "message: No Record Found in DB for this id"
		return token, result
	}

	token.TokenValue = tokenRaw
	token.ExpireAt = expiresAt
	token.UserID = userID

	attrMap := map[string]interface{}{
		"token_value": token.TokenValue,
		"expire_at":   token.ExpireAt,
		"user_id":     token.UserID,
	}

	ObjectFoundError := db.Model(&token).Where("token_value= ?", tokenRaw).Updates(attrMap)

	if ObjectFoundError.RowsAffected == 0 {
		errMessage, _ := json.Marshal(ObjectFoundError.Error)
		var dbErr connectHandler.ErrorModel
		json.Unmarshal(errMessage, &dbErr)

		//get detail and message from the error back
		result = "message :" + dbErr.Message

		return token, result
	}
	result = ""
	return token, result
}

//CreateToken function for the first time after login
//or increase expire time of token and update it in token table
func CreateToken(userID int, userType string) (string, error) {

	claims := JwtClaims{
		userType,
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
	}

	expire := claims.ExpiresAt

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	_, key := config.SigningConfig()

	token, err := rawToken.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	//insert the token value in DB
	errInsertToken := NewTokenDB(token, (int)(expire), userID)

	//convert string to error
	//string that return from newToken function

	errInsert := errors.New(errInsertToken)

	if errInsertToken != "" {
		return "", errInsert
	}

	return token, nil
}

//DeleteToken function
func DeleteToken(tokenValue string) error {
	db := connectHandler.ConnectDB()
	token := new(Token)

	var message string

	ObjectNotFoundError := db.Where("token_value=?", tokenValue).Find(&token).Delete(&token)

	if ObjectNotFoundError.RowsAffected == 0 {
		token = nil
		message = "No Record Found in DB for this id"

		//convert string to error
		errDelete := errors.New(message)

		return errDelete
	}

	return nil
}

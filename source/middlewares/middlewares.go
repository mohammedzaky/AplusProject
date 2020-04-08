package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	tokenObject "gitlab.com/mohamedzaky/aplusProject/source/token"
)

// CheckTokenMiddleware middleware check token if it is valid or expired
// if token is valid return 200 OK
// else will return 401 unautorized and message by Default "invalid or expired jwt"
// Or try to access token that reomved from token table
func CheckTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		request := c.Request()
		authToken := request.Header.Get("Authorization")
		result := strings.SplitAfter(authToken, " ")

		// _, key := config.SigningConfig()

		//get token value only without Bearer word
		tokenValue := result[1]

		//Get the row of this token value from token table in database
		_, errMessage := tokenObject.GetToken(tokenValue)

		//if user try to access token but it is deleted from database
		//handling this case
		if errMessage != "" {

			message := map[string]string{
				"message": "invalid or expired jwt",
			}
			return c.JSON(http.StatusUnauthorized, message)
		}

		return next(c)
	}
}

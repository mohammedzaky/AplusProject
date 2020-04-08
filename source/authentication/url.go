package authentication

import (
	"github.com/labstack/echo"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	// API End Point
	e.POST("/login", login)
	e.POST("/logout", logout)

}

package enroll

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//Create Courses Group
	enrollGroup := e.Group("/enrolls")

	method, key := config.SigningConfig()

	//Set all middleware
	enrollGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	enrollGroup.Use(middlewares.CheckTokenMiddleware)

	//Enroll EndPoints
	enrollGroup.GET("", ShowAllEnrolls)
	enrollGroup.POST("", NewEnrolls)

}

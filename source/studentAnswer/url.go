package studentAnswer

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Choice Group
	studentAnswerGroup := e.Group("/student/answers")

	method, key := config.SigningConfig()

	//check middleware for valid token or not
	studentAnswerGroup.Use(middlewares.CheckTokenMiddleware)

	//Set all middleware
	studentAnswerGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	// studentAnswers End Point
	studentAnswerGroup.GET("", ShowAllStudentAnswers)
	studentAnswerGroup.POST("", NewStudentAnswer)
	studentAnswerGroup.GET("/:studentID/exam/:examID", GetAllStudentAnswer)
}

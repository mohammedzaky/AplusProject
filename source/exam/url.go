package exam

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Exam Group
	examGroup := e.Group("/exams")

	method, key := config.SigningConfig()

	//Set all middleware
	examGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	examGroup.Use(middlewares.CheckTokenMiddleware)

	// Exams
	examGroup.GET("", ShowAllExams)
	examGroup.GET("/:id", ShowExam)
	examGroup.PUT("/:id", UpdateExam)
	examGroup.POST("", NewExams)
	examGroup.DELETE("/:id", DeleteExam)
	//for professor
	examGroup.GET("/:id/questions/answers", GetExamModel)
	examGroup.POST("/:id/questions/answers", CreateExamModel)
	//for student
	examGroup.GET("/:id/questions", GetAllQuestionsChoices)

}

package studentExam

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	"gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//Create Courses Group
	studentExamGroup := e.Group("/student")

	method, key := config.SigningConfig()

	//Set all middleware
	studentExamGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	studentExamGroup.Use(middlewares.CheckTokenMiddleware)

	//StudentExam EndPoint
	studentExamGroup.GET("/exams", ShowAllStudentExams)

	studentExamGroup.POST("/exam", NewStudentExam)

	//update degree for specefic student
	studentExamGroup.PUT("/exam", UpdateStudentExam)

	//calculate and show the student degree after examed

	studentExamGroup.POST("/:studentID/exam/:examID/degree", NewStudentDegree)

	//the student degree after he/she examed
	studentExamGroup.GET("/:studentID/exam/:examID/degree", GetStudentDegree)

	//create Exam Group
	examGroup := e.Group("/exam")

	//Set all middleware
	examGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	examGroup.Use(middlewares.CheckTokenMiddleware)

	//get all students degrees after examed
	examGroup.GET("/:id/degrees", GetStudentDegrees)

	//Set all student degrees of null
	examGroup.PUT("/:id/degrees", ResetStudentDegrees)

	courseGroup := e.Group("/course")

	//Set all middleware
	courseGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	courseGroup.Use(middlewares.CheckTokenMiddleware)

	//get all students degrees after examed for specefic course
	courseGroup.GET("/:id/degrees", GetAllStudentDegrees)

	//Set all student degrees of null for all exams for specefic course
	courseGroup.PUT("/:id/degrees", ResetAllStudentDegrees)
}

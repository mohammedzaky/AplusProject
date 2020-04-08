package course

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	"gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Courses Group
	courseGroup := e.Group("/courses")
	method, key := config.SigningConfig()

	//Set all middleware
	courseGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	courseGroup.Use(middlewares.CheckTokenMiddleware)

	//Courses
	courseGroup.GET("", ShowAllCourses)
	courseGroup.GET("/:id", ShowCourse)
	courseGroup.PUT("/:id", UpdateCourse)
	courseGroup.POST("", NewCourses)
	courseGroup.DELETE("/:id", DeleteCourse)
	courseGroup.GET("/:id/exams", GetAllExams)
	courseGroup.GET("/:id/students", ShowAllStudentIds)
	courseGroup.GET("/:id/exam", GetExamCourse)

}

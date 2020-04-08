package semester

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Semester Group
	semesterGroup := e.Group("/semesters")

	method, key := config.SigningConfig()

	//Set all middleware
	semesterGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	semesterGroup.Use(middlewares.CheckTokenMiddleware)

	// All Semesters EndPoints
	semesterGroup.GET("", ShowAllSemesters)
	semesterGroup.GET("/:id", ShowSemester)
	semesterGroup.PUT("/:id", UpdateSemester)
	semesterGroup.POST("", NewSemesters)
	semesterGroup.DELETE("/:id", DeleteSemester)
	semesterGroup.GET("/:id/courses", GetAllCourses)
	semesterGroup.POST("/:id/courses", newCourseInSemester)
}

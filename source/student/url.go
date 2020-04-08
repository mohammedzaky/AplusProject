package student

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	"gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Student Group
	studentGroup := e.Group("/students")

	method, key := config.SigningConfig()

	//Set all middleware
	studentGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	studentGroup.Use(middlewares.CheckTokenMiddleware)

	//Student
	studentGroup.POST("", NewStudent)
	studentGroup.GET("", ShowAllStudents)
	studentGroup.GET("/:id", ShowStudent)
	studentGroup.PUT("/:id", UpdateStudent)
	studentGroup.DELETE("/:id", DeleteStudent)
	studentGroup.GET("/:id/courses", ShowAllCourseIds)

}

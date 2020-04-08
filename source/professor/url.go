package professor

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Professor Group
	professorGroup := e.Group("/professors")

	method, key := config.SigningConfig()

	//Set all middleware
	professorGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	professorGroup.Use(middlewares.CheckTokenMiddleware)

	// Professors
	professorGroup.GET("", ShowAllProfessors)
	professorGroup.GET("/:id", ShowProfessor)
	professorGroup.PUT("/:id", UpdateProfessors)
	professorGroup.POST("", NewProfessor)
	professorGroup.DELETE("/:id", DeleteProfessors)
	professorGroup.POST("/upload-image/:professorID", UploadImage)
	professorGroup.GET("/:id/courses", GetAllCourses)

}

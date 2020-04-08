package question

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Question Group
	questionGroup := e.Group("/questions")

	method, key := config.SigningConfig()

	//Set all middleware
	questionGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	questionGroup.Use(middlewares.CheckTokenMiddleware)

	//Questions
	questionGroup.GET("", ShowAllQuestions)
	questionGroup.GET("/:id", ShowQuestion)
	questionGroup.PUT("/:id", UpdateQuestion)
	questionGroup.POST("", NewQuestions)
	questionGroup.DELETE("/:id", DeleteQuestion)
}

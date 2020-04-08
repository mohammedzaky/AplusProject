package choice

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Choice Group
	choiceGroup := e.Group("/choices")

	method, key := config.SigningConfig()

	//Set all middleware
	choiceGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	choiceGroup.Use(middlewares.CheckTokenMiddleware)

	// Choices
	choiceGroup.GET("", ShowAllChoices)
	choiceGroup.GET("/:id", ShowChoice)
	choiceGroup.PUT("/:id", UpdateChoices)
	choiceGroup.POST("", NewChoices)
	choiceGroup.DELETE("/:id", DeleteChoices)

}

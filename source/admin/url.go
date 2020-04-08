package admin

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	middlewares "gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Admin Group
	adminGroup := e.Group("/admins")

	method, key := config.SigningConfig()

	//Set all middlewares
	adminGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	adminGroup.Use(middlewares.CheckTokenMiddleware)

	//Admins
	adminGroup.GET("", ShowAllAdmins)
	adminGroup.GET("/:id", ShowAdmin)
	adminGroup.PUT("/:id", UpdateAdmins)
	adminGroup.POST("", NewAdmins)
	adminGroup.DELETE("/:id", DeleteAdmins)

}

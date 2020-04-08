package user

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	config "gitlab.com/mohamedzaky/aplusProject/source/configration"
	"gitlab.com/mohamedzaky/aplusProject/source/middlewares"
)

//MainGroup to call all api endpoints
func MainGroup(e *echo.Echo) {

	//create Choice Group
	userGroup := e.Group("/users")

	method, key := config.SigningConfig()

	//Set all middleware
	userGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: method,
		SigningKey:    []byte(key),
	}))

	//check middleware for valid token or not
	userGroup.Use(middlewares.CheckTokenMiddleware)

	// Users
	userGroup.GET("", ShowAllUsers)
	userGroup.GET("/:id", ShowUser)
	userGroup.PUT("/:id", UpdateUsers)
	userGroup.POST("", NewUsers)
	userGroup.DELETE("/:id", DeleteUsers)
	userGroup.PATCH("/change-password/:id", ChangePassword)
}

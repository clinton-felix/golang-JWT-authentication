package routes

import (
	controller "github.com/clinton-felix/golang-JWT-auth-project/controllers"
	middleware "github.com/clinton-felix/golang-JWT-auth-project/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoute *gin.Engine)  {
	// use middleware to ensure that users are authenticated with Authenticate Fn
	incomingRoute.Use(middleware.Authenticate()) 
	incomingRoute.GET("/users", controller.GetUsers())
	incomingRoute.GET("/users/:user_id", controller.GetUser())
}
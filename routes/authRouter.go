package routes

import (
	controller "github.com/clinton-felix/golang-JWT-auth-project/controllers"
	"github.com/gin-gonic/gin"
)

// creating the AuthRoutes function which accepts incoming routes
func AuthRoutes(incomingRoutes *gin.Engine)  {
	// handle the signup and login routes with controllers
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}
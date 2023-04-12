package routes

import (
	"oprec/controller"
	"oprec/service"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, userController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)
	}

}

package routes

import (
	"mods/controller"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, userController controller.UserController, storeController controller.StoreController, jwtService service.JWTService) {
	userRoutes := router.Group("/user")
	{
		userRoutes.POST("", userController.RegisterUser)
		userRoutes.POST("/login", userController.LoginUser)

		userRoutes.GET("/me", userController.ProfilePage)
		userRoutes.POST("/upload", userController.UploadGame)
		userRoutes.POST("/purchase/:id", userController.PurchaseGame)

	}

	storeMainPage := router.Group("/storeMainPage")
	{
		storeMainPage.GET("/featured", storeController.Featured)
		storeMainPage.GET("/categories", storeController.Categories)

		storeMainPage.GET("/game/:id", storeController.GamePage)
	}

}

package routes

import (
	"mods/controller"
	"mods/middleware"
	"mods/service"

	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine, userController controller.UserController, storeController controller.StoreController, jwtService service.JWTService) {
	userPublic := router.Group("/user")
	{
		// public can access
		userPublic.POST("", userController.RegisterUser)
		userPublic.POST("/login", userController.LoginUser)
	}

	userRoutes := router.Group("/secured/user").Use(middleware.Authenticate())
	{
		// profiles
		userRoutes.GET("/me", userController.ProfilePage)
		userRoutes.GET("/medev", userController.DeveloperProfile)

		// transactional
		userRoutes.POST("/purchase/:id", userController.PurchaseGame)
		userRoutes.POST("/topup", userController.TopUp)
		userRoutes.POST("/purchasedlc/:id", userController.PurchaseDLC)

		// upload new game or dlc
		userRoutes.POST("/upload", userController.UploadGame)
		userRoutes.POST("/updlc", userController.UploadDLC)

		// add tags, language, os
		userRoutes.POST("/add/:method", userController.AddToGame)

	}

	storeMainPage := router.Group("/storeMainPage")
	{
		// public can access
		storeMainPage.GET("/featured", storeController.Featured)
		storeMainPage.GET("/categories", storeController.Categories)
		storeMainPage.GET("/game/:id", storeController.GamePage)
		storeMainPage.GET("/game/all", storeController.GetAllGames)
		storeMainPage.GET("/dlc/:id", storeController.DLCGame)
	}

}

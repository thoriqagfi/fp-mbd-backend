package main

import (
	"log"
	"oprec/FP-MBD/config"
	"oprec/FP-MBD/controller"
	"oprec/FP-MBD/middleware"
	"oprec/FP-MBD/repository"
	"oprec/FP-MBD/routes"
	"oprec/FP-MBD/service"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	db := config.SetupDatabaseConnection()

	jwtService := service.NewJWTService()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService, jwtService)

	defer config.CloseDatabaseConnection(db)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	routes.UserRouter(server, userController, jwtService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server.Run(":" + port)
}

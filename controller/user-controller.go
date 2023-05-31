package controller

import (
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

type UserController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	UploadGame(ctx *gin.Context)
}

func NewUserController(us service.UserService, jwt service.JWTService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwt,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateDTO
	if tx := ctx.ShouldBind(&user); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	isEmailRegistered, _ := uc.userService.IsDuplicateEmail(ctx.Request.Context(), user.Email)
	if isEmailRegistered {
		res := utils.BuildErrorResponse("Email already registered", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := uc.userService.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildErrorResponse("Failed to register user", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse("Success to register user", http.StatusOK, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var userLogin dto.UserLoginDTO
	if tx := ctx.ShouldBind(&userLogin); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	user, err := uc.userService.GetUserByEmail(ctx.Request.Context(), userLogin.Email)
	if err != nil {
		res := utils.BuildErrorResponse("Failed to login, email not registered", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	isValid, _ := uc.userService.VerifyCredential(ctx.Request.Context(), userLogin.Email, userLogin.Password)
	if !isValid {
		res := utils.BuildErrorResponse("Failed to login, email and password do not match", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := uc.jwtService.GenerateToken(user.ID, user.Role)
	res := utils.BuildResponse("Successful login", http.StatusOK, token)
	ctx.JSON(http.StatusOK, res)

}

func (uc *userController) UploadGame(ctx *gin.Context) {
	var gameDTO dto.UploadGame
	if tx := ctx.ShouldBind(&gameDTO); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)
	tokenService := service.NewJWTService()

	idUser, err := tokenService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.UploadGame(ctx, gameDTO, idUser)
	if err != nil {
		res := utils.BuildErrorResponse("Failed to upload game", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("upload game berhasil", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)
}

package controller

import (
	"mods/dto"
	"mods/service"
	"mods/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService service.UserService
	jwtService  service.JWTService
}

type UserController interface {
	// regist login
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)

	// profiles
	ProfilePage(ctx *gin.Context)
	DeveloperProfile(ctx *gin.Context)

	// transactional
	UploadGame(ctx *gin.Context)
	PurchaseGame(ctx *gin.Context)
	TopUp(ctx *gin.Context)
	UploadDLC(ctx *gin.Context)
	PurchaseDLC(ctx *gin.Context)

	// add languages tags OS
	AddToGame(ctx *gin.Context)
}

func NewUserController(us service.UserService, jwt service.JWTService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwt,
	}
}

func (uc *userController) RetrieveID(ctx *gin.Context) (uint64, error) {
	token := ctx.GetHeader("Authorization")
	token = strings.Replace(token, "Bearer ", "", -1)

	return uc.jwtService.GetUserIDByToken(token)
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

	idUser, err := uc.RetrieveID(ctx)
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

func (uc *userController) PurchaseGame(ctx *gin.Context) {
	var transaksi dto.PurchaseGame
	if tx := ctx.ShouldBind(&transaksi); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	gameid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	idUser, err := uc.RetrieveID(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.PurchaseGame(ctx, gameid, idUser, transaksi.MetodeBayar)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("pembelian game berhasil", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)
}

func (uc *userController) ProfilePage(ctx *gin.Context) {
	idUser, err := uc.RetrieveID(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.UserProfile(ctx, idUser)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("success to get profile", http.StatusOK, res)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) TopUp(ctx *gin.Context) {
	var topup dto.PurchaseGame
	if tx := ctx.ShouldBind(&topup); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	nominal, _ := strconv.ParseUint(topup.Nominal, 10, 64)

	idUser, err := uc.RetrieveID(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.TopUp(ctx, idUser, nominal)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("success to add steam wallet funds", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)
}

func (uc *userController) DeveloperProfile(ctx *gin.Context) {
	idDev, err := uc.RetrieveID(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.DeveloperProfile(ctx, idDev)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("success to get developer profile", http.StatusOK, res)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) UploadDLC(ctx *gin.Context) {
	var dlc dto.UploadDLC
	if tx := ctx.ShouldBind(&dlc); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := uc.userService.UploadDLC(ctx, dlc)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("success to upload dlc", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)

}

func (uc *userController) PurchaseDLC(ctx *gin.Context) {
	var transaksi dto.PurchaseGame
	if tx := ctx.ShouldBind(&transaksi); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	dlcid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	idUser, err := uc.RetrieveID(ctx)
	if err != nil {
		response := utils.BuildErrorResponse("gagal memproses request", http.StatusBadRequest)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.PurchaseDLC(ctx, dlcid, idUser, transaksi.MetodeBayar)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("pembelian dlc berhasil", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)
}

func (uc *userController) AddToGame(ctx *gin.Context) {
	chosen := ctx.Param("method")

	var add dto.Add
	if tx := ctx.ShouldBind(&add); tx != nil {
		res := utils.BuildErrorResponse("Failed to process request", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := uc.userService.AddToGame(add.Nama, add.GameID, chosen)
	if err != nil {
		res := utils.BuildErrorResponse(err.Error(), http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	response := utils.BuildResponse("penambahan "+chosen+" berhasil", http.StatusOK, res)
	ctx.JSON(http.StatusCreated, response)
}

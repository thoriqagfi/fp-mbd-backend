package controller

import (
	"mods/service"
	"mods/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type storeController struct {
	storeService service.StoreService
}

type StoreController interface {
	Featured(ctx *gin.Context)
	Categories(ctx *gin.Context)
}

func NewStoreController(ss service.StoreService) StoreController {
	return &storeController{
		storeService: ss,
	}
}

func (sc *storeController) Featured(ctx *gin.Context) {
	featuredInfo, err := sc.storeService.GetFeatured(ctx)
	if err != nil {
		res := utils.BuildErrorResponse("failed to get featured information", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse("success to get featured information", http.StatusOK, featuredInfo)
	ctx.JSON(http.StatusOK, res)
}

func (sc *storeController) Categories(ctx *gin.Context) {
	categoriesInfo, err := sc.storeService.GetCategories(ctx)
	if err != nil {
		res := utils.BuildErrorResponse("failed to get categories info", http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildResponse("success to get categories info", http.StatusOK, categoriesInfo)
	ctx.JSON(http.StatusOK, res)
}

package controllers

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/dig"
	"net/http"
	"shop-bot/config"
	"shop-bot/internal/logger"
	"shop-bot/repositories"
	"shop-bot/types/responses"
)

type ICategoryItemController interface {
	Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
}

type categoryItemControllerDependencies struct {
	dig.In

	Logger                 logger.ILogger                       `name:"Logger"`
	Config                 config.IConfig                       `name:"Config"`
	CategoryItemRepository repositories.ICategoryItemRepository `name:"CategoryItemRepository"`
}

type categoryItemController struct {
	logger                 logger.ILogger
	config                 config.IConfig
	categoryItemRepository repositories.ICategoryItemRepository
}

func NewCategoryItemController(deps categoryItemControllerDependencies) *categoryItemController {
	return &categoryItemController{
		logger:                 deps.Logger,
		config:                 deps.Config,
		categoryItemRepository: deps.CategoryItemRepository,
	}
}

func (h *categoryItemController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) *responses.Response {
	return responses.NewSuccessResponse(nil)
}

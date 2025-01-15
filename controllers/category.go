package controllers

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/dig"
	"net/http"
	"shop-bot/config"
	"shop-bot/internal/logger"
	"shop-bot/repositories"
	"shop-bot/types/requests"
	"shop-bot/types/responses"
	"shop-bot/utils"
)

type ICategoryController interface {
	Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
	GetById(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
	GetAll(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
	UpdateById(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
	DeleteById(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response
}

type categoryControllerDependencies struct {
	dig.In

	Logger             logger.ILogger                   `name:"Logger"`
	Config             config.IConfig                   `name:"Config"`
	CategoryRepository repositories.ICategoryRepository `name:"CategoryRepository"`
}

type categoryController struct {
	logger             logger.ILogger
	config             config.IConfig
	categoryRepository repositories.ICategoryRepository
}

func NewCategoryController(deps categoryControllerDependencies) *categoryController {
	return &categoryController{
		logger:             deps.Logger,
		config:             deps.Config,
		categoryRepository: deps.CategoryRepository,
	}
}

func (h *categoryController) Create(w http.ResponseWriter, r *http.Request, _ httprouter.Params) *responses.Response {
	var req requests.CreateCategoryRequest

	if err := utils.DecodeJSONBody(w, r, &req, h.config.MaxJSONBodySizeInBytes()); err != nil {
		return err
	}

	existingCategory := h.categoryRepository.GetByName(r.Context(), req.Name)

	if existingCategory != nil {
		return responses.NewBadRequestResponse(fmt.Sprintf("Category with name %s already exists", req.Name))
	}

	categoryId := h.categoryRepository.Create(r.Context(), req.ToCategory())

	category := h.categoryRepository.GetById(r.Context(), categoryId)

	return responses.NewSuccessResponse(category)
}

func (h *categoryController) GetAll(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) *responses.Response {
	limit, err := utils.GetLimitQueryParam(r, 20, 100)

	if err != nil {
		return err
	}

	offset, err := utils.GetOffsetQueryParam(r, 0)

	if err != nil {
		return err
	}

	categories := h.categoryRepository.GetAll(r.Context(), limit, offset)

	return responses.NewSuccessResponse(categories)
}

func (h *categoryController) GetById(_ http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response {
	id, err := utils.GetIdParam(p)

	if err != nil {
		return err
	}

	category := h.categoryRepository.GetById(r.Context(), id)

	if category == nil {
		return responses.NewNotFoundResponse(fmt.Sprintf("Category with id %d not found", id))
	}

	return responses.NewSuccessResponse(category)
}

func (h *categoryController) UpdateById(w http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response {
	id, err := utils.GetIdParam(p)

	if err != nil {
		return err
	}

	category := h.categoryRepository.GetById(r.Context(), id)

	if category == nil {
		return responses.NewNotFoundResponse(fmt.Sprintf("Category with id %d not found", id))
	}

	var req requests.UpdateCategoryRequest

	if err := utils.DecodeJSONBody(w, r, &req, h.config.MaxJSONBodySizeInBytes()); err != nil {
		return err
	}

	if req.Name != "" {
		existingCategory := h.categoryRepository.GetByName(r.Context(), req.Name)

		if existingCategory != nil {
			return responses.NewBadRequestResponse(fmt.Sprintf("Category with name %s already exists", req.Name))
		}
	}

	h.categoryRepository.UpdateById(r.Context(), id, req.ToCategory())

	category = h.categoryRepository.GetById(r.Context(), id)

	return responses.NewSuccessResponse(category)
}

func (h *categoryController) DeleteById(_ http.ResponseWriter, r *http.Request, p httprouter.Params) *responses.Response {
	id, err := utils.GetIdParam(p)

	if err != nil {
		return err
	}

	category := h.categoryRepository.GetById(r.Context(), id)

	if category == nil {
		return responses.NewNotFoundResponse(fmt.Sprintf("Category with id %d not found", id))
	}

	h.categoryRepository.DeleteById(r.Context(), id)

	return responses.NewSuccessResponse(category)
}

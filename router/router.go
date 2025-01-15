package router

import (
	"github.com/julienschmidt/httprouter"
	"go.uber.org/dig"
	"shop-bot/config"
	"shop-bot/controllers"
	"shop-bot/internal/logger"
)

type IRouter interface {
	GetHttpRouter() *httprouter.Router
}

type routerDependencies struct {
	dig.In

	Logger logger.ILogger `name:"Logger"`
	Config config.IConfig `name:"Config"`

	CategoryController     controllers.ICategoryController     `name:"CategoryController"`
	CategoryItemController controllers.ICategoryItemController `name:"CategoryItemController"`
}

type router struct {
	httpRouter *httprouter.Router

	logger logger.ILogger
	config config.IConfig

	categoryController     controllers.ICategoryController
	categoryItemController controllers.ICategoryItemController
}

func NewRouter(deps routerDependencies) *router {
	r := &router{
		httpRouter:             httprouter.New(),
		logger:                 deps.Logger,
		config:                 deps.Config,
		categoryController:     deps.CategoryController,
		categoryItemController: deps.CategoryItemController,
	}

	r.registerRoutes()

	return r
}

func (router *router) registerRoutes() {
	router.httpRouter.POST("/api/categories", router.withMiddlewares(router.categoryController.Create))
	router.httpRouter.GET("/api/categories/:id", router.withMiddlewares(router.categoryController.GetById))
	router.httpRouter.GET("/api/categories", router.withMiddlewares(router.categoryController.GetAll))
	router.httpRouter.PUT("/api/categories/:id", router.withMiddlewares(router.categoryController.UpdateById))
	router.httpRouter.DELETE("/api/categories/:id", router.withMiddlewares(router.categoryController.DeleteById))

	router.httpRouter.POST("/api/categories/:categoryId/items", router.withMiddlewares(router.categoryItemController.Create))
}

func (router *router) GetHttpRouter() *httprouter.Router {
	return router.httpRouter
}

package di

import (
	"go.uber.org/dig"
	"shop-bot/config"
	"shop-bot/controllers"
	"shop-bot/handlers"
	"shop-bot/handlers/callback_queries"
	"shop-bot/internal/bot"
	"shop-bot/internal/db"
	"shop-bot/internal/http"
	"shop-bot/internal/logger"
	"shop-bot/repositories"
	"shop-bot/router"
	"shop-bot/services"
	"shop-bot/utils"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	c = AppendDependenciesToContainer(c, getRequiredDependencies())
	c = AppendDependenciesToContainer(c, getRepositoriesDependencies())
	c = AppendDependenciesToContainer(c, getServicesDependencies())
	c = AppendDependenciesToContainer(c, getHandlersDependencies())
	c = AppendDependenciesToContainer(c, getControllersDependencies())
	c = AppendDependenciesToContainer(c, getBotDependencies())
	c = AppendDependenciesToContainer(c, getHttpServerDependencies())

	return c
}

func AppendDependenciesToContainer(container *dig.Container, dependencies []Dependency) *dig.Container {
	for _, dep := range dependencies {
		mustProvideDependency(container, dep)
	}

	return container
}

func mustProvideDependency(container *dig.Container, dependency Dependency) {
	if dependency.Interface == nil {
		utils.PanicIfError(container.Provide(dependency.Constructor, dig.Name(dependency.Token)))
		return
	}

	utils.PanicIfError(container.Provide(
		dependency.Constructor,
		dig.As(dependency.Interface),
		dig.Name(dependency.Token),
	))
}

func getRequiredDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: logger.NewLogger,
			Interface:   new(logger.ILogger),
			Token:       "Logger",
		},
		{
			Constructor: config.NewConfig,
			Interface:   new(config.IConfig),
			Token:       "Config",
		},
		{
			Constructor: db.NewDB,
			Interface:   nil,
			Token:       "DB",
		},
	}
}

func getRepositoriesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: repositories.NewUserRepository,
			Interface:   new(repositories.IUserRepository),
			Token:       "UserRepository",
		},
		{
			Constructor: repositories.NewCategoryRepository,
			Interface:   new(repositories.ICategoryRepository),
			Token:       "CategoryRepository",
		},
		{
			Constructor: repositories.NewCategoryItemRepository,
			Interface:   new(repositories.ICategoryItemRepository),
			Token:       "CategoryItemRepository",
		},
	}
}

func getServicesDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: services.NewTextService,
			Interface:   new(services.ITextService),
			Token:       "TextService",
		},
		{
			Constructor: services.NewUserService,
			Interface:   new(services.IUserService),
			Token:       "UserService",
		},
	}
}

func getHandlersDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: handlers.NewDefaultHandler,
			Interface:   new(handlers.IDefaultHandler),
			Token:       "DefaultHandler",
		},
		{
			Constructor: handlers.NewCommandHandler,
			Interface:   new(handlers.ICommandHandler),
			Token:       "CommandHandler",
		},
		{
			Constructor: callback_queries.NewMenuHandler,
			Interface:   new(callback_queries.IMenuHandler),
			Token:       "MenuHandler",
		},
		{
			Constructor: callback_queries.NewUserAgreementHandler,
			Interface:   new(callback_queries.IUserAgreementHandler),
			Token:       "UserAgreementHandler",
		},
	}
}

func getControllersDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: controllers.NewCategoryController,
			Interface:   new(controllers.ICategoryController),
			Token:       "CategoryController",
		},
		{
			Constructor: controllers.NewCategoryItemController,
			Interface:   new(controllers.ICategoryItemController),
			Token:       "CategoryItemController",
		},
	}
}

func getHttpServerDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: router.NewRouter,
			Interface:   new(router.IRouter),
			Token:       "Router",
		},
		{
			Constructor: http.NewHttpServer,
			Interface:   new(http.IHttpServer),
			Token:       "HttpServer",
		},
	}
}

func getBotDependencies() []Dependency {
	return []Dependency{
		{
			Constructor: bot.NewBot,
			Interface:   new(bot.IBot),
			Token:       "Bot",
		},
	}
}

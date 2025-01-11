package di

import (
	"go.uber.org/dig"
	"shop-bot/config"
	"shop-bot/handlers"
	"shop-bot/handlers/callback_queries"
	"shop-bot/internal/db"
	"shop-bot/internal/logger"
	"shop-bot/repositories"
	"shop-bot/services"
	"shop-bot/utils"
)

func BuildContainer() *dig.Container {
	c := dig.New()

	c = AppendDependenciesToContainer(c, getRequiredDependencies())
	c = AppendDependenciesToContainer(c, getRepositoriesDependencies())
	c = AppendDependenciesToContainer(c, getServicesDependencies())
	c = AppendDependenciesToContainer(c, getHandlersDependencies())

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

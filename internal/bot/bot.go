package bot

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"go.uber.org/dig"
	"shop-bot/config"
	"shop-bot/constants"
	"shop-bot/handlers"
	"shop-bot/handlers/callback_queries"
	"shop-bot/internal/logger"
	"shop-bot/utils"
	"sync"
)

type botDependencies struct {
	dig.In

	ShutdownWaitGroup *sync.WaitGroup `name:"ShutdownWaitGroup"`
	ShutdownContext   context.Context `name:"ShutdownContext"`
	Logger            logger.ILogger  `name:"Logger"`
	Config            config.IConfig  `name:"Config"`

	DefaultHandler       handlers.IDefaultHandler               `name:"DefaultHandler"`
	CommandsHandler      handlers.ICommandHandler               `name:"CommandHandler"`
	MenuHandler          callback_queries.IMenuHandler          `name:"MenuHandler"`
	UserAgreementHandler callback_queries.IUserAgreementHandler `name:"UserAgreementHandler"`
}

type Bot struct {
	shutdownWaitGroup *sync.WaitGroup
	shutdownContext   context.Context

	logger logger.ILogger
	config config.IConfig `name:"Config"`

	bot *tg_bot.Bot

	commandsHandler      handlers.ICommandHandler
	defaultHandler       handlers.IDefaultHandler
	menuHandler          callback_queries.IMenuHandler
	userAgreementHandler callback_queries.IUserAgreementHandler
}

func StartBot(deps botDependencies) {
	defer deps.ShutdownWaitGroup.Done()

	bot := &Bot{
		shutdownWaitGroup:    deps.ShutdownWaitGroup,
		shutdownContext:      deps.ShutdownContext,
		logger:               deps.Logger,
		config:               deps.Config,
		commandsHandler:      deps.CommandsHandler,
		defaultHandler:       deps.DefaultHandler,
		menuHandler:          deps.MenuHandler,
		userAgreementHandler: deps.UserAgreementHandler,
	}

	opts := []tg_bot.Option{
		tg_bot.WithDefaultHandler(bot.defaultHandler.Handle),
		tg_bot.WithMiddlewares(bot.timeoutMiddleware, bot.panicRecoveryMiddleware, bot.logMiddleware),
	}

	tgBot, err := tg_bot.New(deps.Config.BotToken(), opts...)

	utils.PanicIfError(err)

	bot.bot = tgBot

	registerHandlers(bot.bot, deps)

	bot.logger.Log("Bot started")

	bot.bot.Start(bot.shutdownContext)

	select {
	case <-deps.ShutdownContext.Done():
		bot.logger.Log("Shutting down bot gracefully...")
	}

	bot.logger.Log("Bot stopped")
}

func registerHandlers(bot *tg_bot.Bot, deps botDependencies) {
	bot.RegisterHandler(tg_bot.HandlerTypeMessageText, constants.CommandStart, tg_bot.MatchTypeExact, deps.CommandsHandler.Start)

	bot.RegisterHandler(tg_bot.HandlerTypeCallbackQueryData, constants.CallbackDataMenuPrefix, tg_bot.MatchTypePrefix, deps.MenuHandler.Handle)
	bot.RegisterHandler(tg_bot.HandlerTypeCallbackQueryData, constants.CallbackDataUserAgreementPrefix, tg_bot.MatchTypePrefix, deps.UserAgreementHandler.Handle)
}

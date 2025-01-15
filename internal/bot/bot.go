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

type IBot interface {
	Start()
}

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

type bot struct {
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

func NewBot(deps botDependencies) *bot {
	return &bot{
		shutdownWaitGroup:    deps.ShutdownWaitGroup,
		shutdownContext:      deps.ShutdownContext,
		logger:               deps.Logger,
		config:               deps.Config,
		commandsHandler:      deps.CommandsHandler,
		defaultHandler:       deps.DefaultHandler,
		menuHandler:          deps.MenuHandler,
		userAgreementHandler: deps.UserAgreementHandler,
	}
}

func (bot *bot) Start() {
	defer bot.shutdownWaitGroup.Done()

	opts := []tg_bot.Option{
		tg_bot.WithDefaultHandler(bot.defaultHandler.Handle),
		tg_bot.WithMiddlewares(bot.timeoutMiddleware, bot.panicRecoveryMiddleware, bot.logMiddleware),
	}

	tgBot, err := tg_bot.New(bot.config.BotToken(), opts...)

	utils.PanicIfError(err)

	bot.bot = tgBot

	bot.registerHandlers()

	bot.logger.Log("Bot started")

	bot.bot.Start(bot.shutdownContext)

	select {
	case <-bot.shutdownContext.Done():
		bot.logger.Log("Shutting down Bot gracefully...")
	}

	bot.logger.Log("Bot stopped")
}

func (bot *bot) registerHandlers() {
	if bot.bot == nil {
		panic("cannot register handlers without bot instance")
	}

	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeMessageText,
		constants.CommandStart,
		tg_bot.MatchTypeExact,
		bot.commandsHandler.Start,
	)

	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeCallbackQueryData,
		constants.CallbackDataMenuPrefix,
		tg_bot.MatchTypePrefix,
		bot.menuHandler.Handle,
	)

	bot.bot.RegisterHandler(
		tg_bot.HandlerTypeCallbackQueryData,
		constants.CallbackDataUserAgreementPrefix,
		tg_bot.MatchTypePrefix,
		bot.userAgreementHandler.Handle,
	)
}

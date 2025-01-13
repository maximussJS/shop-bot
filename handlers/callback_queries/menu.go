package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"shop-bot/constants"
	"shop-bot/internal/logger"
	"shop-bot/services"
	"shop-bot/utils"
)

type IMenuHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type MenuHandlerDependencies struct {
	dig.In

	Logger      logger.ILogger        `name:"Logger"`
	TextService services.ITextService `name:"TextService"`
	UserService services.IUserService `name:"UserService"`
}

type menuHandler struct {
	logger      logger.ILogger
	textService services.ITextService
	userService services.IUserService
}

func NewMenuHandler(deps MenuHandlerDependencies) *menuHandler {
	return &menuHandler{
		logger:      deps.Logger,
		textService: deps.TextService,
		userService: deps.UserService,
	}
}

func (h *menuHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	answerResult := utils.MustAnswerCallbackQuery(ctx, b, update)

	if !answerResult {
		h.logger.Error(fmt.Sprintf("Failed to answer callback query: %s", update.CallbackQuery.ID))
		utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Message.Chat.ID,
			Text:   "An error occurred while processing the request. Please try again later.",
		})
		return
	}

	h.process(ctx, b, update)
}

func (h *menuHandler) process(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	switch update.CallbackQuery.Data {
	case constants.CallbackDataShowMenu:
		h.showMenu(ctx, b, update)
	}
}

func (h *menuHandler) showMenu(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	userId := utils.GetUserID(update)

	if h.userService.IsUserNotExists(ctx, userId) {
		utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
			Text:      h.textService.UserAgreementNotAccepted(),
			ParseMode: models.ParseModeMarkdown,
		})
		return
	}

	utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "Menu",
	})
}

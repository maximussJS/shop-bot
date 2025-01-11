package callback_queries

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	tg_bot_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"shop-bot/constants"
	"shop-bot/internal/logger"
	"shop-bot/services"
	"shop-bot/utils"
)

type IUserAgreementHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update)
}

type userAgreementHandlerDependencies struct {
	dig.In

	Logger      logger.ILogger        `name:"Logger"`
	TextService services.ITextService `name:"TextService"`
	UserService services.IUserService `name:"UserService"`
}

type UserAgreementHandler struct {
	logger      logger.ILogger
	textService services.ITextService
	userService services.IUserService
}

func NewUserAgreementHandler(deps userAgreementHandlerDependencies) *UserAgreementHandler {
	return &UserAgreementHandler{
		logger:      deps.Logger,
		textService: deps.TextService,
		userService: deps.UserService,
	}
}

func (h *UserAgreementHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update) {
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

func (h *UserAgreementHandler) process(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update) {
	switch update.CallbackQuery.Data {
	case constants.CallbackDataUserAgreementShow:
		h.showUserAgreement(ctx, b, update)
	case constants.CallbackDataUserAgreementAccept:
		h.acceptUserAgreement(ctx, b, update)
	case constants.CallbackDataUserAgreementDecline:
		h.declineUserAgreement(ctx, b, update)
	}
}

func (h *UserAgreementHandler) showUserAgreement(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update) {
	kb := &tg_bot_models.InlineKeyboardMarkup{
		InlineKeyboard: constants.InlineKeyBoardUserAgreement,
	}

	utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        h.textService.UserAgreement(),
		ReplyMarkup: kb,
	})
}

func (h *UserAgreementHandler) acceptUserAgreement(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update) {
	userId := utils.GetUserID(update)
	chatId := utils.GetChatID(update)

	if h.userService.IsUserExists(ctx, userId) {
		utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: chatId,
			Text:   h.textService.UserAgreementAlreadyAccepted(),
		})
		return
	}

	h.userService.CreateFromTelegramUpdate(ctx, update)

	utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID: chatId,
		Text:   h.textService.UserAgreementAccepted(),
	})
}

func (h *UserAgreementHandler) declineUserAgreement(ctx context.Context, b *tg_bot.Bot, update *tg_bot_models.Update) {
	userId := utils.GetUserID(update)
	chatId := utils.GetChatID(update)

	if h.userService.IsUserExists(ctx, userId) {
		utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
			ChatID: chatId,
			Text:   h.textService.UserAgreementAlreadyAccepted(),
		})
		return
	}

	utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
		ChatID: chatId,
		Text:   h.textService.UserAgreementDeclined(),
	})
}

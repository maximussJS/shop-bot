package handlers

import (
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"shop-bot/utils"
)

type IDefaultHandler interface {
	Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type DefaultHandler struct{}

func NewDefaultHandler() *DefaultHandler {
	return &DefaultHandler{}
}

func (h *DefaultHandler) Handle(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &tg_bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Hello, *" + tg_bot.EscapeMarkdown(update.Message.From.FirstName) + "*" + "\n" + "Please use /start command to start the bot",
		ParseMode: models.ParseModeMarkdown,
	})

	utils.PanicIfError(err)
}

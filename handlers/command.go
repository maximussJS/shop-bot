package handlers

import (
	"bytes"
	"context"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"os"
	"shop-bot/constants"
	"shop-bot/services"
	"shop-bot/utils"
)

type ICommandHandler interface {
	Start(ctx context.Context, b *tg_bot.Bot, update *models.Update)
}

type commandHandlerDependencies struct {
	dig.In

	TextService services.ITextService `name:"TextService"`
}

type commandHandler struct {
	startLogoData []byte
	textService   services.ITextService
}

func NewCommandHandler(deps commandHandlerDependencies) *commandHandler {
	fileData, errReadFile := os.ReadFile("./images/facebook.png")
	utils.PanicIfError(errReadFile)

	return &commandHandler{
		startLogoData: fileData,
		textService:   deps.TextService,
	}
}

func (c *commandHandler) Start(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
	captionText := c.textService.WelcomeMessage(update.Message.From.FirstName)

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: constants.InlineKeyboardStart,
	}

	params := &tg_bot.SendPhotoParams{
		ChatID:      update.Message.Chat.ID,
		Photo:       &models.InputFileUpload{Filename: "facebook.png", Data: bytes.NewReader(c.startLogoData)},
		Caption:     captionText,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: kb,
	}

	_, err := b.SendPhoto(ctx, params)

	utils.PanicIfError(err)
}

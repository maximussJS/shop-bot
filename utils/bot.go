package utils

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MustSendMessage(ctx context.Context, b *tg_bot.Bot, params *tg_bot.SendMessageParams) *models.Message {
	msg, err := b.SendMessage(ctx, params)

	PanicIfError(err)

	return msg
}

func MustAnswerCallbackQuery(ctx context.Context, b *tg_bot.Bot, update *models.Update) bool {
	result, err := b.AnswerCallbackQuery(ctx, &tg_bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	PanicIfError(err)

	return result
}

func GetUserID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.From.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}

	panic(fmt.Sprintf("unable to get user id from update: %v", update))
}

func GetChatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}

	panic(fmt.Sprintf("unable to get chat id from update: %v", update))
}

func GetFirstName(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.FirstName
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.FirstName
	}

	panic(fmt.Sprintf("unable to get first name from update: %v", update))
}

func GetLastName(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.LastName
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.LastName
	}

	panic(fmt.Sprintf("unable to get last name from update: %v", update))
}

func GetUsername(update *models.Update) string {
	if update.Message != nil {
		return update.Message.From.Username
	}

	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.Username
	}

	panic(fmt.Sprintf("unable to get username from update: %v", update))
}

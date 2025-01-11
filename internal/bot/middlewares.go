package bot

import (
	"context"
	"fmt"
	tg_bot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"runtime"
	"shop-bot/utils"
)

func (bot *Bot) logMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		if update.Message != nil {
			bot.logger.Log(fmt.Sprintf("%d say: %s", update.Message.From.ID, update.Message.Text))
		}

		if update.CallbackQuery != nil {
			bot.logger.Log(fmt.Sprintf("%d clicked: %s", update.CallbackQuery.From.ID, update.CallbackQuery.Data))
		}

		next(ctx, b, update)
	}
}

func (bot *Bot) requestTimeoutMiddleware(next tg_bot.HandlerFunc) tg_bot.HandlerFunc {
	return func(ctx context.Context, b *tg_bot.Bot, update *models.Update) {
		timeoutDuration := bot.config.RequestTimeout()

		bot.logger.Debug(fmt.Sprintf("Request timeout: %s", timeoutDuration))

		childContext, cancel := context.WithTimeout(ctx, timeoutDuration)
		defer cancel()

		doneCh := make(chan struct{})

		go func() {
			defer func() {
				if err := recover(); err != nil {
					chatId := utils.GetChatID(update)
					stackSize := bot.config.ErrorStackTraceSizeInKb() << 10
					stack := make([]byte, stackSize)
					length := runtime.Stack(stack, true)
					stack = stack[:length]

					if childContext.Err() != nil {
						return
					}

					bot.logger.Error(fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length]))

					b.SendMessage(ctx, &tg_bot.SendMessageParams{
						ChatID: chatId,
						Text:   "An error occurred while processing your request. Please try again later.",
					})
				}
				close(doneCh)
			}()

			next(childContext, b, update)
		}()

		select {
		case <-childContext.Done():
			utils.MustSendMessage(ctx, b, &tg_bot.SendMessageParams{
				ChatID: utils.GetChatID(update),
				Text:   "Request timeout exceeded. Please try again later.",
			})
			return
		case <-doneCh:
			return
		}
	}
}

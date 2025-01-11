package constants

import "github.com/go-telegram/bot/models"

var InlineKeyboardStart = [][]models.InlineKeyboardButton{
	{
		{Text: "📲 Show Menu", CallbackData: CallbackDataShowMenu},
	}, {
		{Text: "📖 User Agreement", CallbackData: CallbackDataUserAgreementShow},
	},
}

var InlineKeyBoardUserAgreement = [][]models.InlineKeyboardButton{
	{
		{Text: "✅ Accept", CallbackData: CallbackDataUserAgreementAccept},
	},
	{
		{Text: "❌ Decline", CallbackData: CallbackDataUserAgreementDecline},
	},
}

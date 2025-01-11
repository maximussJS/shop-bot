package services

import (
	tg_bot "github.com/go-telegram/bot"
	"strings"
)

type ITextService interface {
	WelcomeMessage(firstName string) string

	UserAgreementNotAccepted() string
	UserAgreement() string
	UserAgreementAlreadyAccepted() string
	UserAgreementAccepted() string
	UserAgreementDeclined() string
}

type TextService struct{}

func NewTextService() *TextService {
	return &TextService{}
}

func (s *TextService) WelcomeMessage(firstName string) string {
	var sb strings.Builder

	sb.WriteString("Welcome, *")
	sb.WriteString(tg_bot.EscapeMarkdown(firstName))
	sb.WriteString("*")
	sb.WriteString("\\!\n")
	sb.WriteString("Premium Pharm and products for affiliate marketers\\. \n")
	sb.WriteString("Buy our accounts through a bot automatically 24/7")

	return sb.String()
}

func (s *TextService) UserAgreement() string {
	var sb strings.Builder

	sb.WriteString("*Terms and Conditions*\n\n")

	sb.WriteString("By purchasing products from our store, you automatically accept these terms.\n\n")

	sb.WriteString("All claims regarding returns and exchanges must be made within 24 hours of purchase.\n\n")

	sb.WriteString("Returns and exchanges for Farm accounts are only possible BEFORE linking the payment method or publishing the campaign. Please check the account before linking payment methods!\n\n")

	sb.WriteString("We are not responsible for RISK PAYMENT, POLICY issues, or other reasons for account bans after linking the payment method or publishing the campaign.\n\n")

	sb.WriteString("Returns and exchanges for autorereg are only possible BEFORE logging into the account. Please verify autorereg via the account profile link.\n\n")

	sb.WriteString("Returns, exchanges, and other disputed situations are resolved directly through the manager @bot_maximuss.\n\n")

	sb.WriteString("To resolve disputed situations, please contact @bot_maximuss with your order number and the details of your claim.")

	return sb.String()
}

func (s *TextService) UserAgreementAlreadyAccepted() string {
	var sb strings.Builder

	sb.WriteString("🤔 You have already accepted the terms and conditions. You can use the bot.")

	return sb.String()
}

func (s *TextService) UserAgreementNotAccepted() string {
	var sb strings.Builder

	sb.WriteString("You are not accepted *📖 User Agreement* \n\n")
	sb.WriteString("Please accept the it 😏\n\n")
	sb.WriteString("To accept the agreement\\, follow the instructions below\\:\n\n")
	sb.WriteString("1\\. Tap /start\n")
	sb.WriteString("2\\. Click on the *📖 User Agreement* button\n")
	sb.WriteString("3\\. Click on the *✅ Accept* button\n")
	sb.WriteString("4\\. Use the bot\n")

	return sb.String()
}

func (s *TextService) UserAgreementAccepted() string {
	var sb strings.Builder

	sb.WriteString("😏 You have accepted the terms and conditions. You can now use the bot.")

	return sb.String()
}

func (s *TextService) UserAgreementDeclined() string {
	var sb strings.Builder

	sb.WriteString("😕 You have declined the terms and conditions. You can't use the bot until you accept them.")

	return sb.String()
}

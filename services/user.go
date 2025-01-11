package services

import (
	"context"
	tg_bot_models "github.com/go-telegram/bot/models"
	"go.uber.org/dig"
	"shop-bot/models"
	"shop-bot/repositories"
	"shop-bot/utils"
)

type userServiceDependencies struct {
	dig.In

	UserRepository repositories.IUserRepository `name:"UserRepository"`
}

type IUserService interface {
	IsUserExists(ctx context.Context, userId int64) bool
	IsUserNotExists(ctx context.Context, userId int64) bool
	CreateFromTelegramUpdate(ctx context.Context, update *tg_bot_models.Update) int64
}

type UserService struct {
	userRepository repositories.IUserRepository
}

func NewUserService(deps userServiceDependencies) *UserService {
	return &UserService{
		userRepository: deps.UserRepository,
	}
}

func (s *UserService) IsUserExists(ctx context.Context, userId int64) bool {
	return s.userRepository.GetById(ctx, userId) != nil
}

func (s *UserService) IsUserNotExists(ctx context.Context, userId int64) bool {
	return !s.IsUserExists(ctx, userId)
}

func (s *UserService) CreateFromTelegramUpdate(ctx context.Context, update *tg_bot_models.Update) int64 {
	userId := utils.GetUserID(update)

	s.userRepository.Create(ctx, models.User{
		Id:        userId,
		ChatId:    utils.GetChatID(update),
		Username:  utils.GetUsername(update),
		FirstName: utils.GetFirstName(update),
		LastName:  utils.GetLastName(update),
	})

	return userId
}

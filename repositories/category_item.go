package repositories

import (
	"context"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"shop-bot/config"
	"shop-bot/models"
	"shop-bot/utils"
)

type ICategoryItemRepository interface {
	Create(ctx context.Context, category models.CategoryItem) uint
	GetById(ctx context.Context, id uint) *models.CategoryItem
	GetAll(ctx context.Context, limit, offset int) []models.CategoryItem
	UpdateById(ctx context.Context, id uint, category models.CategoryItem)
	DeleteById(ctx context.Context, id uint)
}

type categoryItemRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type categoryItemRepository struct {
	db *gorm.DB
}

func NewCategoryItemRepository(deps categoryItemRepositoryDependencies) *categoryItemRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.CategoryItem{})

		utils.PanicIfError(err)
	}

	return &categoryItemRepository{
		db: deps.DB,
	}
}

func (r *categoryItemRepository) Create(ctx context.Context, item models.CategoryItem) uint {
	err := r.db.WithContext(ctx).Create(&item).Error

	utils.PanicIfNotContextError(err)

	return item.Id
}

func (r *categoryItemRepository) GetById(ctx context.Context, id uint) *models.CategoryItem {
	var item models.CategoryItem
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Where("id = ?", id).First(&item).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &item
}

func (r *categoryItemRepository) GetAll(ctx context.Context, limit, offset int) []models.CategoryItem {
	var items []models.CategoryItem

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&items).Error

	utils.PanicIfNotContextError(err)

	return items
}

func (r *categoryItemRepository) UpdateById(ctx context.Context, id uint, item models.CategoryItem) {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(&item).Error

	utils.PanicIfNotContextError(err)
}

func (r *categoryItemRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.CategoryItem{}).Error

	utils.PanicIfNotContextError(err)
}

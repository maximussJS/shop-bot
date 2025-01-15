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

type ICategoryRepository interface {
	Create(ctx context.Context, category models.Category) uint
	GetById(ctx context.Context, id uint) *models.Category
	GetAll(ctx context.Context, limit, offset int) []models.Category
	GetByName(ctx context.Context, name string) *models.Category
	UpdateById(ctx context.Context, id uint, category models.Category)
	DeleteById(ctx context.Context, id uint)
}

type categoryRepositoryDependencies struct {
	dig.In

	DB     *gorm.DB       `name:"DB"`
	Config config.IConfig `name:"Config"`
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(deps categoryRepositoryDependencies) *categoryRepository {
	if deps.Config.RunMigrations() {
		err := deps.DB.AutoMigrate(&models.Category{})

		utils.PanicIfError(err)
	}

	return &categoryRepository{
		db: deps.DB,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category models.Category) uint {
	err := r.db.WithContext(ctx).Create(&category).Error

	utils.PanicIfNotContextError(err)

	return category.Id
}

func (r *categoryRepository) GetById(ctx context.Context, id uint) *models.Category {
	var category models.Category
	err := r.db.WithContext(ctx).Clauses(clause.Returning{}).Preload("Items").Where("id = ?", id).First(&category).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &category
}

func (r *categoryRepository) GetAll(ctx context.Context, limit, offset int) []models.Category {
	var categories []models.Category

	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Preload("Items").Find(&categories).Error

	utils.PanicIfNotContextError(err)

	return categories
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) *models.Category {
	var category models.Category
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error

	if err != nil && utils.IsRecordNotFoundError(err) {
		return nil
	}

	utils.PanicIfNotRecordNotFound(err)

	return &category
}

func (r *categoryRepository) UpdateById(ctx context.Context, id uint, category models.Category) {
	err := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(&category).Error

	utils.PanicIfNotContextError(err)
}

func (r *categoryRepository) DeleteById(ctx context.Context, id uint) {
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Category{}).Error

	utils.PanicIfNotContextError(err)
}

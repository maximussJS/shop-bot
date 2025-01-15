package models

import (
	"gorm.io/gorm"
	"time"
)

type CategoryItem struct {
	Id         uint      `gorm:"unique;primaryKey;autoIncrement" json:"id" `
	Url        string    `gorm:"not null" json:"url"`
	CategoryId uint      `gorm:"not null;index" json:"categoryId"`
	Category   Category  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (p *CategoryItem) TableName() string {
	return "category_items"
}

func (p *CategoryItem) BeforeCreate(tx *gorm.DB) (err error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return
}

func (p *CategoryItem) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return
}

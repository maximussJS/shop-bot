package models

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id        uint           `gorm:"primaryKey;autoIncrement" json:"id" `
	Name      string         `gorm:"size:255;not null;unique" json:"name"`
	Price     float64        `gorm:"not null" json:"price"`
	Items     []CategoryItem `gorm:"foreignKey:CategoryId" json:"items"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return
}

func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageUserActivity struct {
	ID        string    `gorm:"column:ID"`
	UserId    string    `gorm:"column:UserID"`
	ImageId   string    `gorm:"column:ImageID"`
	IsLike    bool      `gorm:"column:IsLike"`
	IsComment bool      `gorm:"column:IsComment"`
	CreatedAt time.Time `gorm:"column:CreatedAt; autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt; autoUpdateTime"`
}

func (ImageUserActivity) TableName() string {
	return "ImageUserActivity"
}

func (i *ImageUserActivity) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return
}

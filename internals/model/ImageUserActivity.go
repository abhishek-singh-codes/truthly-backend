package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImageUserActivity struct {
	// primary key
	ID string `gorm:"column:ID;primaryKey"`

	// foreign keys
	ImageId string `gorm:"column:ImageID;not null"`
	UserId  string `gorm:"column:UserID;not null"`

	// activity — tinyint, default false
	IsLike    bool `gorm:"column:IsLike;default:0"`
	IsComment bool `gorm:"column:IsComment;default:0"`

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (ImageUserActivity) TableName() string {
	return "ImageUserActivity"
}

func (i *ImageUserActivity) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		i.ID = id.String()
	}
	return
}

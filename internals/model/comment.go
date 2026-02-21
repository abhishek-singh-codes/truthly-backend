package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Commemts struct {
	// pk
	CommentId string `gorm:"column:CommentId; primaryKey"`

	// fk
	UserId        string `gorm:"column:UserId"`
	ImageId       string `gorm:"column:ImageId"`
	DescriptionId string `gorm:"column:DescriptionId"`
	AnalyticId    string `gorm:"column:AnalyticId"`

	// initially it will be empty
	Comment string `gorm:"column:Comment"`

	// dates
	CreatedAt time.Time `gorm:"column:CreatedAt; autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt; autoUpdateTime"`
}

func (Commemts) TableName() string {
	return "Comments"
}

func (c *Commemts) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CommentId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		c.CommentId = id.String()
	}
	return
}

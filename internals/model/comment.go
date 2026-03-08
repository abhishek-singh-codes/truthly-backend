package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comment struct { // ← typo fix kiya "Commemts" → "Comment"
	// primary key
	CommentId string `gorm:"column:CommentId;primaryKey"`

	// foreign keys
	ImageId string `gorm:"column:ImageId;not null"`
	UserId  string `gorm:"column:UserId;not null"`

	// content — "Commet" schema ka column name hai (typo DB mein hai)
	Comment *string `gorm:"column:Commet;type:text"`

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Comment) TableName() string {
	return "Comments"
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	if c.CommentId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		c.CommentId = id.String()
	}
	return
}

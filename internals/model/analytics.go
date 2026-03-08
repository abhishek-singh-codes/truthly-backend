package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Analytic struct {
	// primary key
	AnalyticId string `gorm:"column:AnalyticId;primaryKey"`

	// foreign key
	ImageId string `gorm:"column:ImageId;not null"`

	// counts — default 0, NULL nahi
	LikeCount    int `gorm:"column:LikeCount;default:0"`
	ShareCount   int `gorm:"column:ShareCount;default:0"`
	CommentCount int `gorm:"column:CommentCount;default:0"`

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Analytic) TableName() string {
	return "Analytics"
}

func (a *Analytic) BeforeCreate(tx *gorm.DB) (err error) {
	if a.AnalyticId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		a.AnalyticId = id.String()
	}
	return
}

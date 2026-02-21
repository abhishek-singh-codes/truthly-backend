package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Analytic struct {

	// fk
	AnalyticId    string `gorm:"column:AnalyticId; primaryKey"`
	ImageId       string `gorm:"column:ImageId"`
	DescriptionId string `gorm:"column:DescriptionId"`
	UserId        string `gorm:"column:UserId"`

	// Initially these field will be 0
	Like    int `gorm:"column:LikeCount"`
	Share   int `gorm:"column:ShareCount"`
	Comment int `gorm:"column:CommentCount"`

	// dates
	CreatedAt time.Time `gorm:"column:CreatedAt; autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt; autoUpdateTime"`
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

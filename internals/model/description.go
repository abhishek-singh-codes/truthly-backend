package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Description struct {
	// primary key
	DescriptionId string `gorm:"column:DescriptionId;primaryKey"`

	// foreign key
	ImageId string `gorm:"column:ImageId;not null"`

	// content
	Description *string `gorm:"column:Description;type:text"`
	Country     *string `gorm:"column:Country"`
	State       *string `gorm:"column:State"`
	City        *string `gorm:"column:City"`

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Description) TableName() string {
	return "Description"
}

func (d *Description) BeforeCreate(tx *gorm.DB) (err error) {
	if d.DescriptionId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		d.DescriptionId = id.String()
	}
	return
}

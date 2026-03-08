package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {

	// primary key
	ImageId string `gorm:"column:ImageId;primaryKey"`

	// foreign key
	UserId string `gorm:"column:UserId;not null"`

	// s3 bucket image url
	ImageUrl string `gorm:"column:ImageUrl;not null;type:text"`

	// lat long — nullable kyunki user location na de
	Latitude  *float64 `gorm:"column:Latitude;type:decimal(9,6)"`  // ← (10,8) fix kiya
	Longitude *float64 `gorm:"column:Longitude;type:decimal(9,6)"` // ← (11,8) fix kiya

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (Image) TableName() string {
	return "Images"
}

func (i *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ImageId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		i.ImageId = id.String()
	}
	return
}

package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId string `gorm:"column:UserId;primaryKey"`

	UserName  string  `gorm:"column:UserName;not null"`
	FirstName *string `gorm:"column:FirstName"`
	LastName  *string `gorm:"column:LastName"`

	Age    *int    `gorm:"column:Age"`
	Gender *string `gorm:"column:Gender;type:enum('Male','Female','Others')"`

	Country      *string `gorm:"column:Country"`
	State        *string `gorm:"column:State"`
	City         *string `gorm:"column:City"`
	AddressLine1 *string `gorm:"column:AddressLine1"`
	AddressLine2 *string `gorm:"column:AddressLine2"`
	PostalCode   *string `gorm:"column:PostalCode"`

	Email        *string `gorm:"column:Email;uniqueIndex"`
	Password     string  `gorm:"column:Password"`
	MobileNumber *string `gorm:"column:MobileNumber;uniqueIndex"`

	ProfileImage *string `gorm:"column:ProfileImage;type:text"`

	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt *time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

// GORM ko table name batao
func (User) TableName() string {
	return "Users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserId == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.UserId = id.String()
	}
	return nil
}

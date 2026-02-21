package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId uuid.UUID `gorm:"column:UserId;"`

	UserName string `gorm:"column:UserName;"`

	FirstName string `gorm:"column:FirstName;"`
	LastName  string `gorm:"column:LastName;"`

	Age int `gorm:"column:Age;"`

	Gender string `gorm:"column:Gender;"`

	Country string `gorm:"column:Country;"`
	State   string `gorm:"column:State;"`
	City    string `gorm:"column:City;"`
	Address string `gorm:"column:Address"`

	Email        string `gorm:"column:Email;"`
	Password     string `gorm:"column:Password;"` // hash stored
	MobileNumber string `gorm:"column:MobileNumber;"`

	CreatedAt time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (User) TableName() string {
	return "Users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserId == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return nil
		}
		u.UserId = id
	}
	return nil
}

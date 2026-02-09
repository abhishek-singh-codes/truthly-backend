package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserId string `gorm:"column:UserId;primaryKey"`

	UserName string `gorm:"column:UserName;unique;"`

	FirstName string `gorm:"column:FirstName;"`
	LastName  string `gorm:"column:LastName;"`

	Age int `gorm:"column:Age;"`

	Gender string `gorm:"column:Gender;"`

	Country string `gorm:"column:Country;"`
	State   string `gorm:"column:State;"`
	City    string `gorm:"column:City;"`
	Address string `gorm:"column:Address"`

	Email        string `gorm:"column:Email;unique;"`
	Password     string `gorm:"column:Password;"` // hash stored
	MobileNumber string `gorm:"column:MobileNumber;unique;"`

	CreatedAt time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (User) TableName() string {
	return "Users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserId == "" {
		u.UserId = uuid.New().String()
	}
	return nil
}

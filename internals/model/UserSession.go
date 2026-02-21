package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	Id        string    `gorm:"column:Id"`
	UserId    string    `gorm:"column:UserId"`
	UserName  string    `gorm:"column:UserName"`
	SessionId string    `gorm:"column:SessionId"`
	Status    string    `gorm:"column:Status"`
	CreatedAt time.Time `gorm:"column:CreatedAt"`
	ExpiredAt time.Time `gorm:"column:ExpiredAt"`
}

func (UserSession) TableName() string {
	return "UserSessions"
}

func (u *UserSession) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == "" {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		u.Id = id.String()
	}
	return nil
}

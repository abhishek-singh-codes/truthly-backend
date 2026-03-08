package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	// primary key
	Id string `gorm:"column:Id;primaryKey"`

	// user info
	UserId   string `gorm:"column:UserId;not null"`
	UserName string `gorm:"column:UserName;not null"`

	// session
	SessionId string `gorm:"column:SessionId;not null"`
	Status    string `gorm:"column:Status;not null"`

	// dates
	CreatedAt *time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	ExpiredAt *time.Time `gorm:"column:ExpiredAt"` // ← nullable, session expire na hua ho
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

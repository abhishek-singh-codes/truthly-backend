package dto

import (
	"truthly/internals/model"
)

type UserRequestDto struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserResponseDto struct {
	Message string `json:"message,omitempty"`
	UserId  string `json:"userId"`
}

type UserDetailsForHome struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	UserName  string `json:"userName,omitempty"`

	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

// DTO → Model

func ToModel(u *UserRequestDto) *model.User {
	return &model.User{

		UserName: u.UserName,
		Password: u.Password,
	}
}

// Model to dto I will do manually

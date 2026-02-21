package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"truthly/internals/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	// Inser a new user
	CreatNewUser(ctx context.Context, user *model.User) (*model.User, error)
	VerifyUser(ctx context.Context, userName string, password string) (string, error)

	GetUserById(ctx context.Context, userId string) (*model.User, error)
}

type userRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

// constructor to get the userRepository
func GetUserRepo(l *slog.Logger, db *gorm.DB) UserRepository {
	return &userRepository{
		db:     db,
		logger: l,
	}
}

// Insert a new user -> signup
func (ur *userRepository) CreatNewUser(ctx context.Context, user *model.User) (*model.User, error) {
	if err := ur.db.WithContext(ctx).Create(user).Error; err != nil {
		ur.logger.Error(err.Error())
		return nil, err
	}
	return user, nil
}

// validate user
func (ur *userRepository) VerifyUser(ctx context.Context, userName string, password string) (string, error) {

	var user model.User

	err := ur.db.WithContext(ctx).
		Where("UserName = ?", userName).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("user not found")
		}

		// actual db error
		ur.logger.Error(err.Error())
		return "", err
	}

	// password check
	if user.Password != password {
		return "", fmt.Errorf("invalid password")
	}

	// success → return userId
	return user.UserId, nil
}

// Get user details by user id to show on home page
func (ur *userRepository) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	var user model.User

	err := ur.db.WithContext(ctx).
		Where("UserId = ?", userId).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.logger.Error("User not found", "userId", userId)
			return nil, err
		}
		ur.logger.Error(err.Error())
		return nil, err
	}

	return &user, nil
}

package service

import (
	"context"
	"log/slog"
	"truthly/internals/dto"
	"truthly/internals/repository"
)

type ProfileService interface {
	ImagesByUserId(
		ctx context.Context,
		userId string,
	) (*dto.ProfileResponseDto, error)
}

type profileService struct {
	logger    *slog.Logger
	imageRepo repository.ImageRepository
	userRepo  repository.UserRepository
}

// constructor
func GetNewProfileService(
	l *slog.Logger,
	ir repository.ImageRepository,
	ur repository.UserRepository,
) ProfileService {
	return &profileService{
		logger:    l,
		imageRepo: ir,
		userRepo:  ur,
	}
}

// Get The Images On Profile Page For RequestID

func (ps *profileService) ImagesByUserId(
	ctx context.Context,
	userId string,
) (*dto.ProfileResponseDto, error) {

	ps.logger.Info("Getting images for userId", "userId", userId)

	//1. get images
	images, err := ps.imageRepo.FetchImagesByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	ps.logger.Info("Total fetched images", "count", len(images))

	//2. bio and userName from user repo
	user, err := ps.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponseDto{
		Images:   images,
		UserName: user.UserName,
	}, nil
}

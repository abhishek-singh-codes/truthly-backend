package service

import (
	"context"
	"log/slog"
	"truthly/internals/repository"
)

type InteractionService interface {
	LikeImage(ctx context.Context, userId, imageId string) error
	AddComment(ctx context.Context, userId, imageID, text string) error
}

type interactionService struct {
	logger          *slog.Logger
	interactionRepo repository.InteractionRepository
	analyticsRepo   repository.AnalyticRepository
}

func GetNewInteractionService(logger *slog.Logger, ir repository.InteractionRepository, analyticRepo repository.AnalyticRepository) InteractionService {
	return &interactionService{
		logger:          logger,
		interactionRepo: ir,
		analyticsRepo:   analyticRepo,
	}
}

func (s *interactionService) LikeImage(ctx context.Context, userId, imageId string) error {
	//1. update the db
	err := s.interactionRepo.LikeImage(ctx, userId, imageId)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	// 2. read updated count
	_, err = s.analyticsRepo.GetAnalyticsByImageId(ctx, imageId)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *interactionService) AddComment(ctx context.Context, userId, imageId, text string) error {

	// 1. Update in the db
	err := s.interactionRepo.AddComment(ctx, userId, imageId, text)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	// 2. read comment count
	_, err = s.analyticsRepo.GetAnalyticsByImageId(ctx, imageId)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

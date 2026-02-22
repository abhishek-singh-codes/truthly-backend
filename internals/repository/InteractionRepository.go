package repository

import (
	"context"
	"errors"
	"log/slog"
	"truthly/internals/model"

	"gorm.io/gorm"
)

type InteractionRepository interface {
	// to incrase the like on image
	LikeImage(ctx context.Context, userId, imageId string) error
	AddComment(ctx context.Context, userId, imageId string, text *string) error
	UnlikeImage(
		ctx context.Context,
		userId string,
		imageId string,
	) error
}

type interactionRepository struct {
	logger *slog.Logger
	Db     *gorm.DB
}

func GetNewInteractionRepository(db *gorm.DB, logger *slog.Logger) InteractionRepository {
	return &interactionRepository{
		Db:     db,
		logger: logger,
	}
}

// helper function increase the like
func (r *interactionRepository) incrementLikeCount(tx *gorm.DB, imageId string) error {
	return tx.
		Model(&model.Analytic{}).
		Where("ImageId=?", imageId).
		UpdateColumn("LikeCount", gorm.Expr("LikeCount + ?", 1)).Error
}

func (r *interactionRepository) decrementLikeCount(
	tx *gorm.DB,
	imageId string,
) error {

	return tx.
		Model(&model.Analytic{}).
		Where("ImageID = ? AND LikeCount > 0", imageId).
		UpdateColumn("LikeCount", gorm.Expr("LikeCount - 1")).Error
}

func (r *interactionRepository) LikeImage(ctx context.Context, userId, imageId string) error {
	return r.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. check if interaction row already exist
		var imageActivity model.ImageUserActivity

		err := tx.
			Where("UserID=? AND ImageID=?", userId, imageId).
			First(&imageActivity).Error

		// Case now row exist than create a row and set true
		if errors.Is(err, gorm.ErrRecordNotFound) {

			newImageUserActivity := model.ImageUserActivity{
				UserId:  userId,
				ImageId: imageId,
				IsLike:  true,
			}

			if err := tx.Create(&newImageUserActivity).Error; err != nil {
				return err
			}

			return r.incrementLikeCount(tx, imageId)
		}

		if err != nil {
			return err
		}

		// row exist and likes
		if imageActivity.IsLike {
			return nil
		}

		/* if row already exists but IsLike = false ,
		   Lets say you have previously lik -> unlike -> Now again like
		*/

		if err := tx.
			Model(&model.ImageUserActivity{}).
			Where("ID", imageActivity.ID).
			Update("IsLike", true).Error; err != nil {
			return err
		}
		return r.incrementLikeCount(tx, imageId)
	})
}

func (r *interactionRepository) UnlikeImage(
	ctx context.Context,
	userId string,
	imageId string,
) error {

	return r.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 1. find interaction row
		var imageActivity model.ImageUserActivity

		err := tx.
			Where("UserID=? AND ImageID=?", userId, imageId).
			First(&imageActivity).Error

		// Case 1: no record → already unliked
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		// Case 2: already unliked
		if !imageActivity.IsLike {
			return nil
		}

		// Case 3: set IsLike = false
		if err := tx.
			Model(&model.ImageUserActivity{}).
			Where("ID", imageActivity.ID).
			Update("IsLike", false).Error; err != nil {
			return err
		}

		// 4. decrement like count
		return r.decrementLikeCount(tx, imageId)
	})
}

func (r *interactionRepository) AddComment(ctx context.Context, userId, imageId string, text *string) error {
	return r.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//1. Get Analytics row
		var analytic *model.Analytic

		err := tx.Where("Imageid = ?", imageId).
			First(&analytic).Error

		if err != nil {
			r.logger.Error("Aanalytics row not found", "error", err.Error, "imageId", imageId)
			return err
		}

		// 2. Add comment
		comment := &model.Commemts{
			UserId:  userId,
			ImageId: imageId,

			DescriptionId: analytic.DescriptionId,
			AnalyticId:    analytic.AnalyticId,
			Comment:       text,
		}

		if err := tx.Create(comment).Error; err != nil {
			r.logger.Error("failed to create comment", "error", err)
			return err
		}

		// 3. Increment comment count
		if err = tx.
			Model(&model.Analytic{}).
			Where("AnalyticId = ?", analytic.AnalyticId).
			UpdateColumn("CommentCount", gorm.Expr("CommentCount + ?", 1)).Error; err != nil {
			r.logger.Error("failed to increase comment count", "error", err.Error())
			return err
		}

		return nil
	})
}

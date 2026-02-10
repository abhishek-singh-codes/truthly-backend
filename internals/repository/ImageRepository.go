package repository

import (
	"context"
	"log/slog"
	"truthly/internals/model"

	"gorm.io/gorm"
)

type ImageRepository interface {
	InsertNewImage(ctx context.Context, data *model.Image) (*model.Image, error)
	FetchImagesByUserId(
		ctx context.Context,
		userId string,
	) ([]string, error)
}

type imageRepository struct {
	Db     *gorm.DB
	logger *slog.Logger
}

// constructor
func GetImageRepo(Db *gorm.DB, logger *slog.Logger) ImageRepository {
	return &imageRepository{
		Db:     Db,
		logger: logger,
	}
}

func (i *imageRepository) InsertNewImage(ctx context.Context, data *model.Image) (*model.Image, error) {
	if err := i.Db.WithContext(ctx).Create(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

// fectch images by the userId
func (i *imageRepository) FetchImagesByUserId(
	ctx context.Context,
	userId string,
) ([]string, error) {

	i.logger.Info("Getting uploaded images", "userId", userId)

	var images []string

	const query = `
		SELECT ImageUrl
		FROM Images
		WHERE UserId = ?
	`

	err := i.Db.
		WithContext(ctx).
		Raw(query, userId).
		Scan(&images).Error

	if err != nil {
		i.logger.Error(err.Error())
		return nil, err
	}

	return images, nil
}

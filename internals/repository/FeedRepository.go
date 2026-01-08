package repository

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type FeedRepository interface {
	GetFeedItems(ctx context.Context, limit int, cursor string) ([]feedRow, string, bool, error)
}

type feedRow struct {
	ImageId   string
	ImageUrl  string
	Caption   string
	CreatedAt time.Time

	UserName string
	UserId   string

	Country string
	City    string
	State   string

	LikeCount    string
	CommentCount string
	ShareCount   string

	IsLike    bool
	IsComment bool
}

type feedRepository struct {
	Db     *gorm.DB
	logger *slog.Logger
}

func GetNewFeedRepository(db *gorm.DB, logger *slog.Logger) FeedRepository {
	return &feedRepository{
		Db:     db,
		logger: logger,
	}
}

// return array of feedRow.  [{}, {}, {}, {}....]
func (fr *feedRepository) GetFeedItems(ctx context.Context, limit int, cursor string) ([]feedRow, string, bool, error) {

	var rows []feedRow

	query := `
			SELECT
				i.ImageId,
				i.ImageUrl,
				i.CreatedAt,

				d.Description AS Caption,
				d.Country,
				d.State,
				d.City,

				a.LikeCount,
				a.CommentCount,
				a.ShareCount,

				u.UserId,
				u.UserName
			FROM Images i
			LEFT JOIN Users u ON u.UserId = i.UserId
			LEFT JOIN Descriptions d ON d.ImageId = i.ImageId
			LEFT JOIN Analytics a ON a.ImageId = i.ImageId
			ORDER BY i.CreatedAt DESC
			LIMIT ?

	`

	err := fr.Db.WithContext(ctx).Raw(query, limit+1).Scan(&rows).Error

	if err != nil {
		fr.logger.Error(err.Error())
		return nil, "", false, err
	}

	hasMore := false
	if len(rows) > limit {
		hasMore = true
		rows = rows[:limit]
	}

	var nextCursor string
	if hasMore {
		nextCursor = rows[len(rows)-1].CreatedAt.Format(time.RFC3339)
	}

	return rows, nextCursor, hasMore, nil
}

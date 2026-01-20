package service

import (
	"context"
	"log/slog"
	"strconv"
	"truthly/internals/dto"
	"truthly/internals/repository"
)

func mustAtoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// utility function that will transform data FeedItemDto --> FeedResponseDto
func conversion(
	rows []repository.FeedRow,
) []dto.FeedItemDto {
	items := make([]dto.FeedItemDto, 0, len(rows))
	for _, r := range rows {
		items = append(items, dto.FeedItemDto{
			ImageId:   r.ImageId,
			ImageUrl:  r.ImageUrl,
			Caption:   r.Caption,
			CreatedAt: r.CreatedAt,

			UserName: r.UserName,
			UserId:   r.UserId,
			// these grouping is difference than feedRow struct
			Location: dto.LocationDto{
				City:    r.City,
				State:   r.State,
				Country: r.Country,
			},
			Analytics: dto.AnalyticsDto{
				Like:      mustAtoi(r.LikeCount),
				Comment:   mustAtoi(r.CommentCount),
				Share:     mustAtoi(r.ShareCount),
				IsLike:    r.IsLike,
				IsComment: r.IsComment,
			},
		})
	}
	return items
}

type FeedService interface {
	GetFeed(ctx context.Context, limit int, cursor string, userId string) (*dto.FeedResponseDto, error)
	GetFeedByRange(
		ctx context.Context,
		reqData *dto.GetFeedByRangeDto,
		userId string,
	) (*dto.FeedResponseDto, error)
}

type feedService struct {
	feedRepo repository.FeedRepository
	logger   *slog.Logger
	// Added geo repo to interact with tile38
	geoRepository repository.GeoRepository
}

func GetNewFeedService(fr repository.FeedRepository, l *slog.Logger, gr repository.GeoRepository) FeedService {
	return &feedService{
		feedRepo:      fr,
		logger:        l,
		geoRepository: gr,
	}
}

func (fs *feedService) GetFeed(ctx context.Context, limit int, cursor string, userId string) (*dto.FeedResponseDto, error) {

	//1. Get FeedRows [{imageId: xyz}, {imageId: abc},]
	rows, nextCursor, hasMore, err := fs.feedRepo.GetFeedItems(ctx, limit, cursor, userId)
	if err != nil {
		fs.logger.Error(err.Error())
		return nil, err
	}

	//2. Image Ids collect
	imageIds := make([]string, 0, len(rows))
	for _, r := range rows {
		imageIds = append(imageIds, r.ImageId)
	}

	// call the conversion function
	items := conversion(rows)

	fs.logger.Info("Feed responseded")

	return &dto.FeedResponseDto{
		Items: items,
		Pagination: dto.PaginationDto{
			NextCursor: nextCursor,
			HasMore:    hasMore,
		},
	}, nil
}

func (fs *feedService) GetFeedByRange(
	ctx context.Context,
	reqData *dto.GetFeedByRangeDto,
	userId string,
) (*dto.FeedResponseDto, error) {
	// 1. Get the imageIds from tile38
	imageIds, nextCursor, hasMore, err := fs.geoRepository.NearByImages(
		ctx,
		reqData.Long, reqData.Lat, reqData.Radius,
		reqData.UserId, reqData.Cursor, reqData.Limit,
	)
	if err != nil {
		fs.logger.Error(err.Error())
		return nil, err
	}

	// 2. Get data from feed Repo
	rows, err := fs.feedRepo.GetFeedItemsByImageIds(
		ctx,
		userId,
		imageIds,
	)
	if err != nil {
		fs.logger.Error("Error in getting geed data by imagedIds", "userId", userId)
		return nil, err
	}

	// convert data into FeedItemsDto
	// call the conversion function
	items := conversion(rows)

	fs.logger.Info("Feed By range responseded", "nextCursor", nextCursor)

	return &dto.FeedResponseDto{
		Items: items,
		Pagination: dto.PaginationDto{
			NextCursor: strconv.Itoa(nextCursor),
			HasMore:    hasMore,
		},
	}, nil
}

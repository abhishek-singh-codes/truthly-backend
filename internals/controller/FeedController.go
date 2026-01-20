package controller

import (
	"log/slog"
	"strconv"
	"truthly/internals/dto"
	"truthly/internals/service"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
	logger      *slog.Logger
	feedService service.FeedService
}

func GetNewFeedController(l *slog.Logger, fs service.FeedService) *FeedController {
	return &FeedController{
		logger:      l,
		feedService: fs,
	}
}

func (c *FeedController) GetFeed(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "4"))
	cursor := ctx.Query("cursor")
	userId := ctx.GetString("userId")

	resp, err := c.feedService.GetFeed(ctx, limit, cursor, userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.ResponseDto[any]{
		Status:    "sucess",
		Message:   "paginated feed",
		ResultObj: resp,
	})
}

// Method for GetFeedByRange, This function will take the data from tile38
func (c *FeedController) GetFeedByRange(ctx *gin.Context) {

	lat, err := strconv.ParseFloat(ctx.Query("lat"), 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid or missing lat"})
		return
	}

	long, err := strconv.ParseFloat(ctx.Query("long"), 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid or missing lon"})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "4"))
	cursor, _ := strconv.Atoi(ctx.DefaultQuery("cursor", "0"))
	radius, _ := strconv.Atoi(ctx.DefaultQuery("radius", "20"))

	userId := ctx.GetString("userId")

	requestDto := &dto.GetFeedByRangeDto{
		Limit:      limit,
		Cursor:     cursor,
		UserId:     userId,
		Radius:     float64(radius),
		Collection: "images",
		Long:       long,
		Lat:        lat,
	}

	c.logger.Info("Request data for feed by range",
		"requestData", requestDto,
	)

	resp, err := c.feedService.GetFeedByRange(ctx, requestDto, userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.ResponseDto[any]{
		Status:    "success",
		Message:   "Paginated Feed By Range",
		ResultObj: resp,
	})
}

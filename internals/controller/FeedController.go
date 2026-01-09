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

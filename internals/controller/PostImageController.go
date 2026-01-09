package controller

import (
	"log/slog"
	"truthly/internals/dto"
	"truthly/internals/service"

	"github.com/gin-gonic/gin"
)

type PostImageController struct {
	logger      *slog.Logger
	postService service.PostService
}

// constructor
func GetNewPostImageController(logger *slog.Logger, postService service.PostService) *PostImageController {
	return &PostImageController{
		logger:      logger,
		postService: postService,
	}
}

func (h *PostImageController) PostImage(ctx *gin.Context) {
	userId := ctx.GetString("userId")

	var postReqDto dto.PostRequestDto
	if err := ctx.ShouldBind(&postReqDto); err != nil {
		h.logger.Error(err.Error())
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//  Reverse geocode only if lat/lng present
	if postReqDto.Latitude != 0 && postReqDto.Longitude != 0 {

		extractor := service.GetNewExtractLocation(h.logger)

		location, err := extractor.FromLatLong(
			postReqDto.Latitude,
			postReqDto.Longitude,
		)

		if err != nil {
			h.logger.Error(
				"failed to extract location",
				"error", err,
			)
		} else {
			postReqDto.City = location.City
			postReqDto.State = location.State
			postReqDto.Country = location.Country
		}
	}

	resp, err := h.postService.UploadPost(ctx, &postReqDto, userId)
	if err != nil {
		h.logger.Error(err.Error())
		ctx.JSON(500, resp)
		return
	}

	ctx.JSON(200, resp)
}

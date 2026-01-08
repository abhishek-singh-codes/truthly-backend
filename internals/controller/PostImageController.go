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
	// 1. Read values from dto
	var postReqDto dto.PostRequestDto
	if err := ctx.ShouldBind(&postReqDto); err != nil {
		h.logger.Error(err.Error())
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 2. Call the service layer
	resp, err := h.postService.UploadPost(ctx, &postReqDto, userId)
	if err != nil {
		h.logger.Error(err.Error())
		ctx.JSON(500, resp)
		return
	}

	// 3. return the response
	ctx.JSON(200, resp)
}

package controller

import (
	"log/slog"
	"truthly/internals/service"

	"github.com/gin-gonic/gin"
)

type InteractionController struct {
	logger             *slog.Logger
	interactionService service.InteractionService
}

func GetNewInteractionController(logger *slog.Logger, is service.InteractionService) *InteractionController {
	return &InteractionController{
		logger:             logger,
		interactionService: is,
	}
}

// POST /api/v1/interactions/images/:imageId/like

func (c *InteractionController) LikeImage(ctx *gin.Context) {

	// 1. imageId from URL
	imageId := ctx.Param("imageId")
	if imageId == "" {
		ctx.JSON(400, gin.H{
			"error": "imageId is required",
		})
		return
	}

	// 2. userId from context (auth middleware)
	userId := ctx.GetString("userId")
	if userId == "" {
		ctx.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// 3. call service
	err := c.interactionService.LikeImage(ctx.Request.Context(), userId, imageId)
	if err != nil {
		c.logger.Error("failed to like image", "error", err)

		ctx.JSON(500, gin.H{
			"error": "failed to like image",
		})
		return
	}

	// 4. success response
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "image liked",
	})
}

// POST /api/v1/interactions/images/:imageId/comments
type addCommentReq struct {
	Comment *string `json:"comment"`
}

func (c *InteractionController) AddComment(ctx *gin.Context) {

	// 1. imageId from URL
	imageId := ctx.Param("imageId")
	if imageId == "" {
		ctx.JSON(400, gin.H{
			"error": "imageId is required",
		})
		return
	}

	// 2. userId from context
	userId := ctx.GetString("userId")
	if userId == "" {
		ctx.JSON(401, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// 3. parse body
	var req addCommentReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if req.Comment == nil {
		ctx.JSON(400, gin.H{
			"error": "comment cannot be empty",
		})
		return
	}

	// 4. call service
	err := c.interactionService.AddComment(
		ctx.Request.Context(),
		userId,
		imageId,
		req.Comment,
	)
	if err != nil {
		c.logger.Error("failed to add comment", "error", err)

		ctx.JSON(500, gin.H{
			"error": "failed to add comment",
		})
		return
	}

	// 5. success response
	ctx.JSON(201, gin.H{
		"status":  "success",
		"message": "comment added",
	})
}

// DELETE /api/v1/interactions/images/:imageId/like
func (c *InteractionController) UnlikeImage(ctx *gin.Context) {

	// 1. imageId
	imageId := ctx.Param("imageId")
	if imageId == "" {
		ctx.JSON(400, gin.H{"error": "imageId is required"})
		return
	}

	// 2. userId
	userId := ctx.GetString("userId")
	if userId == "" {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	// 3. service call
	err := c.interactionService.UnlikeImage(
		ctx.Request.Context(),
		userId,
		imageId,
	)
	if err != nil {
		c.logger.Error("failed to unlike image", "error", err)
		ctx.JSON(500, gin.H{"error": "failed to unlike image"})
		return
	}

	// 4. success
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "image unliked",
	})
}

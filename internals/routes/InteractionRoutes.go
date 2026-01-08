package routes

import (
	"truthly/internals/controller"
	"truthly/internals/middleware"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
)

type InteractionRoutes struct {
	interactionController *controller.InteractionController
	authMiddleware        gin.HandlerFunc
}

func GetNewInteractionRoutes(ic *controller.InteractionController, authToken *auth.AuthToken) *InteractionRoutes {
	return &InteractionRoutes{
		interactionController: ic,
		authMiddleware:        middleware.AuthMiddleware(authToken),
	}
}

// POST /api/v1/interactions/images/:imageId/like
func (i *InteractionRoutes) RegisterRoutes(router *gin.RouterGroup) {

	// group
	interactionImage := router.Group("interactions/images")

	// for like the image
	interactionImage.POST("/:imageId/like", i.authMiddleware, i.interactionController.LikeImage)

	// for unlike the image
	interactionImage.DELETE("/:imageId/like", i.authMiddleware, i.interactionController.UnlikeImage)

	// for comment the image
	interactionImage.POST("/:imageId/comment", i.authMiddleware, i.interactionController.AddComment)
}

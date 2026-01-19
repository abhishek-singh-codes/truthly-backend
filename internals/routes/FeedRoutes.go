package routes

import (
	"truthly/internals/controller"
	"truthly/internals/middleware"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
)

type FeedRoutes struct {
	feedController *controller.FeedController
	authMiddleware gin.HandlerFunc
}

func GetNewFeedRoutes(feedController *controller.FeedController, authToken *auth.AuthToken) *FeedRoutes {
	return &FeedRoutes{
		feedController: feedController,
		authMiddleware: middleware.AuthMiddleware(authToken),
	}
}

// api/v1/feed
func (fr *FeedRoutes) RegisterRoutes(router *gin.RouterGroup) {

	// versioned group
	feedGroup := router.Group("/feed")

	// get feed
	feedGroup.GET("", fr.authMiddleware, fr.feedController.GetFeed)

	//get Feed by range
	feedGroup.GET("/nearby", fr.authMiddleware, fr.feedController.GetFeedByRange)
}

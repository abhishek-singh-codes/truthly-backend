package routes

import (
	"truthly/internals/controller"
	"truthly/internals/middleware"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
)

type PostImageRoutes struct {
	postImageController *controller.PostImageController
	authMiddleware      gin.HandlerFunc
}

func GetNewPostImageRoutes(postImageController *controller.PostImageController, authToken *auth.AuthToken) *PostImageRoutes {
	return &PostImageRoutes{
		postImageController: postImageController,
		authMiddleware:      middleware.AuthMiddleware(authToken),
	}
}

// api/v1/posts
func (p *PostImageRoutes) RegisterRoutes(router *gin.RouterGroup) {

	// group
	postImageGroup := router.Group("/posts")

	// create a new post
	postImageGroup.POST("", p.authMiddleware, p.postImageController.PostImage)
}

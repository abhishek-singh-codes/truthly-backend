package routes

import (
	"truthly/internals/controller"
	"truthly/internals/middleware"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
)

type ProfileRoutes struct {
	profileController *controller.ProfileController
	authMiddleware    gin.HandlerFunc
}

// object creation
func GetNewProfileRoutes(
	pc *controller.ProfileController,
	authToken *auth.AuthToken,
) *ProfileRoutes {
	return &ProfileRoutes{
		profileController: pc,
		authMiddleware:    middleware.AuthMiddleware(authToken),
	}
}

// api/v1/profile
func (p *ProfileRoutes) RegisterRoutes(router *gin.RouterGroup) {

	// group
	ProfileGroup := router.Group("/profile")

	// create a new post
	ProfileGroup.GET("", p.authMiddleware, p.profileController.GetProfileDetails)
}

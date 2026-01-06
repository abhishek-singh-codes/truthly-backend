package routes

import (
	"log/slog"
	"truthly/internals/controller"
	"truthly/internals/realtime"
	"truthly/internals/repository"
	"truthly/internals/service"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAll(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger, hub *realtime.Hub) {
	registerPost(router, db, logger)
	registerFeed(router, db, logger)
	registerAuth(router, db, logger)
	registerInteraction(router, db, logger, hub)
	registerWebsocket(router, hub, db, logger)
	registerUser(router, db, logger)
}

// post
func registerPost(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {

	imageRepo := repository.GetImageRepo(db, logger)
	descriptionRepo := repository.GetDescriptionRepository(db, logger)
	analyticsRepo := repository.GetAnalyticRepository(db, logger)
	commentRepo := repository.GetCommentRepository(db, logger)
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)

	s3Uploader, err := service.NewS3Uploader("truthly-images", logger)
	if err != nil {
		logger.Error(err.Error())
	}

	postService := service.GetPostService(logger, analyticsRepo, commentRepo,
		descriptionRepo, imageRepo, s3Uploader,
	)
	postImageController := controller.GetNewPostImageController(logger, postService)

	authToken := auth.GetNewAuthToken(logger, userSessionRepo)

	GetNewPostImageRoutes(postImageController, authToken).RegisterRoutes(router)

}

// feed
// Get the images in feed section in pagination form
func registerFeed(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {

	// repo's
	// Get feed items [{}, {}]
	feedRepo := repository.GetNewFeedRepository(db, logger)

	feedService := service.GetNewFeedService(feedRepo, logger)
	feedController := controller.GetNewFeedController(logger, feedService)

	// protected routes
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)
	// for verify the token
	authToken := auth.GetNewAuthToken(logger, userSessionRepo)
	// register routes in /GET
	GetNewFeedRoutes(feedController, authToken).RegisterRoutes(router)
}

// auth
func registerAuth(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	//repos
	// to delete last ACTIVE session and add new session
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)
	// to add the new user in signup
	userRepo := repository.GetUserRepo(logger, db)

	// auth service
	// User signup, verify mail and add session
	authService := service.GetNewAuthService(logger, userSessionRepo, userRepo)

	// utils
	// generate jwt token
	authUtil := auth.GetNewAuthToken(logger, userSessionRepo)

	// auth controller
	// handle signup and login
	authController := controller.GetNewAuthController(logger, authService, authUtil)

	// routes
	// auth/signup
	// auth/login
	GetNewAuthRoutes(authController).RegisterRoutes(router)
}

// interaction
func registerInteraction(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger, hub *realtime.Hub) {
	// repo required
	interactionRepo := repository.GetNewInteractionRepository(db, logger)
	analyticRepo := repository.GetAnalyticRepository(db, logger)
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)

	// service required
	interactionService := service.GetNewInteractionService(logger, interactionRepo, analyticRepo, hub)

	// auth
	authToken := auth.GetNewAuthToken(logger, userSessionRepo)

	// controller
	interactionController := controller.GetNewInteractionController(logger, interactionService)

	// routes
	GetNewInteractionRoutes(interactionController, authToken).RegisterRoutes(router)
}

// websocket
// websocket
func registerWebsocket(
	router *gin.RouterGroup,
	hub *realtime.Hub,
	db *gorm.DB,
	logger *slog.Logger,
) {

	// repo
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)

	// auth token
	authToken := auth.GetNewAuthToken(logger, userSessionRepo)

	// route
	router.GET("/ws", func(c *gin.Context) {
		controller.ServeWS(
			hub,
			c.Writer,
			c.Request,
			authToken,
			logger,
		)
	})
}

// user
func registerUser(router *gin.RouterGroup, db *gorm.DB, logger *slog.Logger) {
	// repo
	userRepo := repository.GetUserRepo(logger, db)
	userSessionRepo := repository.GetNewUserSessionRepo(logger, db)

	// service
	userService := service.GetNewUserService(logger, userRepo)

	// token
	authToken := auth.GetNewAuthToken(logger, userSessionRepo)

	// controller
	userController := controller.GetNewUserController(userService)

	GetNewUserRoutes(userController, authToken).RegisterRoutes(router)
}

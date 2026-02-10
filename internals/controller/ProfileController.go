package controller

import (
	"log/slog"
	"truthly/internals/dto"
	"truthly/internals/service"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	logger         *slog.Logger
	profileService service.ProfileService
}

func GetNewProfileController(
	l *slog.Logger,
	ps service.ProfileService,
) *ProfileController {
	return &ProfileController{
		logger:         l,
		profileService: ps,
	}
}

// handler
func (pc *ProfileController) GetProfileDetails(ctx *gin.Context) {

	userId := ctx.GetString("userId")

	resp, err := pc.profileService.ImagesByUserId(ctx, userId)
	if err != nil {
		ctx.JSON(500, dto.ResponseDto[any]{
			Status: "ERROR",
			Error:  err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.ResponseDto[any]{
		Status:    "Success",
		Message:   "profile data for user",
		ResultObj: resp,
	})

}

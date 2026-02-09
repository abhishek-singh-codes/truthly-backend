package controller

import (
	"log/slog"
	"truthly/internals/dto"
	"truthly/internals/service"
	"truthly/internals/util/auth"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	logger      *slog.Logger
	authService service.AuthService
	authUtil    *auth.AuthToken
}

func GetNewAuthController(logger *slog.Logger, authService service.AuthService, authUtil *auth.AuthToken) *AuthController {
	return &AuthController{
		logger:      logger,
		authService: authService,
		authUtil:    authUtil,
	}
}

func (c *AuthController) UserSignup(ctx *gin.Context) {
	// variable to taking data from request
	var user dto.UserRequestDto

	// binding the data
	if err := ctx.ShouldBindJSON(&user); err != nil {
		c.logger.Error(err.Error())
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// call the service to signup the user
	resp, err := c.authService.UserSignup(ctx, &user)
	if err != nil {
		c.logger.Error(err.Error())
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, resp)
}

func (c *AuthController) UserLogin(ctx *gin.Context) {
	// get the data from request
	var loginReq dto.LoginReq
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		c.logger.Error(err.Error())
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := c.authService.VerifyUser(ctx, &loginReq)
	if err != nil {
		ctx.JSON(400, gin.H{
			"Messsage": "Invalid Cardentials",
		})
		return
	}

	// set userId and userName in ctx to access it later
	ctx.Set("userId", data.ResultObj.UserId)
	ctx.Set("userName", loginReq.UserName)

	// generate token
	token, sessionId, err := c.authUtil.GenerateJwtToken(loginReq.UserName, data.ResultObj.UserId)
	if err != nil {
		c.logger.Error("Error in token generation", "error", err.Error())
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	// SET JWT IN COOKIE
	ctx.SetCookie(
		"access_token",
		token,
		86400,
		"/",
		"",
		false,
		true,
	)

	// Add this session in UserSession table
	res, err := c.authService.AddSession(ctx, sessionId, data.ResultObj.UserId, loginReq.UserName, token)
	if err != nil {
		ctx.JSON(500, res)
	}

	ctx.JSON(200, res)
}

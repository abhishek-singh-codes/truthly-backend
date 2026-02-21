package service

import (
	"context"
	"log/slog"
	"time"
	"truthly/internals/dto"
	"truthly/internals/model"
	"truthly/internals/repository"
)

type AuthService interface {
	UserSignup(ctx context.Context, user *dto.UserRequestDto) (*dto.ResponseDto[*dto.LogInRes], error)
	VerifyUser(ctx context.Context, loginReq *dto.LoginReq) (*dto.ResponseDto[*dto.LogInRes], error)
	AddSession(ctx context.Context, sessionId string, userId string, userName string, token string) (*dto.ResponseDto[*dto.LogInRes], error)
}

type authService struct {
	logger          *slog.Logger
	userSessionRepo repository.UserSessionRepository
	userRepo        repository.UserRepository
}

func GetNewAuthService(logger *slog.Logger, userSessionRepo repository.UserSessionRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		logger:          logger,
		userSessionRepo: userSessionRepo,
		userRepo:        userRepo,
	}
}

func (s *authService) UserSignup(ctx context.Context, userReq *dto.UserRequestDto) (*dto.ResponseDto[*dto.LogInRes], error) {
	//1. dto -> model
	user := dto.ToModel(userReq)

	//2. creating a new user
	savedUser, err := s.userRepo.CreatNewUser(ctx, user)
	if err != nil {
		s.logger.Error("Error while creating a new user, Error: " + err.Error())
		return nil, err
	}

	s.logger.Info("New user created", "userId", savedUser.UserId)

	//3. model -> dto
	return &dto.ResponseDto[*dto.LogInRes]{
		Status:  "success",
		Message: "User created",
		ResultObj: &dto.LogInRes{
			UserId: savedUser.UserId,
		},
	}, nil
}

// user name and password login
func (s *authService) VerifyUser(ctx context.Context, logInReq *dto.LoginReq) (*dto.ResponseDto[*dto.LogInRes], error) {

	res, err := s.userRepo.VerifyUser(ctx, logInReq.UserName, logInReq.Password)
	if err != nil {
		return &dto.ResponseDto[*dto.LogInRes]{
			Error: err.Error(),
		}, err
	}

	return &dto.ResponseDto[*dto.LogInRes]{
		Status:  "Success",
		Message: "User Existed",
		ResultObj: &dto.LogInRes{
			UserId: res,
		},
	}, nil
}

// Add session
func (s *authService) AddSession(ctx context.Context, sessionId string, userId string, userName string, token string) (*dto.ResponseDto[*dto.LogInRes], error) {
	// data ---> model
	userSession := model.UserSession{
		UserId:    userId,
		SessionId: sessionId,
		UserName:  userName,
		Status:    "ACTIVE",
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}
	// old session Id ko expired mark kardo
	err := s.userSessionRepo.ExpireLastActiveSession(ctx, userId)
	if err != nil {
		return &dto.ResponseDto[*dto.LogInRes]{
			Error: err.Error(),
		}, nil
	}
	// repo calling
	err = s.userSessionRepo.CreateNewSession(ctx, &userSession)
	if err != nil {
		return &dto.ResponseDto[*dto.LogInRes]{
			Status: "Error",
			Error:  err.Error(),
		}, nil
	}

	// return
	return &dto.ResponseDto[*dto.LogInRes]{
		Status:  "Success",
		Message: "Log in success",
		ResultObj: &dto.LogInRes{
			Token: token,
		},
	}, nil
}

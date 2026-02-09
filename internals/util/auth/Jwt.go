package auth

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"
	"truthly/internals/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	UserName  string `json:"userName"`
	UserId    string `json:"userId"`
	SessionId string `json:"sessionId"`

	jwt.RegisteredClaims
}

type AuthToken struct {
	logger          *slog.Logger
	secret          []byte
	userSessionRepo repository.UserSessionRepository
}

func GetNewAuthToken(logger *slog.Logger, userSessionRepo repository.UserSessionRepository) *AuthToken {
	// get secret from .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}

	return &AuthToken{
		logger:          logger,
		secret:          []byte(secret),
		userSessionRepo: userSessionRepo,
	}
}

// Generate jwt token
func (at *AuthToken) GenerateJwtToken(userName string, userId string) (string, string, error) {
	sessionId := uuid.New().String()
	if sessionId == "" {
		at.logger.Error("faild to generate the session Id")
		return "", "", errors.New("faild to generate the session Id")
	}
	claims := JwtClaims{
		UserName:  userName,
		UserId:    userId,
		SessionId: sessionId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(at.secret)
	if err != nil {
		at.logger.Error("failed to sign jwt token", "error", err)
		return "", "", err
	}

	return signedToken, sessionId, nil
}

// verify token
func (at *AuthToken) VerifyJwtToken(tokenString string, ctx context.Context) (*JwtClaims, error) {
	// get claims and verify

	token, err := jwt.ParseWithClaims(tokenString,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// signing method verify
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				at.logger.Error(
					"Unexpected signing method",
					"alg", token.Header["alg"],
				)
				return nil, errors.New("Unexpected signing method")
			}
			return at.secret, nil
		},
	)

	if err != nil {
		at.logger.Error("Invalid jwt token", "error", err.Error())
		return nil, err
	}

	// extract claims
	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		at.logger.Error("Invalid JWT Token")
		return nil, errors.New("Invalid JWT Token")
	}

	// check jwt authentication time
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		at.logger.Warn("jwt token expired", "userId", claims.UserId)
		return nil, errors.New("token expired")
	}

	// based on SessionId if Status 'ACTIVE' Valid , 'EXPIRED'
	_, err = at.userSessionRepo.GetActiveSession(ctx, claims.SessionId)
	if err != nil {
		return nil, errors.New("Session not found or expired")
	}

	return claims, nil
}

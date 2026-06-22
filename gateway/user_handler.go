package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hubit-space/gateway/config"
	"hubit-space/gateway/model"
	"hubit-space/gateway/repository"
	"hubit-space/gateway/utility"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthHandler struct {
	repo        repository.UserRepository
	redisClient *redis.Client
}

func NewAuthHandler(repo repository.UserRepository, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var requestBody map[string]any
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	appsToken, appsTokenExists := requestBody["apps_token"].(string)
	if !appsTokenExists || appsToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Apps token is required"})
		return
	}

	payload, err := h.VerifyAppsToken(appsToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	email, emailExists := payload["email"].(string)
	if !emailExists || email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
	}

	// Generate tokens
	response, err := h.GenerateAllToken(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	var requestBody map[string]any
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, refreshTokenExists := requestBody["refresh_token"].(string)
	if !refreshTokenExists || refreshToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	email, ok := utility.GetHeader(ctx, "Email")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email header is required"})
		return
	}

	// Validate refresh token
	key := fmt.Sprintf("scb-refresh-token:%s", email)
	existingToken, err := h.redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch refresh token: " + err.Error()})
		return
	}

	if refreshToken != existingToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Generate tokens
	response, err := h.GenerateAllToken(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GenerateAllToken(email string) (model.AuthResponse, error) {
	// Validate user credentials
	user, err := h.repo.GetIdentity(email)
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Generate access token
	accessToken, err := h.GenerateAccessToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Generate refresh token
	refreshToken, err := h.GenerateRefreshToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	response := model.AuthResponse{
		UserID:   user.ID.Hex(),
		UserCode: user.UserCode,
		Name:     user.Name,
		RoleCode: user.Role.RoleCode,

		Email:        user.Email,
		Menus:        user.Menus,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (h *AuthHandler) VerifyAppsToken(appsToken string) (jwt.MapClaims, error) {
	jwtSecretApps := []byte(config.GetEnv("JWT_SECRET_KEY_APPS", ""))

	token, err := jwt.Parse(appsToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecretApps, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid or expired token")
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid or expired token")
	}

	return payload, nil
}

func (h *AuthHandler) GenerateAccessToken(user model.User) (string, error) {
	tokenExpire, err := strconv.Atoi(config.GetEnv("TOKEN_EXPIRE", "3600"))
	if err != nil {
		fmt.Println("Error:", err)
	}

	claims := model.AccessToken{
		UserId:   user.ID.Hex(),
		UserCode: user.UserCode,
		Name:     user.Name,
		RoleCode: user.Role.RoleCode,
		Email:    user.Email,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenExpire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := []byte(config.GetEnv("JWT_SECRET_KEY", ""))

	return token.SignedString(jwtSecret)
}

func (h *AuthHandler) GenerateRefreshToken(user model.User) (string, error) {
	token := ""
	key := fmt.Sprintf("scb-refresh-token:%s", user.Email)
	existingToken, err := h.redisClient.Get(context.Background(), key).Result()

	if err == redis.Nil {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}

		timeBytes := []byte(time.Now().String())
		token = base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(append(bytes, timeBytes...))

	} else {
		token = existingToken
	}

	refreshExpire := config.GetEnv("REFRESH_TOKEN_EXPIRE", "86400")
	expire, err := strconv.Atoi(refreshExpire)
	if err != nil {
		fmt.Println("Error:", err)
	}

	expiresAt := time.Duration(expire) * time.Second

	err = h.redisClient.Set(context.Background(), key, token, expiresAt).Err()
	if err != nil {
		return "", err
	}

	return token, nil
}

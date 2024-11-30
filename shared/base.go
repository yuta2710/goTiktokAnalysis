package shared

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
}

type RouterGroup struct {
	Protected *echo.Group `json:"protected"`
	Public    *echo.Group `json:"public"`
}

func Response(c echo.Context, success bool, statusCode int, message string, data interface{}, headers map[string]string) error {
	if headers != nil {
		MapHeaders(c, headers)
	}

	if data != nil {
		return c.JSON(statusCode, BaseResponse{Success: success, Message: message, Status: statusCode, Data: data})
	}

	return c.JSON(statusCode, BaseResponse{Success: success, Message: message, Status: statusCode})
}

func MapHeaders(ctx echo.Context, headers map[string]string) error {
	for key, val := range headers {
		ctx.Response().Header().Set(key, val)
	}

	return nil
}

func TokenProvider(userId int, authId string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"authId": authId,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})

	accessSecret := os.Getenv("ACCESS_SECRET")
	accessTokenString, err := accessToken.SignedString([]byte(accessSecret))

	if err != nil {
		return "", "", fmt.Errorf("Login failed, something wrong due to processing access token to String type")
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"authId": authId,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	refreshSecret := os.Getenv("REFRESH_SECRET")
	refreshTokenString, err := refreshToken.SignedString([]byte(refreshSecret))

	if err != nil {
		return "", "", fmt.Errorf("Login failed, something wrong due to processing refresh token to String type")
	}

	return accessTokenString, refreshTokenString, nil
}

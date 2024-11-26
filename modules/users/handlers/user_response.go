package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type UserResponse struct {
	shared.BaseResponse
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AuthId       string `json:"auth_id"`
}

func response(ctx echo.Context, responseCode int, message string) error {
	return ctx.JSON(responseCode, &shared.BaseResponse{
		Message: message,
	})
}

package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/auth/models"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/auth/usecases"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type AuthHttp struct {
	AuthUsecase usecases.AuthUsecase
}

func (ahp *AuthHttp) Login(ctx echo.Context) error {
	// Check user exist in the repo
	body := new(models.LoginRequest)

	if err := ctx.Bind(body); err != nil {
		return shared.Response(ctx, false, http.StatusBadRequest, "Invalid login request body", nil, nil)

	}

	authResp, err := ahp.AuthUsecase.Login(body)

	if err != nil {
		return shared.Response(ctx, false, http.StatusBadRequest, err.Error(), nil, nil)
	}

	authorization := fmt.Sprintf("Bearer %s", authResp.AccessToken)

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": authorization,
		"Custom-Header": "CustomValue",
	}

	shared.MapHeaders(ctx, headers)

	refreshTokenCookie := &http.Cookie{
		Name:     "REFRESH_TOKEN",
		Value:    authResp.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
	}

	ctx.SetCookie(refreshTokenCookie)

	return shared.Response(
		ctx,
		true,
		http.StatusOK,
		"Login successfully",
		&models.AuthResponse{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
		},
		headers,
	)
}

func (ahp *AuthHttp) SignUp(ctx echo.Context) error {
	body := new(models.SignUpRequest)

	if err := ctx.Bind(body); err != nil {
		return err
	}

	// get repo
	authResp, err := ahp.AuthUsecase.SignUp(body)

	if err != nil {
		return err
	}

	authorization := fmt.Sprintf("Bearer %s", authResp.AccessToken)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": authorization,
		"Custom-Header": "CustomValue",
	}
	shared.MapHeaders(ctx, headers)
	refreshTokenCookie := &http.Cookie{
		Name:     "REFRESH_TOKEN",
		Value:    authResp.RefreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
	}

	ctx.SetCookie(refreshTokenCookie)

	return shared.Response(
		ctx,
		true,
		http.StatusOK,
		"Sign Up successfully",
		&models.AuthResponse{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
		},
		headers,
	)
}
func (ahp *AuthHttp) Profile(ctx echo.Context) error {
	user := ctx.Get("user").(*entities.FetchUserDto)

	fmt.Println(user)

	if user == nil {
		return shared.Response(ctx, false, http.StatusNotFound, "Profile not found", nil, nil)
	}

	return shared.Response(ctx, true, http.StatusOK, "Profile is here", user, nil)
}

func (ahp *AuthHttp) SignOut(ctx echo.Context) error {
	return nil
}

func NewAuthHttp(auc usecases.AuthUsecase) AuthHandler {
	return &AuthHttp{
		AuthUsecase: auc,
	}
}

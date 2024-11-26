package handlers

import "github.com/labstack/echo/v4"

type AuthHandler interface {
	Login(ctx echo.Context) error
	SignUp(ctx echo.Context) error
	Profile(ctx echo.Context) error
	SignOut(ctx echo.Context) error
}

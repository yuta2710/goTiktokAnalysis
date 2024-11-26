package handlers

import "github.com/labstack/echo/v4"

type UserHandler interface {
	CreateNewUser(ctx echo.Context) error
	GetUserById(ctx echo.Context) error
	GetUsers(ctx echo.Context) error
}

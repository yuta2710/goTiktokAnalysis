package handlers

import (
	"github.com/labstack/echo/v4"
)

type TodoHandler interface {
	CreateNewTodo(ctx echo.Context) error
	InsertBatch(ctx echo.Context) error
	FindById(ctx echo.Context) error
	FindAllByUserId(ctx echo.Context) error
	UpdateTodo(ctx echo.Context) error
	UpdateAvatarOfTodo(ctx echo.Context) error
	DeleteTodo(ctx echo.Context) error
}

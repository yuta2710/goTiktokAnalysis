package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/models"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/usecases"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type TodoHttp struct {
	Uc usecases.TodoUsecase
}

func (tdHttp *TodoHttp) CreateNewTodo(ctx echo.Context) error {
	body := new(models.InsertTodoSample)

	if err := ctx.Bind(body); err != nil {
		return err
	}

	stdCxt := ctx.Request().Context()
	todoId, err := tdHttp.Uc.Insert(stdCxt, body)

	if err != nil {
		return shared.Response(ctx, false, http.StatusBadRequest, "Error inserting todo", nil, nil)
	}

	return shared.Response(ctx, true, http.StatusOK, "Inserted todo successfully", todoId, nil)
}

func (tdHttp *TodoHttp) InsertBatch(ctx echo.Context) error {
	return nil
}

func (tdHttp *TodoHttp) FindById(ctx echo.Context) error {
	id := ctx.Param("id")

	fmt.Println(id)
	stdCtx := ctx.Request().Context()

	todo, err := tdHttp.Uc.FindById(stdCtx, id)

	if err != nil {
		return shared.Response(ctx, false, http.StatusBadRequest, "Todo not found", nil, nil)
	}
	return shared.Response(ctx, true, http.StatusOK, "Get todo successfully", todo, nil)
}

func (tdHttp *TodoHttp) FindAllByUserId(ctx echo.Context) error {
	id := ctx.Param("authId")
	stdCtx := ctx.Request().Context()
	todos, err := tdHttp.Uc.FindAllByUserId(stdCtx, id)

	if err != nil {
		return shared.Response(ctx, false, http.StatusBadRequest, "Todo not found or empty", nil, nil)
	}

	return shared.Response(ctx, true, http.StatusOK, "Get todos successfully", todos, nil)
}

func (tdHttp *TodoHttp) UpdateTodo(ctx echo.Context) error {
	return nil
}

func (tdHttp *TodoHttp) UpdateAvatarOfTodo(ctx echo.Context) error {
	return nil
}

func (tdHttp *TodoHttp) DeleteTodo(ctx echo.Context) error {
	return nil
}

func NewTodoHttp(uc usecases.TodoUsecase) TodoHandler {
	return &TodoHttp{
		Uc: uc,
	}
}

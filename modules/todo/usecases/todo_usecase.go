package usecases

import (
	"context"

	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/models"
)

type TodoUsecase interface {
	Insert(ctx context.Context, in *models.InsertTodoSample) (int, error)
	InsertBatch(ctx context.Context, in []*models.InsertTodoSample) error
	FindById(ctx context.Context, id string) (*models.FetchTodoDto, error)
	// FindAllByUserId(ctx context.Context, fakeId string) ([]*entities.Todo, error)
	FindAllByUserId(ctx context.Context, fakeId string) ([]*models.FetchTodoDto, error)
	UpdateTodo(ctx context.Context, id string, sample *models.UpdateTodoSample) error
	UpdateAvatarOfTodo(ctx context.Context, id string, sample *models.UpdateTodoAvatarSample) error
	DeleteTodo(ctx context.Context, id string) error
}

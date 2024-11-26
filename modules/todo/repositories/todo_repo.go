package repositories

import (
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/models"
)

type TodoRepository interface {
	// Insert(in *entities.Todo) (int, error)
	InsertTodo(in *entities.Todo) (int, error)
	InsertTodoMembers(todoId int, members []entities.TodoMember) error

	InsertBatch(in []*entities.Todo) error
	FindById(id int) (*entities.Todo, error)
	FindAllByUserId(userId int) ([]*entities.Todo, error)

	UpdateTodo(id int, sample *models.UpdateTodoSample) error

	UpdateAvatarOfTodo(id string, sample *models.UpdateTodoAvatarSample) error
	DeleteTodo(id string) error
	AddUserForTodo(userId string) error
}

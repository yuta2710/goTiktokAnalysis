package usecases

import (
	"context"
	"fmt"

	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/models"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/repositories"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type TodoUsecaseImpl struct {
	Repo repositories.TodoRepository
	// AuthIdProvider shared.AuthIdProvider
}

func (tduc *TodoUsecaseImpl) Insert(ctx context.Context, in *models.InsertTodoSample) (int, error) {
	authId, ok := ctx.Value("authId").(string)

	if !ok {
		return 0, fmt.Errorf("authId not found in context")
	}

	fmt.Println("Clm no hahaha %s", authId)
	// if tduc.AuthIdProvider == nil {
	// 	return 0, fmt.Errorf("AuthIdProvider is not initialized")
	// } else {
	// 	fmt.Println("Let's go")
	// }

	decodedId, _ := shared.DecomposeUidV2(string(authId))
	fmt.Println(in.Members)

	// in.Priority = priority

	for i, _ := range in.Members {
		in.Members[i].TodoId = 1
	}

	toEntity := &entities.Todo{
		// AuthId:      authId,
		UserId:      int(decodedId.GetLocalID()),
		Title:       in.Title,
		Description: in.Description,
		IsCompleted: in.IsCompleted,
		DueDate:     in.DueDate,
		Priority:    in.Priority,
		// Members:     in.Members,
	}

	fmt.Println("Todo entity data inserted")

	newTodoId, err := tduc.Repo.InsertTodo(toEntity)

	if err != nil || newTodoId == 0 {
		return 0, fmt.Errorf("[DB Layer]: failed to insert todo members: %v", err)
	}

	members := []entities.TodoMember{}

	for _, mem := range in.Members {
		// decodedId, _ := shared.DecomposeUidV2(string(authId))

		members = append(members, entities.TodoMember{
			UserId: mem.UserId,
			TodoId: newTodoId,
			Role:   mem.Role,
		})
	}

	if len(members) > 0 {
		err = tduc.Repo.InsertTodoMembers(newTodoId, members)

		if err != nil {
			return 0, fmt.Errorf("failed to insert todo members: %v", err)
		}
	}

	// Update members to database

	return newTodoId, nil
}

func (tduc *TodoUsecaseImpl) InsertBatch(ctx context.Context, in []*models.InsertTodoSample) error {
	return nil
}

func (tduc *TodoUsecaseImpl) FindById(ctx context.Context, id string) (*models.FetchTodoDto, error) {
	// id := ctx.Value()
	fmt.Println("Alo alo alo")
	b64ap := &shared.Base64AuthIdProvider{}
	decodedTodoId, err := b64ap.Decode(id, "todo")

	if err != nil {
		return nil, fmt.Errorf("Todo not found due to fail-decoded id", err)
	}

	todo, err := tduc.Repo.FindById(decodedTodoId)

	if err != nil {
		return nil, fmt.Errorf("Todo not found", err)
	}

	fmt.Println(todo.Id)
	// var dto *models.FetchTodoDto

	fetchedTodo := &models.FetchTodoDto{}

	if todo != nil {
		fetchedTodo.FakeId = b64ap.Encode(todo.Id, "todo")
		fetchedTodo.UserId = b64ap.Encode(todo.UserId, "user")
		fetchedTodo.Title = todo.Title
		fetchedTodo.Description = todo.Description
		fetchedTodo.IsCompleted = todo.IsCompleted
		fetchedTodo.DueDate = todo.DueDate
		fetchedTodo.Priority = todo.Priority
		fetchedTodo.CreatedAt = todo.CreatedAt
		fetchedTodo.UpdatedAt = todo.UpdatedAt
		fetchedTodo.Members = todo.Members
	}

	return fetchedTodo, nil
}

func (tduc *TodoUsecaseImpl) FindAllByUserId(ctx context.Context, fakeId string) ([]*models.FetchTodoDto, error) {
	componentId, _ := shared.DecomposeUidV2(fakeId)
	localId := componentId.GetLocalID()

	fmt.Printf("Decoded ID: %d\n", localId)

	// Fetch todos from the repository
	preprocessedTodos, err := tduc.Repo.FindAllByUserId(int(localId))

	if err != nil {
		return nil, err
	}

	// fake := &models.Fake{}

	// Prepare the fetchedTodos slice
	fetchedTodos := make([]*models.FetchTodoDto, len(preprocessedTodos))

	for i, todo := range preprocessedTodos {
		if todo == nil {
			fmt.Printf("Skipped nil todo at index %d\n", i)
			continue
		}

		// Map todo to FetchTodoDto
		// fetchedTodos[i] = models.NewFetchTodoDto(todo)
		if todo != nil {
			b64ap := &shared.Base64AuthIdProvider{}
			fetchedTodos[i] = &models.FetchTodoDto{}

			fmt.Println(todo.Id)
			// var dto *models.FetchTodoDto
			fetchedTodos[i].FakeId = b64ap.Encode(todo.Id, "todo")
			fetchedTodos[i].UserId = b64ap.Encode(todo.UserId, "user")
			fetchedTodos[i].Title = todo.Title
			fetchedTodos[i].Description = todo.Description
			fetchedTodos[i].IsCompleted = todo.IsCompleted
			fetchedTodos[i].DueDate = todo.DueDate
			fetchedTodos[i].Priority = todo.Priority
			fetchedTodos[i].CreatedAt = todo.CreatedAt
			fetchedTodos[i].UpdatedAt = todo.UpdatedAt
			fetchedTodos[i].Members = todo.Members

			// Debug log to ensure mapping is successful
			if fetchedTodos[i] == nil {
				fmt.Printf("Failed to map todo[%d] to fetchedTodoDto\n", i)
			} else {
				fmt.Printf("Mapped todo[%d] to fetchedTodo[%d]: %+v\n", i, fetchedTodos[i])
			}
		}
		// dto.Title = todo.Title
	}

	fmt.Println("Mapped fetchedTodos: ", fetchedTodos)
	return fetchedTodos, nil
}

func (tduc *TodoUsecaseImpl) UpdateTodo(ctx context.Context, id string, sample *models.UpdateTodoSample) error {
	return nil
}

func (tduc *TodoUsecaseImpl) UpdateAvatarOfTodo(ctx context.Context, id string, sample *models.UpdateTodoAvatarSample) error {
	return nil
}

func (tduc *TodoUsecaseImpl) DeleteTodo(ctx context.Context, id string) error {
	return nil
}

func NewTodoUsecaseImpl(repo repositories.TodoRepository) TodoUsecase {
	return &TodoUsecaseImpl{
		Repo: repo,
		// AuthIdProvider: authIdProvider,
	}
}

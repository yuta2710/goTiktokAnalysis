package models

import (
	"time"

	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
)

type (
	InsertTodoSample struct {
		AuthId      string                `json:"authId"`
		Title       string                `json:"title"`
		Description string                `json:"description"`
		IsCompleted bool                  `json:"isCompleted"`
		DueDate     time.Time             `json:"dueDate"`
		Priority    entities.Priority     `json:"priority"`
		Members     []entities.TodoMember `json:"members"`
	}

	InsertTodoMemberDto struct {
		UserId int                 `json:"userId"`
		TodoId int                 `json:"todoId"`
		Role   entities.MemberRole `json:"role"`
	}

	UpdateTodoSample struct {
		Title       *string            `json:"title,omitempty"`
		Description *string            `json:"description,omitempty"`
		IsCompleted *bool              `json:"isCompleted,omitempty"`
		DueDate     *time.Time         `json:"dueDate,omitempty"`
		Priority    *entities.Priority `json:"priority,omitempty"`
	}

	UpdateTodoAvatarSample struct {
	}
)

type Fake struct {
	FakeId string `json:"fakeId"`
}

type FetchTodoDto struct {
	FakeId      string                `json:"fakeId,omitempty"`
	UserId      string                `json:"userId"`      // Matches the user_id foreign key in the database
	Title       string                `json:"title"`       // Matches the title column in the database
	Description string                `json:"description"` // Matches the description column in the database
	IsCompleted bool                  `json:"isCompleted"` // Matches the is_completed column in the database
	DueDate     time.Time             `json:"dueDate"`     // Matches the due_date column in the database
	Priority    entities.Priority     `json:"priority"`    // Matches the priority column as ENUM
	CreatedAt   time.Time             `json:"createdAt"`   // Tracks the creation time
	UpdatedAt   time.Time             `json:"updatedAt"`   // Tracks the update time
	Members     []entities.TodoMember `json:"members"`     // For TodoMember association
}

// func NewFetchTodoDto(todo *entities.Todo) *FetchTodoDto {
// 	// Ensure that todo is not nil before proceeding
// 	if todo == nil {
// 		fmt.Println("NewFetchTodoDto received a nil todo")
// 		return nil
// 	}

// 	// Debug log for incoming todo
// 	fmt.Printf("Mapping todo to FetchTodoDto: %+v\n", todo)

// 	// Map fields from todo to FetchTodoDto
// 	return &FetchTodoDto{
// 		FakeId:      todo.FakeId.String(), // Handle possible nil FakeId
// 		UserId:      (todo.UserId),        // Map UserId
// 		Title:       (todo.Title),         // Map Title
// 		Description: (todo.Description),   // Map Description
// 		IsCompleted: (todo.IsCompleted),   // Map IsCompleted
// 		DueDate:     (todo.DueDate),       // Map DueDate
// 		Priority:    (todo.Priority),      // Map Priority
// 		CreatedAt:   (todo.CreatedAt),     // Map CreatedAt
// 		UpdatedAt:   (todo.UpdatedAt),     // Map UpdatedAt
// 		Members:     (todo.Members),       // Map Members
// 	}
// }

// // Helper functions for handling nils
// func stringOrNil(value string) *string {
// 	if value == "" {
// 		return nil
// 	}
// 	return &value
// }

// func intOrNil(value int) *int {
// 	if value == 0 {
// 		return nil
// 	}
// 	return &value
// }

// func boolOrNil(value bool) *bool {
// 	return &value
// }

// func timeOrNil(value time.Time) *time.Time {
// 	if value.IsZero() {
// 		return nil
// 	}
// 	return &value
// }

// func priorityOrNil(value entities.Priority) *entities.Priority {
// 	return &value
// }

// // Map members to pointers
// func mapTodoMembers(members []entities.TodoMember) []*entities.TodoMember {
// 	if members == nil {
// 		return nil
// 	}
// 	mapped := make([]*entities.TodoMember, len(members))
// 	for i, member := range members {
// 		mapped[i] = &member
// 	}
// 	return mapped
// }

// result := tdr.db.GetDb().
// Preload("Members").
// Where("todos.user_id = ?", userId).
// Find(&preprocessedTodos)

// if result.Error != nil {
// return nil, result.Error
// }

// for i, todo := range preprocessedTodos {
// fetchedTodos[i].FakeId = todo.FakeId.String()
// fetchedTodos[i].UserId = todo.UserId
// fetchedTodos[i].Title = todo.Title
// fetchedTodos[i].Description = todo.Description
// fetchedTodos[i].IsCompleted = todo.IsCompleted
// fetchedTodos[i].DueDate = todo.DueDate
// fetchedTodos[i].Priority = todo.Priority
// fetchedTodos[i].CreatedAt = todo.CreatedAt
// fetchedTodos[i].UpdatedAt = todo.UpdatedAt
// fetchedTodos[i].Members = todo.Members
// }

// func NewFetchTodoDto(todo *entities.Todo) *FetchTodoDto {
// 	if todo == nil {
// 		return nil
// 	}

//		return &FetchTodoDto{
//			FakeId:      todo.FakeId.String(),
//			UserId:      todo.UserId,
//			Title:       todo.Title,
//			Description: todo.Description,
//			IsCompleted: todo.IsCompleted,
//			DueDate:     todo.DueDate,
//			Priority:    todo.Priority,
//			CreatedAt:   todo.CreatedAt,
//			UpdatedAt:   todo.UpdatedAt,
//			Members:     todo.Members,
//		}
//	}

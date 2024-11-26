package repositories

import (
	"fmt"

	"github.com/yuta_2710/go-clean-arc-reviews/database"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/todo/models"
)

type TodoPostgresRepository struct {
	db database.Database
}

// func (tdr *TodoPostgresRepository) Insert(in *entities.Todo) (int, error) {
// 	// Transaction đầu tiên: Lưu vào bảng todos
// 	txTodo := tdr.db.GetDb().Begin()

// 	// Thêm bản ghi vào todos
// 	if err := txTodo.Create(in).Error; err != nil {
// 		txTodo.Rollback()
// 		return 0, fmt.Errorf("failed to insert todo: %v", err)
// 	}

// 	// Lấy ID của bản ghi todo vừa được tạo
// 	if in.Id == 0 {
// 		txTodo.Rollback()
// 		return 0, fmt.Errorf("failed to fetch inserted todo ID")
// 	}

// 	// Commit transaction đầu tiên
// 	if err := txTodo.Commit().Error; err != nil {
// 		return 0, fmt.Errorf("failed to commit todo transaction: %v", err)
// 	}

// 	// Transaction thứ hai: Lưu vào bảng todo_members
// 	if len(in.Members) > 0 {
// 		txMembers := tdr.db.GetDb().Begin()

// 		for _, member := range in.Members {
// 			member.TodoId = in.Id
// 			if err := txMembers.Create(&member).Error; err != nil {
// 				txMembers.Rollback()
// 				return 0, fmt.Errorf("failed to insert todo member: %v", err)
// 			}
// 		}

// 		// Commit transaction thứ hai
// 		if err := txMembers.Commit().Error; err != nil {
// 			return 0, fmt.Errorf("failed to commit todo members transaction: %v", err)
// 		}
// 	}

// 	fmt.Println("[INSERTED DATA SUCCESSFULLY]")
// 	return in.Id, nil
// }

func (tdr *TodoPostgresRepository) InsertTodo(in *entities.Todo) (int, error) {
	tx := tdr.db.GetDb().Begin()

	fmt.Println(in)

	// Lưu todo
	if err := tx.Create(in).Error; err != nil {
		fmt.Println("Yoooo bro oi")
		tx.Rollback()
		return 0, fmt.Errorf("failed to insert todo: %v", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return 0, fmt.Errorf("failed to commit todo transaction: %v", err)
	}

	fmt.Println("[TODO INSERTED SUCCESSFULLY]")
	return in.Id, nil
}

func (tdr *TodoPostgresRepository) InsertTodoMembers(todoId int, members []entities.TodoMember) error {
	tx := tdr.db.GetDb().Begin()

	// Lưu từng member vào todo_members
	for _, member := range members {
		member.TodoId = todoId
		if err := tx.Create(&member).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert todo member: %v", err)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit todo members transaction: %v", err)
	}

	fmt.Println("[TODO MEMBERS INSERTED SUCCESSFULLY]")
	return nil
}

func (tdr *TodoPostgresRepository) InsertBatch(in []*entities.Todo) error {
	return nil
}

func (tdr *TodoPostgresRepository) FindById(id int) (*entities.Todo, error) {
	var todo *entities.Todo

	result := tdr.db.GetDb().
		Preload("Members").
		Where("id = ?", id).
		First(&todo)

	if result.Error != nil {
		return nil, result.Error
	}

	return todo, nil
}

func (tdr *TodoPostgresRepository) FindAllByUserId(userId int) ([]*entities.Todo, error) {
	var todos []*entities.Todo

	result := tdr.db.GetDb().
		Preload("Members").
		Where("todos.user_id = ?", userId).
		Find(&todos)

	if result.Error != nil {
		return nil, result.Error
	}

	return todos, nil
}

func (tdr *TodoPostgresRepository) UpdateTodo(id int, sample *models.UpdateTodoSample) error {
	sampleMap := map[string]interface{}{}

	if sample.Title != nil {
		sampleMap["title"] = *sample.Title
	}

	if sample.Description != nil {
		sampleMap["description"] = *sample.Description
	}

	if sample.IsCompleted != nil {
		sampleMap["is_completed"] = *sample.IsCompleted
	}

	if sample.DueDate != nil {
		sampleMap["due_date"] = *sample.DueDate
	}

	if sample.Priority != nil {
		sampleMap["priority"] = *sample.Priority
	}

	if len(sampleMap) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := tdr.db.GetDb().Model(&entities.Todo{}).Where("id = ?", id).Updates(sampleMap)

	if result.Error != nil {
		return fmt.Errorf("failed to update todo: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("Todo not found to update")
	}

	return nil
}

func (tdr *TodoPostgresRepository) UpdateAvatarOfTodo(id string, sample *models.UpdateTodoAvatarSample) error {
	return nil
}

func (tdr *TodoPostgresRepository) DeleteTodo(id string) error {
	return nil
}

func (tdr *TodoPostgresRepository) AddUserForTodo(userId string) error {
	return nil
}

func NewTodoPostgresRepository(db database.Database) TodoRepository {
	return &TodoPostgresRepository{
		db: db,
	}
}

func convertPriorityToEnum(priority entities.Priority) string {
	switch priority {
	case entities.Low:
		return "Low"
	case entities.Medium:
		return "Medium"
	case entities.High:
		return "High"
	default:
		return "Unknown" // Or handle invalid case appropriately
	}
}

package repositories

import "github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"

type UserRepository interface {
	Insert(in *entities.InsertUserDto) (int, string, error)
	InsertBatch(in []*entities.InsertUserDto) error
	FindAll() ([]*entities.User, error)
	FindById(id int) (*entities.User, error)
	FindByEmail(email string) (*entities.User, error)
}

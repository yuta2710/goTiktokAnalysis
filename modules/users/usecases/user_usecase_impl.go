package usecases

import (
	"fmt"

	// TodoEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/models"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/repositories"
	UserRepo "github.com/yuta_2710/go-clean-arc-reviews/modules/users/repositories"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type UserUsecaseImpl struct {
	userRepo UserRepo.UserRepository
}

func (uui *UserUsecaseImpl) InsertNewUser(mod *models.InsertUserRequest) (int, string, error) {
	insertDto := entities.NewInsertUserRequest(mod.FirstName, mod.LastName, mod.Email, mod.Password, mod.Role)
	userId, fakeId, err := uui.userRepo.Insert(insertDto)

	if err != nil {
		return 0, "", err
	}

	fmt.Println("[CREATE ACCOUNT FROM USECASE LAYER SUCCESSFULLY]")

	return userId, fakeId, nil
}

func (uui *UserUsecaseImpl) FindById(id int) (*entities.FetchUserDto, error) {
	// Get user from repo
	user, err := uui.userRepo.FindById(id)

	// check err is not nil
	if err != nil {
		panic("User not found")
		// return nil, nil
	}

	uid := shared.NewUID(uint32(user.BaseSQLModel.Id), int(shared.DbTypeUser), 1)
	// convert to dto
	userDto := &entities.FetchUserDto{
		FakeId:    uid.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	return userDto, nil
}

func (uui *UserUsecaseImpl) FindByEmail(email string) (*entities.FetchUserDto, error) {
	// Get user from repo
	user, err := uui.userRepo.FindByEmail(email)

	if err != nil {
		panic("\nUser not found")
	}

	userDto := PreprocessUserDto(user)

	// Check err is not nil
	// Mask the ID
	return userDto, nil
}

func (uui *UserUsecaseImpl) FindAll() ([]*entities.FetchUserDto, error) {
	result, err := uui.userRepo.FindAll()
	users := make([]*entities.FetchUserDto, len(result))

	if err != nil {
		panic("\nUsers empty or not found")
	}

	for i, _ := range users {
		users[i] = PreprocessUserDto(result[i])
	}

	return users, nil
}

func PreprocessUserDto(ent *entities.User) *entities.FetchUserDto {
	uid := shared.NewUID(uint32(ent.BaseSQLModel.Id), int(shared.DbTypeUser), 1)
	return &entities.FetchUserDto{
		FakeId:    uid.String(),
		FirstName: ent.FirstName,
		LastName:  ent.LastName,
		Email:     ent.Email,
		Role:      ent.Role,
	}
}

func NewUserUsecaseImpl(userRepo repositories.UserRepository) UserUseCase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

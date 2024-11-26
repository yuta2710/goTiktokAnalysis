package usecases

import (
	"fmt"
	"time"

	"github.com/yuta_2710/go-clean-arc-reviews/modules/auth/models"
	TodoEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	TokenRepo "github.com/yuta_2710/go-clean-arc-reviews/modules/token/repositories"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	UserEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/repositories"
	UserRepo "github.com/yuta_2710/go-clean-arc-reviews/modules/users/repositories"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type AuthUsecaseImpl struct {
	UserRepo  UserRepo.UserRepository
	TokenRepo TokenRepo.TokenRepository
}

func (aui *AuthUsecaseImpl) Login(mod *models.LoginRequest) (*models.AuthResponse, error) {
	u, err := aui.UserRepo.FindByEmail(mod.Email)

	if err != nil {
		return nil, fmt.Errorf("[Login failed]: Email is not valid")
	}

	isValidPassword := shared.CheckPasswordHash(mod.Password, u.Password)

	if !isValidPassword {
		return nil, fmt.Errorf("[Login failed]: Password is not correct")
	}

	u.Mask(shared.DbTypeUser)
	authId := u.FakeId.String()
	accessTokenString, refreshTokenString, err := shared.TokenProvider(u.Id, authId)

	if err != nil {
		return nil, fmt.Errorf("Login failed, something wrong due to processing refresh token to String type")
	}

	// Save to db
	err = aui.TokenRepo.CreateTokens(u.Id, accessTokenString, refreshTokenString, time.Now().Add(7*24*time.Hour))

	if err != nil {
		return nil, fmt.Errorf("Login failed, something wrong due to saving token to database")
	}

	return &models.AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
func (aui *AuthUsecaseImpl) SignUp(mod *models.SignUpRequest) (*models.AuthResponse, error) {
	insertData := &UserEntities.InsertUserDto{
		FirstName: mod.FirstName,
		LastName:  mod.LastName,
		Email:     mod.Email,
		Password:  mod.Password,
		Role:      mod.Role,
		IsActive:  true,
		IsAdmin:   false,
		IsBlocked: false,
		Todos:     []TodoEntities.Todo{},
	}

	userId, fakeId, err := aui.UserRepo.Insert(insertData)

	if err != nil {
		return nil, err
	}

	accessTokenString, refreshTokenString, err := shared.TokenProvider(userId, fakeId)

	if err != nil {
		return nil, fmt.Errorf("Login failed, something wrong due to processing refresh token to String type")
	}
	// fmt.Println("Hahahahahaha")
	// fmt.Println(accessTokenString, refreshTokenString)

	// Save to db
	err = aui.TokenRepo.CreateTokens(userId, accessTokenString, refreshTokenString, time.Now().Add(7*24*time.Hour))

	if err != nil {
		return nil, fmt.Errorf("Login failed, something wrong due to saving token to database")
	}

	return &models.AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (aui *AuthUsecaseImpl) Profile() (*entities.FetchUserDto, error) {
	return nil, nil
}
func (aui *AuthUsecaseImpl) SignOut() error {
	return nil
}

func NewAuthUsecaseImpl(userRepo repositories.UserRepository, tokenRepo TokenRepo.TokenRepository) AuthUsecase {
	return &AuthUsecaseImpl{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	}
}

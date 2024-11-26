package usecases

import (
	"github.com/yuta_2710/go-clean-arc-reviews/modules/auth/models"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
)

type AuthUsecase interface {
	Login(mod *models.LoginRequest) (*models.AuthResponse, error)
	SignUp(mod *models.SignUpRequest) (*models.AuthResponse, error)
	Profile() (*entities.FetchUserDto, error)
	SignOut() error
}

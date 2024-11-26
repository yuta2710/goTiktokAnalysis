package models

import "github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"

type (
	LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	AuthResponse struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}

	SignUpRequest struct {
		FirstName string            `json:"firstName"`
		LastName  string            `json:"lastName"`
		Email     string            `json:"email"`
		Password  string            `json:"password"`
		Role      entities.UserRole `json:"role"`
	}
)

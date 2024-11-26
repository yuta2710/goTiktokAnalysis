package repositories

import (
	"fmt"
	"time"

	"github.com/yuta_2710/go-clean-arc-reviews/database"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/token/entities"
)

type TokenPostgresRepository struct {
	db database.Database
}

func (tpr *TokenPostgresRepository) CreateTokens(userId int, accessToken string, refreshToken string, expiredAt time.Time) error {
	tokenProvider := entities.TokenProvider{
		UserId:       userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
	}

	result := tpr.db.GetDb().Create(&tokenProvider)

	if result.Error != nil {
		fmt.Printf("Error inserting token provider: %v\n", result.Error)
	}

	return result.Error
}

func (tpr *TokenPostgresRepository) ValidateRefreshToken(refreshToken string) error {
	return nil
}
func (tpr *TokenPostgresRepository) DeleteTokens(authId string) error {
	return nil
}

func NewTokenPostgresRepository(db database.Database) TokenRepository {
	return &TokenPostgresRepository{
		db: db,
	}
}

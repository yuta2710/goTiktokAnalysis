package repositories

import "time"

type TokenRepository interface {
	CreateTokens(userId int, accessToken string, refreshToken string, expiredAt time.Time) error
	ValidateRefreshToken(refreshToken string) error
	DeleteTokens(authId string) error
}

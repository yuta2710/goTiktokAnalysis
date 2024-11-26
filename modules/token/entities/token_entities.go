package entities

import "time"

type TokenProvider struct {
	Id           int       `gorm:"column:id;primaryKey"`
	UserId       int       `gorm:"column:user_id;not null"`
	AccessToken  string    `gorm:"column:access_token;not null"`
	RefreshToken string    `gorm:"column:refresh_token;not null"`
	ExpiredAt    time.Time `gorm:"column:expired_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

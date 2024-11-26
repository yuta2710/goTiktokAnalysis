package models

import (
	TodoEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	UserEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
)

type InsertUserRequest struct {
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Email     string                `json:"email"`
	Password  string                `json:"password"`
	Role      UserEntities.UserRole `json:"role"`
	IsActive  bool                  `json:"isActive"`
	IsAdmin   bool                  `json:"isAdmin"`
	IsBlocked bool                  `json:"isBlocked"`
	IsDeleted bool                  `gorm:"column:is_deleted" json:"isDeleted"`
	Todos     []TodoEntities.Todo   `gorm:"foreignKey:AuthId;constraint:OnDelete:CASCADE;" json:"todos"`
}

type GetUserByIdRequest struct {
	Id string `json:"id"`
}

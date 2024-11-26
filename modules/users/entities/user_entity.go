package entities

import (
	TodoEntities "github.com/yuta_2710/go-clean-arc-reviews/modules/todo/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type UserRole string

// superadmin, admin, moderator, user, guest
const (
	superadmin UserRole = "superadmin"
	admin      UserRole = "admin"
	moderator  UserRole = "moderator"
	user       UserRole = "user"
	guest      UserRole = "guest"
)

// IsValid checks if a given role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case superadmin, admin, moderator, user, guest:
		return true
	default:
		return false
	}
}

func (r UserRole) HasPermission(required UserRole) bool {
	roles := map[UserRole]int{
		superadmin: 5,
		admin:      4,
		moderator:  3,
		user:       2,
		guest:      1,
	}

	return roles[r] >= roles[required]
}

type (
	User struct {
		shared.BaseSQLModel
		FirstName string              `gorm:"column:first_name" json:"firstName"`
		LastName  string              `gorm:"column:last_name" json:"lastName"`
		Email     string              `gorm:"column:email;unique" json:"email"`
		Password  string              `gorm:"column:password" json:"password"`
		Role      UserRole            `gorm:"column:role" json:"role"`
		IsActive  bool                `gorm:"column:is_active" json:"isActive"`
		IsAdmin   bool                `gorm:"column:is_admin" json:"isAdmin"`
		IsBlocked bool                `gorm:"column:is_blocked" json:"isBlocked"`
		IsDeleted bool                `gorm:"column:is_deleted" json:"isDeleted"`
		Todos     []TodoEntities.Todo `gorm:"many2many:todo_members;joinForeignKey:UserId;joinReferences:TodoId" json:"todos"`
	}

	InsertUserDto struct {
		FirstName string              `json:"firstName"`
		LastName  string              `json:"lastName"`
		Email     string              `json:"email"`
		Password  string              `json:"password"`
		Role      UserRole            `json:"role"`
		IsActive  bool                `json:"isActive"`
		IsAdmin   bool                `json:"isAdmin"`
		IsBlocked bool                `json:"isBlocked"`
		IsDeleted bool                `gorm:"column:is_deleted" json:"isDeleted"`
		Todos     []TodoEntities.Todo `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;" json:"todos"`
	}

	FetchUserDto struct {
		FakeId    string   `json:"fakeId"`
		FirstName string   `json:"firstName"`
		LastName  string   `json:"lastName"`
		Email     string   `json:"email"`
		Role      UserRole `json:"role"`
	}
)

func NewInsertUserRequest(firstName, lastName, email, password string, role UserRole) *InsertUserDto {
	// Default role to "user" if not provided
	if role == "" {
		role = UserRole("user")
	} else {
		role = UserRole(role)
	}

	return &InsertUserDto{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Role:      role,
		IsActive:  true, // Set any other default values here
		IsAdmin:   false,
		IsBlocked: false,
		IsDeleted: false,
		Todos:     []TodoEntities.Todo{},
	}
}

func (u *User) Mask(dbType shared.DbType) {
	u.BaseSQLModel.Mask(dbType)
}

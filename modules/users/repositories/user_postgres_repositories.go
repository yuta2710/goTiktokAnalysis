package repositories

import (
	"fmt"
	"log"

	"github.com/yuta_2710/go-clean-arc-reviews/database"
	"github.com/yuta_2710/go-clean-arc-reviews/modules/users/entities"
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type UserPostgresRepository struct {
	db database.Database
}

// func (ur *UserPostgresRepository) Insert(dto *entities.InsertUserDto) (int, string, error) {
// 	insert := &entities.User{
// 		FirstName: dto.FirstName,
// 		LastName:  dto.LastName,
// 		Email:     dto.Email,
// 		Password:  dto.Password, // Ensure password is hashed before calling this
// 		Role:      dto.Role,
// 		Todos:     dto.Todos,
// 	}

// 	hash, err := shared.HashPassword(insert.Password)

// 	if err != nil {
// 		log.Fatal("Error hashing password")
// 	}
// 	if hash != "" {
// 		insert.Password = hash
// 	}

// 	result := ur.db.GetDb().Create(insert)

// 	if result.Error != nil {
// 		if strings.Contains(result.Error.Error(), "unique constraint") {
// 			fmt.Println("This email already exists")
// 			panic(result.Error.Error())
// 		}
// 	}

// 	fmt.Println("[INSERTED DATA SUCCESSFULLY]")
// 	insert.Mask(shared.DbTypeUser)

// 	fmt.Printf("User ID : %s, Fake ID: %s", insert.Id, insert.FakeId.String())

// 	return insert.Id, insert.FakeId.String(), nil
// }

func (ur *UserPostgresRepository) Insert(dto *entities.InsertUserDto) (int, string, error) {
	// Prepare the User entity
	insert := &entities.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      dto.Role,
		Todos:     dto.Todos,
	}

	// Hash the password
	hash, err := shared.HashPassword(insert.Password)
	if err != nil {
		log.Fatal("Error hashing password")
	}
	insert.Password = hash

	var id int
	// Use raw SQL to insert and get the ID
	sql := `
		INSERT INTO users (first_name, last_name, email, password, role, is_active, is_admin, is_blocked, is_deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id
	`

	// Execute the query and scan the ID
	err = ur.db.GetDb().Raw(sql,
		insert.FirstName,
		insert.LastName,
		insert.Email,
		insert.Password,
		insert.Role,
		insert.IsActive,  // is_active default
		insert.IsAdmin,   // is_admin default
		insert.IsBlocked, // is_blocked default
		insert.IsDeleted, // is_deleted default
	).Scan(&id).Error

	if err != nil {
		return 0, "", fmt.Errorf("failed to insert user: %v", err)
	}

	// Assign the ID to the user entity
	insert.Id = id

	// Mask the ID into a FakeId
	insert.Mask(shared.DbTypeUser)

	fmt.Printf("[INSERTED DATA SUCCESSFULLY] User ID: %d, Fake ID: %s\n", insert.Id, insert.FakeId.String())

	return insert.Id, insert.FakeId.String(), nil
}

func (ur *UserPostgresRepository) InsertBatch(dtos []*entities.InsertUserDto) error {
	return nil
}

func (ur *UserPostgresRepository) FindById(id int) (*entities.User, error) {
	// Debug log
	fmt.Printf("Looking for user with ID: %d\n", id)

	// Use a struct, not an uninitialized pointer
	var u *entities.User

	// Enable GORM debugging to log SQL
	result := ur.db.GetDb().Debug().Where("id = ?", id).First(&u)

	if result.Error != nil {
		fmt.Printf("Error finding user: %v\n", result.Error)
		return nil, result.Error
	}

	return u, nil
}

func (ur *UserPostgresRepository) FindByEmail(email string) (*entities.User, error) {
	var u entities.User // Use a struct instead of a pointer
	fmt.Printf("Email is %s\n", email)

	// Use LOWER to ensure case-insensitive comparison
	// result := ur.db.GetDb().Debug().Where("LOWER(email) = ?", strings.ToLower(email)).First(&u)
	result := ur.db.GetDb().Where("email = ?", email).First(&u)

	if result.Error != nil {
		fmt.Printf("Error finding user by email: %v\n", result.Error)
		return nil, result.Error
	}

	fmt.Printf("Found user: %+v\n", u)
	return &u, nil
}

func (ur *UserPostgresRepository) FindAll() ([]*entities.User, error) {
	var users []*entities.User
	result := ur.db.GetDb().Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func NewUserPostgresRepository(db database.Database) UserRepository {
	return &UserPostgresRepository{
		db: db,
	}
}

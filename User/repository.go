package user

import (
	"errors"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
	UpdateUser(user *User, existingUID uint) error
	GetAllUsers() ([]User, error)
	AuthenticateUser(username, passwd string) (*User, error)
	GetUserRole(username string) (string, error)
}

type SQLUserRepository struct {
	DB *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{DB: db}
}
func (repo *SQLUserRepository) CreateUser(user *User) error {
	insertQuery := "INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, user.Username, user.Email, user.Password, user.Role).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) GetUserByID(userID uint) (*User, error) {
	var user User
	if err := repo.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (repo *SQLUserRepository) DeleteUser(userID uint) error {
	deleteQuery := "DELETE FROM users WHERE id = ?"
	if err := repo.DB.Exec(deleteQuery, userID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) UpdateUser(user *User, existingUID uint) error {
	updateQuery := "UPDATE users SET username = ?, email = ?, password = ?,  role = ? WHERE id = ?"
	if err := repo.DB.Exec(updateQuery, user.Username, user.Email, user.Password, user.Role, existingUID).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLUserRepository) GetAllUsers() ([]User, error) {
	var users []User
	if err := repo.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *SQLUserRepository) AuthenticateUser(username, password string) (*User, error) {
	var user User
	if err := repo.DB.Where("username = ? AND password = ?", username, password).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

func (repo *SQLUserRepository) GetUserRole(username string) (string, error) {
	var user User

	// Find the user by username
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // User not found, return empty role
		}
		return "", err
	}

	return user.Role, nil
}

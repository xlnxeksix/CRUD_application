package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(userID uint) (*User, error)
	DeleteUser(userID uint) error
	UpdateUser(user *User, existingUID uint) error
	GetAllUsers() ([]User, error)
}

type SQLUserRepository struct {
	DB *gorm.DB
}

func NewSQLUserRepository(db *gorm.DB) UserRepository {
	return &SQLUserRepository{DB: db}
}
func (repo *SQLUserRepository) CreateUser(user *User) error {
	insertQuery := "INSERT INTO users (id, username, email, role) VALUES (?, ?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, user.ID, user.Username, user.Email, user.Role).Error; err != nil {
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
	updateQuery := "UPDATE users SET id = ?, username = ?, email = ?, role = ? WHERE id = ?"
	if err := repo.DB.Exec(updateQuery, user.ID, user.Username, user.Email, user.Role, existingUID).Error; err != nil {
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

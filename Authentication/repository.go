package Authentication

import (
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	AuthenticateUser(username, passwd string) (string, error)
}

type SQLAuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &SQLAuthRepository{DB: db}
}
func (repo *SQLAuthRepository) AuthenticateUser(username, password string) (string, error) {
	var role string
	query := "SELECT role FROM users WHERE username = ? AND password = ?"

	result := repo.DB.Raw(query, username, password).Scan(&role)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // User not found
		}
		return "", result.Error
	}

	return role, nil
}

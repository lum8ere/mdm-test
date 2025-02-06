package repositories

import (
	"errors"
	"mdm/libs/2_generated_models/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByUsername(username string) (*model.User, error)
}

type user_repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &user_repository{db: db}
}

func (r *user_repository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

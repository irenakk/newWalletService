package repository

import (
	"newWalletService/config"
	"newWalletService/internal/model"
)

type InterfaceUserRepository interface {
	Find(username string) (model.User, error)
}

type UserRepository struct {
	db *config.Database
}

func NewUserRepository(db *config.Database) InterfaceUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Find(username string) (model.User, error) {
	var user model.User
	err := r.db.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

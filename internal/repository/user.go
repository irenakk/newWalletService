package repository

import (
	"newWalletService/config"
	"newWalletService/internal/model"
)

type InterfaceUserRepository interface {
	Create(user model.UserRegister) (int, error)
	Find(username string) (model.User, error)
	Exists(username string) (bool, error)
}

type UserRepository struct {
	db *config.Database
}

func NewUserRepository(db *config.Database) InterfaceUserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(user model.UserRegister) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`,
		user.Username, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) Find(username string) (model.User, error) {
	var user model.User
	err := r.db.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = $1`, username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Exists(username string) (bool, error) {
	var exists bool
	err := r.db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

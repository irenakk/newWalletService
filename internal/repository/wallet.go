package repository

import (
	"walletService/config"
)

type InterfaceWalletRepository interface {
	Create(userId int) (int, error)
	FindByUserId(userId int) (int, error)
}

type WalletRepository struct {
	db *config.Database
}

func NewWalletRepository(db *config.Database) InterfaceWalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) Create(userId int) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`INSERT INTO wallet (user_id) VALUES ($1) RETURNING id`, userId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *WalletRepository) FindByUserId(userId int) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`SELECT id FROM wallet WHERE user_id = $1`, userId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

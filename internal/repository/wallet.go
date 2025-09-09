package repository

import (
	"context"
	"walletService/config"
)

type InterfaceWalletRepository interface {
	Create(ctx context.Context, userId int) (int, error)
	Delete(ctx context.Context, userId int) error
	FindByUserId(userId int) (int, error)
}

type WalletRepository struct {
	db *config.Database
}

func NewWalletRepository(db *config.Database) InterfaceWalletRepository {
	return &WalletRepository{db}
}

func (r *WalletRepository) Create(ctx context.Context, userId int) (int, error) {
	var id int
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	err = r.db.DB.QueryRow(`INSERT INTO wallet (user_id) VALUES ($1) RETURNING id`, userId).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (r *WalletRepository) Delete(ctx context.Context, userId int) error {
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM wallet WHERE id = $1`, userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *WalletRepository) FindByUserId(userId int) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`SELECT id FROM wallet WHERE user_id = $1`, userId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

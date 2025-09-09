package repository

import (
	"context"
	"walletService/config"
	"walletService/internal/model"
)

type InterfaceAccountRepository interface {
	Create(ctx context.Context, currency string, walletId int) (int, error)
	Delete(ctx context.Context, currency string, walletId int) (int, error)
	FindByWalletCurrency(walletId int, currency string) (int, error)
	Add(id int, amount int) (model.Account, error)
	Subtraction(id int, amount int) (model.Account, error)
}

type AccountRepository struct {
	db *config.Database
}

func NewAccountRepository(db *config.Database) InterfaceAccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) Create(ctx context.Context, currency string, walletId int) (int, error) {
	var id int
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	err = r.db.DB.QueryRow(`INSERT INTO account (balance, currency, wallet_id) VALUES ($1, $2, $3) RETURNING id`,
		0, currency, walletId).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (r *AccountRepository) Delete(ctx context.Context, currency string, walletId int) (int, error) {
	var id int
	tx, err := r.db.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	err = r.db.DB.QueryRow(`DELETE FROM wallet WHERE wallet_id = $1 AND currency = $2 RETURNING id`,
		walletId, currency).Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (r *AccountRepository) FindByWalletCurrency(walletId int, currency string) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`SELECT id FROM account WHERE wallet_id = $1 AND currency = $2`,
		walletId, currency).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountRepository) Add(id int, amount int) (model.Account, error) {
	var updatedAccount model.Account
	err := r.db.DB.QueryRow(`UPDATE account SET balance = balance + $1 WHERE id = $2 RETURNING id, balance, currency, wallet_id`,
		amount, id).Scan(&updatedAccount.ID, &updatedAccount.Balance, &updatedAccount.Currency, &updatedAccount.WalletId)
	if err != nil {
		return model.Account{}, err
	}
	return updatedAccount, nil
}

func (r *AccountRepository) Subtraction(id int, amount int) (model.Account, error) {
	var updatedAccount model.Account
	err := r.db.DB.QueryRow(`UPDATE account SET balance = balance - $1 WHERE id = $2 RETURNING id, balance, currency, wallet_id`,
		amount, id).Scan(&updatedAccount.ID, &updatedAccount.Balance, &updatedAccount.Currency, &updatedAccount.WalletId)
	if err != nil {
		return model.Account{}, err
	}
	return updatedAccount, nil
}

package repository

import (
	"newWalletService/config"
	"newWalletService/internal/model"
)

type InterfaceAccountRepository interface {
	Create(currency string, walletId int) (int, error)
	FindByAccountId(id int) (model.Account, error)
	FindByWalletId(walletId int) ([]model.Account, error)
	Add(id int, newBalance int) (model.Account, error)
	Transfer(idFrom int, idTo int, transferSum int) (model.Account, error)
}

type AccountRepository struct {
	db *config.Database
}

func NewAccountRepository(db *config.Database) InterfaceAccountRepository {
	return &AccountRepository{db}
}

func (r *AccountRepository) Create(currency string, walletId int) (int, error) {
	var id int
	err := r.db.DB.QueryRow(`INSERT INTO account (balance, currency, wallet_id) VALUES ($1, $2, $3) RETURNING id`,
		0, currency, walletId).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AccountRepository) FindByAccountId(id int) (model.Account, error) {
	var account model.Account
	err := r.db.DB.QueryRow(`SELECT id, balance, currency, wallet_id FROM accounts WHERE id = $1`, id).Scan(
		&account.ID, &account.Balance, &account.Currency, &account.WalletId)
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (r *AccountRepository) FindByWalletId(walletId int) ([]model.Account, error) {
	var accounts []model.Account
	rows, err := r.db.DB.Query(`SELECT id, balance, currency, wallet_id FROM accounts WHERE wallet_id = $1`, walletId)
	if err != nil {
		return []model.Account{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var account model.Account
		err := rows.Scan(&account.ID, &account.Balance, &account.Currency, &account.WalletId)
		if err != nil {
			return []model.Account{}, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (r *AccountRepository) Add(id int, newBalance int) (model.Account, error) {
	var updatedAccount model.Account
	err := r.db.DB.QueryRow(`UPDATE wallet SET balance = $1 WHERE id = $2`, newBalance, id).Scan(
		&updatedAccount.ID, &updatedAccount.Balance, &updatedAccount.Currency, &updatedAccount.WalletId)
	if err != nil {
		return model.Account{}, err
	}
	return updatedAccount, nil
}

func (r *AccountRepository) Transfer(idFrom int, idTo int, transferSum int) (model.Account, error) {
	var updatedAccount model.Account
	err := r.db.DB.QueryRow(`UPDATE wallet SET balance = balance - $1 WHERE id = $2`, transferSum, idFrom).Scan(
		&updatedAccount.ID, &updatedAccount.Balance, &updatedAccount.Currency, &updatedAccount.WalletId)
	if err != nil {
		return model.Account{}, err
	}
	err1 := r.db.DB.QueryRow(`UPDATE wallet SET balance = balance + $1 WHERE id = $2`, transferSum, idTo).Scan(
		&updatedAccount.ID, &updatedAccount.Balance, &updatedAccount.Currency, &updatedAccount.WalletId)
	if err1 != nil {
		return model.Account{}, err
	}
	return updatedAccount, nil
}

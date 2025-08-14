package usecase

import (
	"newWalletService/internal/model"
	"newWalletService/internal/repository"
)

type InterfaceAccountUsecase interface {
	Create(currency string, walletId int) (int, error)
	FindByAccountId(id int) (model.Account, error)
	FindByWalletId(walletId int) ([]model.Account, error)
	Add(account model.Account, newBalance int) (model.Account, error)
	Transfer(account model.Account, newBalance int) (model.Account, error)
}

type AccountUsecase struct {
	accountRepository repository.InterfaceAccountRepository
}

func NewAccountUsecase(accountRepository repository.InterfaceAccountRepository) *AccountUsecase {
	return &AccountUsecase{accountRepository}
}

func (usecase AccountUsecase) Create(currency string, walletId int) (int, error) {
	id, err := usecase.accountRepository.Create(currency, walletId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (usecase AccountUsecase) FindByAccountId(id int) (model.Account, error) {
	account, err := usecase.accountRepository.FindByAccountId(id)
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (usecase AccountUsecase) FindByWalletId(walletId int) ([]model.Account, error) {
	account, err := usecase.accountRepository.FindByWalletId(walletId)
	if err != nil {
		return []model.Account{}, err
	}
	return account, nil
}

func (usecase AccountUsecase) Add(account model.Account, newBalance int) (model.Account, error) {
	account, err := usecase.accountRepository.Add(account.ID, newBalance)
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (usecase AccountUsecase) Transfer(account model.Account, newBalance int) (model.Account, error) {
	account, err := usecase.accountRepository.AddBalance(account.ID, newBalance)
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

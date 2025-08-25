package usecase

import (
	"errors"
	"newWalletService/internal/model"
	"newWalletService/internal/repository"
	"strings"
)

type InterfaceWalletUsecase interface {
	Add(username string, currency string, amount int) (model.Account, error)
	Transfer(username string, currency string, amount int) (model.Account, error)
}

type WalletUsecase struct {
	userRepository    repository.InterfaceUserRepository
	walletRepository  repository.InterfaceWalletRepository
	accountRepository repository.InterfaceAccountRepository
}

func NewWalletUsecase(userRepository repository.InterfaceUserRepository,
	walletRepository repository.InterfaceWalletRepository,
	accountRepository repository.InterfaceAccountRepository) *WalletUsecase {
	return &WalletUsecase{userRepository, walletRepository, accountRepository}
}

func (usecase WalletUsecase) Add(username string, currency string, amount int) (model.Account, error) {
	user, err := usecase.userRepository.Find(username)
	if err != nil {
		return model.Account{}, err
	}

	walletId, err := usecase.walletRepository.FindByUserId(user.ID)
	if err != nil {
		return model.Account{}, err
	}

	if strings.ToUpper(currency) == "USD" || strings.ToUpper(currency) == "EUR" || strings.ToUpper(currency) == "RUB" {
		accountId, err := usecase.accountRepository.FindByWalletCurrency(walletId, strings.ToUpper(currency))
		if err != nil {
			return model.Account{}, err
		}

		account, err := usecase.accountRepository.Add(accountId, amount)
		if err != nil {
			return model.Account{}, err
		}

		return account, nil
	}

	return model.Account{}, errors.New("invalid currency")
}

func (usecase WalletUsecase) Transfer(usernameFrom string, usernameTo string, currency string, amount int) (model.Account, error) {
	userFrom, err := usecase.userRepository.Find(usernameFrom)
	if err != nil {
		return model.Account{}, err
	}

	walletIdFrom, err := usecase.walletRepository.FindByUserId(userFrom.ID)
	if err != nil {
		return model.Account{}, err
	}

	userTo, err := usecase.userRepository.Find(usernameTo)
	if err != nil {
		return model.Account{}, err
	}

	walletIdTo, err := usecase.walletRepository.FindByUserId(userTo.ID)
	if err != nil {
		return model.Account{}, err
	}

	if strings.ToUpper(currency) == "USD" || strings.ToUpper(currency) == "EUR" || strings.ToUpper(currency) == "RUB" {
		accountIdTo, err := usecase.accountRepository.FindByWalletCurrency(walletIdTo, strings.ToUpper(currency))
		if err != nil {
			return model.Account{}, err
		}

		_, err = usecase.accountRepository.Add(accountIdTo, amount)
		if err != nil {
			return model.Account{}, err
		}

		accountIdFrom, err := usecase.accountRepository.FindByWalletCurrency(walletIdFrom, strings.ToUpper(currency))
		if err != nil {
			return model.Account{}, err
		}

		accountFrom, err := usecase.accountRepository.Subtraction(accountIdFrom, amount)
		if err != nil {
			return model.Account{}, err
		}

		return accountFrom, nil
	}

	return model.Account{}, errors.New("invalid currency")
}

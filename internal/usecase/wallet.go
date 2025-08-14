package usecase

import (
	"newWalletService/internal/model"
	"newWalletService/internal/repository"
)

type InterfaceWalletUsecase interface {
	Create(userId int) (int, error)
	FindByWalletId(id int) (model.Wallet, error)
	FindByUserId(userId int) (model.Wallet, error)
}

type WalletUsecase struct {
	walletRepository repository.InterfaceWalletRepository
}

func NewWalletUsecase(walletRepository repository.InterfaceWalletRepository) *WalletUsecase {
	return &WalletUsecase{walletRepository}
}

func (usecase WalletUsecase) Create(userId int) (int, error) {
	id, err := usecase.walletRepository.Create(userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (usecase WalletUsecase) FindByWalletId(id int) (model.Wallet, error) {
	wallet, err := usecase.walletRepository.FindByWalletId(id)
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}

func (usecase WalletUsecase) FindByUserId(userId int) (model.Wallet, error) {
	wallet, err := usecase.walletRepository.FindByUserId(userId)
	if err != nil {
		return model.Wallet{}, err
	}
	return wallet, nil
}

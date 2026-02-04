package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
	"time"
	"walletService/internal/model"
	"walletService/internal/repository"
	"walletService/internal/service"
)

type InterfaceWalletUsecase interface {
	CreateWallet(ctx context.Context, userId int) (int, error)
	DeleteWallet(ctx context.Context, userId int) error
	CreateAccount(ctx context.Context, walletId int, currency string) (int, error)
	DeleteAccount(ctx context.Context, walletId int, currency string) (int, error)
	Add(username string, currency string, amount int) (model.Account, error)
	Transfer(username string, currency string, amount int) (model.Account, error)
	Notify(username string, message string) error
}

type WalletUsecase struct {
	userService       service.UserGrpcService
	walletRepository  repository.InterfaceWalletRepository
	accountRepository repository.InterfaceAccountRepository
	kafkaWriter       *kafka.Writer
}

func NewWalletUsecase(userService service.UserGrpcService,
	walletRepository repository.InterfaceWalletRepository,
	accountRepository repository.InterfaceAccountRepository,
	kafkaWriter *kafka.Writer) *WalletUsecase {
	return &WalletUsecase{userService, walletRepository, accountRepository, kafkaWriter}
}

func (usecase WalletUsecase) CreateWallet(ctx context.Context, userId int) (int, error) {
	id, err := usecase.walletRepository.Create(ctx, userId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (usecase WalletUsecase) DeleteWallet(ctx context.Context, userId int) error {
	err := usecase.walletRepository.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (usecase WalletUsecase) CreateAccount(ctx context.Context, walletId int, currency string) (int, error) {
	id, err := usecase.accountRepository.Create(ctx, currency, walletId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (usecase WalletUsecase) DeleteAccount(ctx context.Context, walletId int, currency string) (int, error) {
	id, err := usecase.accountRepository.Delete(ctx, currency, walletId)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (usecase WalletUsecase) Add(username string, currency string, amount int) (model.Account, error) {
	user, err := usecase.userService.Find(username)
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

		err = usecase.Notify(username, fmt.Sprintf("Ваш счет %s пополнен на %d", currency, amount))
		if err != nil {
			return model.Account{}, err
		}

		return account, nil
	}

	return model.Account{}, errors.New("invalid currency")
}

func (usecase WalletUsecase) Transfer(usernameFrom string, usernameTo string, currency string, amount int) (model.Account, error) {
	userFrom, err := usecase.userService.Find(usernameFrom)
	if err != nil {
		return model.Account{}, err
	}

	walletIdFrom, err := usecase.walletRepository.FindByUserId(userFrom.ID)
	if err != nil {
		return model.Account{}, err
	}

	userTo, err := usecase.userService.Find(usernameTo)
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

		err = usecase.Notify(usernameTo, fmt.Sprintf("Вам перевели %d %s от %s", amount, currency, usernameFrom))
		if err != nil {
			return model.Account{}, err
		}
		err = usecase.Notify(usernameFrom, fmt.Sprintf("Вы отправили %d %s пользователю %s", amount, currency, usernameTo))
		if err != nil {
			return model.Account{}, err
		}

		return accountFrom, nil
	}

	return model.Account{}, errors.New("invalid currency")
}

func (usecase WalletUsecase) Notify(username string, message string) error {
	messages := []kafka.Message{
		{Key: []byte(username), Value: []byte(message)},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, msg := range messages {
		err := usecase.kafkaWriter.WriteMessages(ctx, msg)
		if err != nil {
			log.Fatalf("Ошибка записи: %v", err)
			return err
		}
		fmt.Printf("Отправлено: key=%s, value=%s\n", msg.Key, msg.Value)
	}

	fmt.Println("Все сообщения отправлены!")
	return nil
}

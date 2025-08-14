package usecase

import (
	"github.com/golang-jwt/jwt/v5"
	"newWalletService/internal/model"
	"newWalletService/internal/repository"
	"newWalletService/internal/utils"
	"time"
)

type InterfaceUserUsecase interface {
	Create(user model.UserRegister) (int, error)
	Find(username string) (model.User, error)
	Exists(username string) (bool, error)
	CheckPassword(loginPassword string, userPassword string) bool
	GenerateJWT(user model.User, tokenExpiration time.Duration, jwtSecret []byte) (string, error)
}

type UserUsecase struct {
	userRepository    repository.InterfaceUserRepository
	walletRepository  repository.InterfaceWalletRepository
	accountRepository repository.InterfaceAccountRepository
}

func NewUserUsecase(userRepository repository.InterfaceUserRepository,
	walletRepository repository.InterfaceWalletRepository,
	accountRepository repository.InterfaceAccountRepository) *UserUsecase {
	return &UserUsecase{userRepository, walletRepository, accountRepository}
}

func (usecase UserUsecase) Create(user model.UserRegister) (int, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}

	u := model.UserRegister{
		Username: user.Username,
		Password: hashedPassword,
	}

	id, err := usecase.userRepository.Create(u)
	if err != nil {
		return 0, err
	}

	walletId, walletErr := usecase.walletRepository.Create(id)
	if walletErr != nil {
		return 0, walletErr
	}

	_, accountUsdError := usecase.accountRepository.Create("USD", walletId)
	if accountUsdError != nil {
		return 0, walletErr
	}

	_, accountEurError := usecase.accountRepository.Create("EUR", walletId)
	if accountEurError != nil {
		return 0, walletErr
	}

	_, accountRubError := usecase.accountRepository.Create("RUB", walletId)
	if accountRubError != nil {
		return 0, walletErr
	}

	return id, nil
}

func (usecase UserUsecase) Find(username string) (model.User, error) {
	user, err := usecase.userRepository.Find(username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (usecase UserUsecase) Exists(username string) (bool, error) {
	exists, err := usecase.userRepository.Exists(username)
	if err != nil {
		return true, err
	}
	return exists, nil
}

func (usecase UserUsecase) CheckPassword(loginPassword string, userPassword string) bool {
	return utils.CheckPasswordHash(loginPassword, userPassword)
}

func (usecase UserUsecase) GenerateJWT(user model.User, tokenExpiration time.Duration, jwtSecret []byte) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"username": user.Username,
		"iat":      now.Unix(),
		"exp":      now.Add(tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}

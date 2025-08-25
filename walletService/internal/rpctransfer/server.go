package rpctransfer

import (
	"newWalletService/internal/usecase"
	"newWalletService/proto/wallet"
)

type Handlers struct {
	Usecase *usecase.WalletUsecase
	wallet.UnimplementedWalletServiceServer
}

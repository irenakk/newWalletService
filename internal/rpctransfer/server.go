package rpctransfer

import (
	"walletService/internal/usecase"
	"walletService/proto/server"
)

type Handlers struct {
	Usecase *usecase.WalletUsecase
	server.UnimplementedWalletServiceServer
}

func (h *Handlers) mustEmbedUnimplementedWalletServiceServer() {
	//TODO implement me
	panic("implement me")
}

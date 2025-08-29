package rpctransfer

import (
	"newWalletService/internal/usecase"
	"newWalletService/proto"
)

type Handlers struct {
	Usecase *usecase.UserUsecase
	proto.UnimplementedUserServiceServer
}

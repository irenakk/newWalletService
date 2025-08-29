package rpctransfer

import (
	"context"
	"newWalletService/proto"
)

func (h *Handlers) CreateAccount(ctx context.Context, in *proto.AccountRequest) (*proto.AccountResponse, error) {
	id, err := h.Usecase.CreateAccount(int(in.WalletId), in.Currency)
	if err != nil {
		return &proto.AccountResponse{}, err
	}

	return &proto.AccountResponse{AccountId: int64(id)}, nil
}

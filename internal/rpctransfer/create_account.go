package rpctransfer

import (
	"context"
	"walletService/proto/server"
)

func (h *Handlers) CreateAccount(ctx context.Context, in *server.AccountRequest) (*server.AccountResponse, error) {
	id, err := h.Usecase.CreateAccount(ctx, int(in.WalletId), in.Currency)
	if err != nil {
		return &server.AccountResponse{}, err
	}

	return &server.AccountResponse{AccountId: int64(id)}, nil
}

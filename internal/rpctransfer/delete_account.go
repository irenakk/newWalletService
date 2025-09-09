package rpctransfer

import (
	"context"
	"walletService/proto/server"
)

func (h *Handlers) DeleteAccount(ctx context.Context, in *server.AccountRequest) (*server.AccountResponse, error) {
	id, err := h.Usecase.DeleteAccount(ctx, int(in.WalletId), in.Currency)
	if err != nil {
		return &server.AccountResponse{}, err
	}

	return &server.AccountResponse{AccountId: int64(id)}, nil
}

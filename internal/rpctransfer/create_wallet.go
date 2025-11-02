package rpctransfer

import (
	"context"
	"walletService/proto/server"
)

func (h *Handlers) CreateWallet(ctx context.Context, in *server.WalletRequest) (*server.WalletResponse, error) {
	id, err := h.Usecase.CreateWallet(ctx, int(in.UserId))
	if err != nil {
		return &server.WalletResponse{}, err
	}

	return &server.WalletResponse{WalletId: int64(id)}, nil
}

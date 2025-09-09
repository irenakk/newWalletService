package rpctransfer

import (
	"context"
	"walletService/proto/server"
)

func (h *Handlers) DeleteWallet(ctx context.Context, in *server.WalletRequest) (*server.WalletResponse, error) {
	err := h.Usecase.DeleteWallet(ctx, int(in.UserId))
	if err != nil {
		return &server.WalletResponse{}, err
	}

	return &server.WalletResponse{WalletId: 0}, nil
}

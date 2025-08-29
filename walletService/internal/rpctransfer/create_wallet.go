package rpctransfer

import (
	"context"
	"newWalletService/proto"
)

func (h *Handlers) CreateWallet(ctx context.Context, in *proto.WalletRequest) (*proto.WalletResponse, error) {
	id, err := h.Usecase.CreateWallet(int(in.UserId))
	if err != nil {
		return &proto.WalletResponse{}, err
	}

	return &proto.WalletResponse{WalletId: int64(id)}, nil
}

package rpctransfer

import (
	"context"
	"newWalletService/proto/wallet"
)

func (h *Handlers) Add(ctx context.Context, in *wallet.AddRequest) (*wallet.AddResponse, error) {
	_, err := h.Usecase.Add(in.Username, in.Currency, int(in.Amount))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

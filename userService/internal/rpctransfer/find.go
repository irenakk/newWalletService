package rpctransfer

import (
	"context"
	"newWalletService/internal/model"
	"newWalletService/proto"
)

func (h *Handlers) Find(ctx context.Context, in *proto.FindRequest) (*proto.FindResponse, error) {
	var user model.User
	user, err := h.Usecase.Find(in.Username)
	if err != nil {
		return &proto.FindResponse{}, err
	}

	return &proto.FindResponse{UserId: int64(user.ID), Username: user.Username}, nil
}

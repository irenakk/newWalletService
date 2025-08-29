package grpcclient

import (
	"context"
	"newWalletService/internal/model"
	"newWalletService/proto"
)

type UserGrpcRepository struct {
	client proto.UserServiceClient
}

func NewUserGrpcRepository(client proto.UserServiceClient) *UserGrpcRepository {
	return &UserGrpcRepository{client: client}
}

func (r *UserGrpcRepository) Find(username string) (model.User, error) {
	resp, err := r.client.Find(context.Background(), &proto.FindRequest{Username: username})
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:       int(resp.UserId),
		Username: resp.Username,
	}, nil
}

package service

import (
	"context"
	"walletService/internal/model"
	user "walletService/proto/client"
)

type UserGrpcService struct {
	client user.UserServiceClient
}

func NewUserService(client user.UserServiceClient) *UserGrpcService {
	return &UserGrpcService{client: client}
}

func (r *UserGrpcService) Find(username string) (model.User, error) {
	resp, err := r.client.Find(context.Background(), &user.FindRequest{Username: username})
	if err != nil {
		return model.User{}, err
	}

	return model.User{
		ID:       int(resp.UserId),
		Username: resp.Username,
	}, nil
}

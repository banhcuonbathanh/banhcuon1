package handler

import (
	"context"
    pb "english-ai-full/ecomm-grpc/proto"
)

type handlercontroller struct {
	ctx        context.Context
	client     pb.EcommUserClient
	// TokenMaker *token.JWTMaker
}

func NewHandler(client pb.EcommUserClient, secretKey string) *handlercontroller {
	return &handlercontroller{
		ctx:        context.Background(),
		client:     client,
		// TokenMaker: token.NewJWTMaker(secretKey),
	}
}

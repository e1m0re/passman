package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/proto"
	"google.golang.org/grpc"
)

// AuthClient is a client to call authentication RPC.
type AuthClient struct {
	service proto.AuthServiceClient
}

// Login does login user and returns the access token.
func (client *AuthClient) Login(credentials model.Credentials) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &proto.LoginRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}

// SignUp registers new user on the server.
func (client *AuthClient) SignUp(credentials model.Credentials) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	request := &proto.SignUpRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	resp, err := client.service.SignUp(ctx, request)
	if err != nil {
		return err
	}

	if resp.Status != proto.StatusCode_SUCCESS {
		return errors.New(resp.Message)
	}

	return nil
}

// NewAuthClient initiates a new instance of AuthClient.
func NewAuthClient(cc *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		service: proto.NewAuthServiceClient(cc),
	}
}

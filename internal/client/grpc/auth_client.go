package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/proto"
)

// AuthClient is a client to call authentication RPC.
type AuthClient struct {
	service  proto.AuthServiceClient
	username string
	password string
}

// Login does login user and returns the access token.
func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &proto.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	return res.GetAccessToken(), nil
}

// SignUp registers new user on the server.
func (client *AuthClient) SignUp(ctx context.Context, credentials model.Credentials) error {
	request := &proto.SignUpRequest{
		Username: credentials.Username,
		Password: credentials.Password,
	}

	resp, err := client.service.SignUp(ctx, request)
	if err != nil {
		return err
	}

	if resp.Status == proto.StatusCode(codes.Internal) {
		return errors.New(resp.Message)
	}

	return nil
}

// NewAuthClient initiates a new instance of AuthClient.
func NewAuthClient(cc *grpc.ClientConn, username string, password string) *AuthClient {
	return &AuthClient{
		service:  proto.NewAuthServiceClient(cc),
		username: username,
		password: password,
	}
}

package grpc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/server/repository"
	"github.com/e1m0re/passman/internal/server/service/jwt"
	"github.com/e1m0re/passman/internal/server/service/users"
	"github.com/e1m0re/passman/proto"
)

type authController struct {
	jwtManager  jwt.JWTManager
	userManager users.UserManager

	proto.UnimplementedAuthServiceServer
}

func (ac *authController) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	credentials := model.Credentials{
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	}

	_, err := ac.userManager.CreateUser(ctx, credentials)
	if err != nil {
		return &proto.SignUpResponse{
			Status:  proto.StatusCode(codes.Internal),
			Message: fmt.Sprintf("cannot create user: %v", err),
		}, nil
	}

	return &proto.SignUpResponse{Status: proto.StatusCode_SUCCESS}, nil
}

func (ac *authController) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	user, err := ac.userManager.FindUserByUsername(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, repository.ErrorEntityNotFound) {
			return &proto.LoginResponse{
				Status: proto.StatusCode(codes.NotFound),
			}, err
		}

		return &proto.LoginResponse{
			Status: proto.StatusCode(codes.Internal),
		}, err
	}

	ok, err := ac.userManager.CheckPassword(ctx, *user, req.GetPassword())
	if err != nil {
		return &proto.LoginResponse{
			Status: proto.StatusCode(codes.Internal),
		}, err
	}

	if !ok {
		return &proto.LoginResponse{
			Status: proto.StatusCode(codes.InvalidArgument),
		}, err
	}

	token, err := ac.jwtManager.Generate(user)
	if err != nil {
		return &proto.LoginResponse{
			Status: proto.StatusCode(codes.Internal),
		}, err
	}

	slog.Info("the user has successfully logged into the system", slog.String("username", user.Username))
	return &proto.LoginResponse{
		Status:      proto.StatusCode_SUCCESS,
		AccessToken: token,
	}, nil
}

var _ proto.AuthServiceServer = (*authController)(nil)

// NewAuthController initiates new instance of AuthServiceServer.
func NewAuthController(jwtManager jwt.JWTManager, userManager users.UserManager) proto.AuthServiceServer {
	return &authController{
		jwtManager:  jwtManager,
		userManager: userManager,
	}
}

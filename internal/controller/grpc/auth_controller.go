package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/service/jwt"
	"github.com/e1m0re/passman/internal/service/users"
	"github.com/e1m0re/passman/proto"
)

type authController struct {
	jwtManager  jwt.JWTManager
	userManager users.UserManager

	proto.UnimplementedAuthServiceServer
}

func (uc *authController) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	credentials := model.Credentials{
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	}

	_, err := uc.userManager.CreateUser(ctx, credentials)
	if err != nil {
		return &proto.SignUpResponse{
			Status:  proto.StatusCode(codes.Internal),
			Message: fmt.Sprintf("cannot create user: %v", err),
		}, nil
	}

	return &proto.SignUpResponse{Status: proto.StatusCode_SUCCESS}, nil
}

var _ proto.AuthServiceServer = (*authController)(nil)

// NewAuthController initiates new instance of AuthServiceServer.
func NewAuthController(jwtManager jwt.JWTManager, userManager users.UserManager) proto.AuthServiceServer {
	return &authController{
		jwtManager:  jwtManager,
		userManager: userManager,
	}
}

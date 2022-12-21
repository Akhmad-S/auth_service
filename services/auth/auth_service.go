package auth

import (
	"context"
	"errors"
	"log"
	"time"

	ecom "github.com/uacademy/e_commerce/auth_service/proto-gen/e_commerce"
	"github.com/uacademy/e_commerce/auth_service/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Login ...
func (s *authService) Login(ctx context.Context, req *ecom.LoginRequest) (*ecom.TokenResponse, error) {
	log.Println("Login...")

	errAuth := errors.New("username or password wrong")

	user, err := s.stg.GetUserByUsername(req.Username)
	if err != nil {
		log.Println(err.Error())
		return nil, status.Errorf(codes.Unauthenticated, errAuth.Error())
	}

	match, err := util.ComparePassword(user.Password, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "util.ComparePassword: %s", err.Error())
	}

	if !match {
		return nil, status.Errorf(codes.Unauthenticated, errAuth.Error())
	}

	m := map[string]interface{}{
		"user_id":  user.Id,
		"username": user.Username,
	}

	tokenStr, err := util.GenerateJWT(m, 10*time.Minute, s.cfg.SecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "util.GenerateJWT: %s", err.Error())
	}

	return &ecom.TokenResponse{
		Token: tokenStr,
	}, nil
}

// HasAccess ...
func (s *authService) HasAccess(ctx context.Context, req *ecom.TokenRequest) (*ecom.HasAccessResponse, error) {
	log.Println("HasAccess...")

	result, err := util.ParseClaims(req.Token, s.cfg.SecretKey)
	if err != nil {
		log.Println(status.Errorf(codes.Unauthenticated, "util.ParseClaims: %s", err.Error()))
		return &ecom.HasAccessResponse{
			User:      nil,
			HasAccess: false,
		}, nil
	}

	log.Println(result.Username)

	user, err := s.stg.GetUserById(result.UserID)
	if err != nil {
		log.Println(status.Errorf(codes.Unauthenticated, "s.stg.GetUserById: %s", err.Error()))
		return &ecom.HasAccessResponse{
			User:      nil,
			HasAccess: false,
		}, nil
	}

	return &ecom.HasAccessResponse{
		User:      user,
		HasAccess: true,
	}, nil
}

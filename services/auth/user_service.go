package auth

import (
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/uacademy/e_commerce/auth_service/config"
	ecom "github.com/uacademy/e_commerce/auth_service/proto-gen/e_commerce"
	"github.com/uacademy/e_commerce/auth_service/storage"
	"github.com/uacademy/e_commerce/auth_service/util"

	"context"
)

type authService struct {
	stg storage.StorageI
	cfg config.Config
	ecom.UnimplementedAuthServiceServer
}

// NewAuthService ...
func NewAuthService(cfg config.Config, stg storage.StorageI) *authService {
	return &authService{
		cfg: cfg,
		stg: stg,
	}
}

func (s *authService) CreateUser(ctx context.Context, req *ecom.CreateUserRequest) (*ecom.User, error) {
	id := uuid.New()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword

	err = s.stg.CreateUser(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateUser: %s", err.Error())
	}

	user, err := s.stg.GetUserById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetUserById: %s", err.Error())
	}

	return user, nil
}

func (s *authService) GetUserList(ctx context.Context, req *ecom.GetUserListRequest) (*ecom.GetUserListResponse, error) {
	res, err := s.stg.GetUserList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetUserList: %s", err.Error())
	}

	return res, nil
}

func (s *authService) UpdateUser(ctx context.Context, req *ecom.UpdateUserRequest) (*ecom.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "util.HashPassword: %s", err.Error())
	}

	req.Password = hashedPassword

	err = s.stg.UpdateUser(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateUser: %s", err.Error())
	}

	user, err := s.stg.GetUserById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetUserById: %s", err.Error())
	}

	return user, nil
}

func (s *authService) DeleteUser(ctx context.Context, req *ecom.DeleteUserRequest) (*ecom.User, error) {
	user, err := s.stg.GetUserById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetUserById: %s", err.Error())
	}

	err = s.stg.DeleteUser(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteUser: %s", err.Error())
	}

	return user, nil
}

func (s *authService) GetUserByID(ctx context.Context, req *ecom.GetUserByIDRequest) (*ecom.User, error) {
	user, err := s.stg.GetUserById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetUserById: %s", err.Error())
	}
	return user, nil
}

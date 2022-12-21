package storage

import ecom "github.com/uacademy/e_commerce/auth_service/proto-gen/e_commerce"

type StorageI interface {
	CreateUser(id string, input *ecom.CreateUserRequest) error
	GetUserList(offset, limit int, search string) (resp *ecom.GetUserListResponse, err error)
	UpdateUser(input *ecom.UpdateUserRequest) error
	GetUserById(id string) (resp *ecom.User, err error)
	DeleteUser(id string) error
	GetUserByUsername(username string) (*ecom.User, error)
}

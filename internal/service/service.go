package service

import (
	"context"
	"errors"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewUserServiceService)

type UserHandler interface {
	CreateUser(ctx context.Context, email string, password string) (uint64, error)
	VerifyCode(ctx context.Context, email string, code string) bool
	SendVerifyCode(ctx context.Context, email string) (string, error)
	GetUserInfoByID(ctx context.Context, userID uint64) (model.User, error)
	UpdateUserInfo(ctx context.Context, user model.User) error
	CheckEmailExist(ctx context.Context, email string) bool
	GetUserInfoByEmail(ctx context.Context, email string) (model.User, error)
	Logout(ctx context.Context, userID uint64) error
	DeleteUser(ctx context.Context, userID uint64) error
}

var (
	ErrPasswordsDoNotMatch = errors.New("passwords do not equal confirm password")
	ErrPasswordNotValid    = errors.New("password is invalid")
	ErrEmailVerifyCode     = errors.New("email verify code is not valid")
	ErrCreateUser          = errors.New("create user failed")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrSendVerifyCode      = errors.New("send verify code failed")
	ErrUserNotFound        = errors.New("user not found")
	ErrGetUserInfo         = errors.New("get user info failed")
	ErrUpdateUser          = errors.New("update user info failed")
	ErrLogin               = errors.New("login failed")
	ErrPasswordIncorrect   = errors.New("password incorrect")
	ErrLogout              = errors.New("logout failed")
	ErrDeleteUser          = errors.New("delete user failed")
	ErrEmailExist          = errors.New("email already exists")
)

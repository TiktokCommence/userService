package biz

import (
	"context"
	"errors"
	"fmt"
	"github.com/TiktokCommence/userService/internal/errcode"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/TiktokCommence/userService/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var _ service.UserHandler = (*UserHandler)(nil)

const (
	InvalidID uint64 = 0
)

type GenerateID interface {
	GenerateUserID(ctx context.Context) (uint64, error)
}

type EmailWorker interface {
	VerifyEmailCode(ctx context.Context, email, code string) bool
	SendEmailCode(ctx context.Context, email string) (string, error)
}

type RedisWorker interface {
	GetUserByID(ctx context.Context, id uint64) (model.User, error)
	SetNULLUser(ctx context.Context, id uint64) error
	SetUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id uint64) error
	EnableRead(ctx context.Context, id uint64) error
	DisableRead(ctx context.Context, id uint64) error
}
type DBWorker interface {
	CreateUser(ctx context.Context, user model.User) error
	GetUserByID(ctx context.Context, id uint64) (model.User, error)
	SaveUserInfo(ctx context.Context, user model.User) error
	CheckEmailExist(ctx context.Context, email string) bool
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
	DeleteUser(ctx context.Context, id uint64) error
}

type UserHandler struct {
	g GenerateID
	r RedisWorker
	d DBWorker
	e EmailWorker
	h *log.Helper
}

func NewUserHandler(g GenerateID, r RedisWorker, d DBWorker, e EmailWorker, logger log.Logger) *UserHandler {
	return &UserHandler{
		g: g,
		r: r,
		d: d,
		e: e,
		h: log.NewHelper(logger),
	}
}

func (u *UserHandler) CreateUser(ctx context.Context, email string, password string) (uint64, error) {
	id, err := u.g.GenerateUserID(ctx)
	if err != nil {
		return InvalidID, err
	}
	user := model.User{
		ID:       id,
		Email:    email,
		Password: password,
	}
	err = u.d.CreateUser(ctx, user)
	if err != nil {
		return InvalidID, err
	}
	return id, nil
}

func (u *UserHandler) VerifyCode(ctx context.Context, email string, code string) bool {
	return u.e.VerifyEmailCode(ctx, email, code)
}

func (u *UserHandler) SendVerifyCode(ctx context.Context, email string) (string, error) {

	return u.e.SendEmailCode(ctx, email)
}

func (u *UserHandler) GetUserInfoByID(ctx context.Context, userID uint64) (model.User, error) {
	user, err := u.r.GetUserByID(ctx, userID)
	if err != nil && !errors.Is(err, errcode.CacheMiss) && !errors.Is(err, errcode.CacheNullValue) {
		return model.User{}, err
	}
	if errors.Is(err, errcode.CacheNullValue) {
		return model.User{}, nil
	}
	if errors.Is(err, errcode.CacheMiss) {
		user, err = u.d.GetUserByID(ctx, userID)
		if err != nil && !errors.Is(err, errcode.UserNotFound) {
			return model.User{}, fmt.Errorf("get user %d from db err:%w", userID, err)
		}
		if errors.Is(err, errcode.UserNotFound) {
			err1 := u.r.SetNULLUser(ctx, userID)
			return model.User{}, fmt.Errorf("user %d not found and set null user err:%w", userID, err1)
		}
		err = u.r.SetUser(ctx, user)
		if err != nil {
			return user, fmt.Errorf("get user %d from db success but set user err:%w", userID, err)
		}
		return user, nil
	}
	return user, nil
}

func (u *UserHandler) UpdateUserInfo(ctx context.Context, user model.User) error {
	id := user.ID
	// 1 针对 key 维度禁用读流程写缓存机制
	if err := u.r.DisableRead(ctx, id); err != nil {
		return fmt.Errorf("disable cache read userID %d error:%w", id, err)
	}
	defer func() {
		go func() {
			tctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err1 := u.r.EnableRead(tctx, id)
			if err1 != nil {
				u.h.Warnf("enable cache read userID %d error:%v", id, err1)
			}
		}()
	}()
	// 2 删除 key 维度对应缓存
	if err := u.r.DeleteUser(ctx, id); err != nil {
		return fmt.Errorf("delete user %d error:%w", id, err)
	}
	// 3 数据写入 db
	err := u.d.SaveUserInfo(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserHandler) CheckEmailExist(ctx context.Context, email string) bool {
	return u.d.CheckEmailExist(ctx, email)
}

func (u *UserHandler) GetUserInfoByEmail(ctx context.Context, email string) (model.User, error) {
	return u.d.GetUserByEmail(ctx, email)
}
func (u *UserHandler) Logout(ctx context.Context, userID uint64) error {
	return u.r.DeleteUser(ctx, userID)
}

func (u *UserHandler) DeleteUser(ctx context.Context, userID uint64) error {
	err := u.r.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user %d from cache failed:%w", userID, err)
	}
	err = u.d.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("delete user %d from cache success but delete user in db failed:%w", userID, err)
	}
	return nil
}

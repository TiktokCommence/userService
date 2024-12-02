package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/TiktokCommence/userService/internal/biz"
	"github.com/TiktokCommence/userService/internal/errcode"
	"github.com/TiktokCommence/userService/internal/foundation/DB"
	"github.com/TiktokCommence/userService/internal/foundation/common"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
)

var _ biz.DBWorker = (*UserRepo)(nil)

type UserRepo struct {
	d common.DB
	h *log.Helper
}

func NewUserRepo(d common.DB, logger log.Logger) *UserRepo {
	return &UserRepo{d: d, h: log.NewHelper(logger)}
}

func (D *UserRepo) CreateUser(ctx context.Context, user model.User) error {
	err := D.d.Put(ctx, &user)
	defer func() {
		if err != nil {
			D.h.Errorf("put user{%v} to db error {%v}", user, err)
		}
	}()
	if errors.Is(err, DB.ErrorDBLocateTable) {
		return fmt.Errorf("create user error:%w", err)
	}
	if errors.Is(err, DB.ErrorDBDuplicateEntry) {
		return errcode.UserAlreadyExists
	}
	return err
}

func (D *UserRepo) GetUserByID(ctx context.Context, id uint64) (model.User, error) {
	user := model.User{ID: id}
	err := D.d.Query(ctx, &user, map[string]interface{}{
		"id": id,
	})
	if errors.Is(err, DB.ErrorDBMiss) {
		return user, errcode.UserNotFound
	}
	return user, err
}

func (D *UserRepo) SaveUserInfo(ctx context.Context, user model.User) error {
	err := D.d.Update(ctx, &user)
	defer func() {
		if err != nil {
			D.h.Errorf("update user{%v} to db error {%v}", user, err)
		}
	}()
	return err
}

func (D *UserRepo) CheckEmailExist(ctx context.Context, email string) bool {
	ok, _ := D.d.Exist(ctx, &model.User{}, map[string]interface{}{
		"email": email,
	})
	return ok
}

func (D *UserRepo) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	err := D.d.Query(ctx, &user, map[string]interface{}{
		"email": email,
	})
	if errors.Is(err, DB.ErrorDBMiss) {
		return user, errcode.UserNotFound
	}
	return user, err
}

func (D *UserRepo) DeleteUser(ctx context.Context, id uint64) error {
	err := D.d.Delete(ctx, &model.User{}, map[string]interface{}{
		"id": id,
	})
	defer func() {
		if err != nil {
			D.h.Errorf("delete user %d from db error {%v}", id, err)
		}
	}()
	return err
}

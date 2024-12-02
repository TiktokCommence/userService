package data

import (
	"context"
	"errors"
	"github.com/TiktokCommence/userService/internal/conf"
	"github.com/TiktokCommence/userService/internal/errcode"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"testing"
)

func initUserRepo() (*UserRepo, error) {
	db, err := NewDB(&conf.Data{
		Database: &conf.Data_Database{
			Source: "root:12345678@tcp(127.0.0.1:13306)/user?parseTime=True&loc=Local",
		}})
	if err != nil {
		return nil, err
	}
	logger := log.NewStdLogger(os.Stdout)
	return NewUserRepo(db, logger), nil
}

func TestUserRepo_CreateUser(t *testing.T) {
	ctx := context.Background()
	userRepo, err := initUserRepo()
	if err != nil {
		t.Fatal(err)
	}

	user := model.User{
		ID:       uint64(11111),
		Email:    "test@example1.com",
		Password: "123456",
	}
	err = userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	user = model.User{
		ID:       uint64(11112),
		Email:    "test@example2.com",
		Password: "123456",
	}
	err = userRepo.CreateUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserRepo_GetUserByID(t *testing.T) {
	ctx := context.Background()
	userRepo, err := initUserRepo()
	if err != nil {
		t.Fatal(err)
	}
	user, err := userRepo.GetUserByID(ctx, 11112)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
	user, err = userRepo.GetUserByID(ctx, 11113)
	if errors.Is(err, errcode.UserNotFound) {
		t.Log("User not found")
	}

}

func TestUserRepo_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	userRepo, err := initUserRepo()
	if err != nil {
		t.Fatal(err)
	}
	user, err := userRepo.GetUserByEmail(ctx, "test@example2.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestUserRepo_DeleteUser(t *testing.T) {
	ctx := context.Background()
	userRepo, err := initUserRepo()
	if err != nil {
		t.Fatal(err)
	}
	err = userRepo.DeleteUser(ctx, 11112)
	if err != nil {
		t.Fatal(err)
	}
}

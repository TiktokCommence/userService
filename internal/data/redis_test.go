package data

import (
	"context"
	"github.com/TiktokCommence/userService/internal/conf"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"testing"
)

func initRedisWorkerImplement() *RedisWorkerImplement {
	cache := NewCache(&conf.Data{Redis: &conf.Data_Redis{
		Addr:        "127.0.0.1:16379",
		Password:    "",
		MaxIdle:     10,
		IdleTimeout: 2,
		MaxActive:   15,
		Wait:        true,
	}})
	options := NewOptions(&conf.Data{Redis: &conf.Data_Redis{ExpirationSeconds: 300}})
	logger := log.NewStdLogger(os.Stdout)
	ri := NewRedisWorkerImplement(cache, options, logger)
	return ri
}

func TestRedisWorkerImplement_DeleteUser(t *testing.T) {
	ri := initRedisWorkerImplement()
	err := ri.SetUser(context.Background(), model.User{ID: 1112, Password: "1231231", Email: "296313@qq.com"})
	if err != nil {
		t.Fatal(err)
	}
	err = ri.DeleteUser(context.Background(), 1112)
	if err != nil {
		t.Fatal(err)
	}
	err = ri.DeleteUser(context.Background(), 1112)
	if err != nil {
		t.Error(err)
	}
}

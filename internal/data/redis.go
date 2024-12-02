package data

import (
	"context"
	"errors"
	"github.com/TiktokCommence/userService/internal/biz"
	"github.com/TiktokCommence/userService/internal/errcode"
	cache2 "github.com/TiktokCommence/userService/internal/foundation/cache"
	"github.com/TiktokCommence/userService/internal/foundation/common"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gomodule/redigo/redis"
	"sync"
)

var _ biz.RedisWorker = (*RedisWorkerImplement)(nil)

const NullData = "Err_Syntax_Null_Data"

const StartValue int64 = 0

const IDKey = "id"

type RedisWorkerImplement struct {
	c    common.Cache
	opt  *cache2.Options
	once sync.Once
	h    *log.Helper
}

func (r *RedisWorkerImplement) GenerateUserID(ctx context.Context) (uint64, error) {
	r.once.Do(func() {
		r.c.Set(ctx, IDKey, StartValue)
	})
	id, err := redis.Uint64(r.c.IncrBy(ctx, IDKey, 1))
	if err != nil {
		return 0, err
	}
	return id, nil
}

func NewRedisWorkerImplement(c common.Cache, opt *cache2.Options, logger log.Logger) *RedisWorkerImplement {
	return &RedisWorkerImplement{c: c, opt: opt, h: log.NewHelper(logger)}
}

func (r *RedisWorkerImplement) EnableRead(ctx context.Context, id uint64) error {
	key := GenerateKey(id)
	return r.c.Enable(ctx, key, r.opt.EnableDelayMills)
}

func (r *RedisWorkerImplement) DisableRead(ctx context.Context, id uint64) error {
	key := GenerateKey(id)
	return r.c.Disable(ctx, key, r.opt.DisableExpireSeconds)
}

func (r *RedisWorkerImplement) GetUserByID(ctx context.Context, id uint64) (model.User, error) {
	key := GenerateKey(id)
	val, err := r.c.Get(ctx, key)
	if errors.Is(err, cache2.ErrorCacheMiss) {
		return model.User{}, errcode.CacheMiss
	}
	if err != nil {
		return model.User{}, err
	}
	if val == NullData {
		return model.User{}, errcode.CacheNullValue
	}

	user := model.User{}
	err = user.Read(val)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *RedisWorkerImplement) SetNULLUser(ctx context.Context, id uint64) error {
	key := GenerateKey(id)
	ok, err := r.c.PutWhenEnable(ctx, key, NullData, r.opt.GetCacheExpireSeconds())
	if err != nil {
		r.h.Errorf("put null data into cache fail, key: %s, err: %v", key, err)
		return err
	}
	r.h.Infof("put null data into cache resp, key: %s, ok: %t", key, ok)
	return nil
}

func (r *RedisWorkerImplement) SetUser(ctx context.Context, user model.User) error {
	key := GenerateKey(user.ID)
	val, err := user.Write()
	if err != nil {
		return err
	}
	ok, err := r.c.PutWhenEnable(ctx, key, val, r.opt.GetCacheExpireSeconds())
	if err != nil {
		r.h.Errorf("put data into cache fail, key: %s, data: %v, err: %v", key, val, err)
		return err
	}
	r.h.Infof("put data into cache resp, key: %s, v: %v, ok: %t", key, val, ok)
	return nil
}

func (r *RedisWorkerImplement) DeleteUser(ctx context.Context, id uint64) error {
	key := GenerateKey(id)
	return r.c.Del(ctx, key)
}

package data

import (
	"fmt"
	"github.com/TiktokCommence/userService/internal/conf"
	DB2 "github.com/TiktokCommence/userService/internal/foundation/DB"
	cache2 "github.com/TiktokCommence/userService/internal/foundation/cache"
	"github.com/TiktokCommence/userService/internal/foundation/common"
	"github.com/TiktokCommence/userService/internal/model"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewCache, NewOptions, NewUserRepo, NewEmailWorker, NewRedisWorkerImplement)

func NewDB(data *conf.Data) (common.DB, error) {
	tables := []interface{}{&model.User{}}
	return DB2.NewDB(&DB2.Config{Tables: tables, Dsn: data.Database.Source}, DB2.WithDuplicateEntry(false))
}
func NewCache(c *conf.Data) common.Cache {
	cf := &cache2.Config{
		Address:            c.Redis.Addr,
		Password:           c.Redis.Password,
		MaxIdle:            int(c.Redis.MaxIdle),
		IdleTimeoutSeconds: int(c.Redis.IdleTimeout),
		MaxActive:          int(c.Redis.MaxActive),
		Wait:               c.Redis.Wait,
	}
	client := cache2.NewRClient(cf)
	cac := cache2.NewCache(client)
	return cac
}

func NewOptions(c *conf.Data) *cache2.Options {
	options := cache2.NewOptions(
		cache2.WithCacheExpireSeconds(c.Redis.ExpirationSeconds),
		cache2.WithCacheExpireRandomMode(),
		cache2.WithDisableExpireSeconds(2),
		cache2.WithEnableDelayMilis(150),
	)
	return options
}

func GenerateKey(id uint64) string {
	return fmt.Sprintf("user:%d", id)
}

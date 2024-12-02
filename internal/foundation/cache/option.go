package cache

import (
	"math/rand"
	"time"
)

type Options struct {
	// 缓存过期时间，单位：秒
	CacheExpireSeconds int64
	// 是否启用过期时间扰动
	CacheExpireRandomMode bool
	// 禁用读流程写缓存模式过期时间，单位：秒
	DisableExpireSeconds int64
	// 写流程 disable 操作后延时多长时间进行 enable 操作，单位：毫秒
	EnableDelayMills int64
	// 随机数生成器
	rander *rand.Rand
}

func (o *Options) GetCacheExpireSeconds() int64 {
	if !o.CacheExpireRandomMode {
		return o.CacheExpireSeconds
	}

	// 过期时间在 1~2倍之间取随机值
	return o.CacheExpireSeconds + o.rander.Int63n(o.CacheExpireSeconds+1)
}

type Option func(*Options)

const (
	// 默认的缓存过期时间为 60 s
	DefaultCacheExpireSeconds = 60
	// 默认的禁用写缓存时间为 10 s
	DefaultDisableExpireSeconds = 10
	// 默认的延时 enable 时间为 1 s
	DefaultEnableDelayMilis = 1000
)

func NewOptions(opts ...Option) *Options {
	options := &Options{
		CacheExpireSeconds:    DefaultCacheExpireSeconds,
		CacheExpireRandomMode: true,
		DisableExpireSeconds:  DefaultDisableExpireSeconds,
		EnableDelayMills:      DefaultEnableDelayMilis,
	}
	for _, opt := range opts {
		opt(options)
	}
	repair(options)
	return options
}

func WithCacheExpireSeconds(cacheExpireSeconds int64) Option {
	return func(o *Options) {
		o.CacheExpireSeconds = cacheExpireSeconds
	}
}
func WithCacheExpireRandomMode() Option {
	return func(o *Options) {
		o.CacheExpireRandomMode = true
		o.rander = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}
func WithDisableExpireSeconds(disableExpireSeconds int64) Option {
	return func(o *Options) {
		o.DisableExpireSeconds = disableExpireSeconds
	}
}
func WithEnableDelayMilis(enableDelayMilis int64) Option {
	return func(o *Options) {
		o.EnableDelayMills = enableDelayMilis
	}
}

func repair(o *Options) {
	if o.CacheExpireSeconds <= 0 {
		o.CacheExpireSeconds = DefaultCacheExpireSeconds
	}

	if o.DisableExpireSeconds <= 0 {
		o.DisableExpireSeconds = DefaultDisableExpireSeconds
	}

	if o.EnableDelayMills <= 0 {
		o.EnableDelayMills = DefaultEnableDelayMilis
	}
}

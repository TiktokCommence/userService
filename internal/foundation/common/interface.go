package common

import "context"

//将基础模块抽象

// 缓存模块的抽象接口定义
type Cache interface {
	// 启用某个 key 对应读流程写缓存机制（默认情况下为启用状态）
	Enable(ctx context.Context, key string, delayMilis int64) error
	// 禁用某个 key 对应读流程写缓存机制
	Disable(ctx context.Context, key string, expireSeconds int64) error
	// 读取 key 对应缓存
	Get(ctx context.Context, key string) (string, error)
	// 删除 key 对应缓存
	Del(ctx context.Context, key string) error
	// 校验某个 key 对应读流程写缓存机制是否启用，倘若启用则写入缓存（默认情况下为启用状态）
	PutWhenEnable(ctx context.Context, key, value string, expireSeconds int64) (bool, error)

	Set(ctx context.Context, key string, value interface{}) error

	IncrBy(ctx context.Context, key string, step int64) (int64, error)
	SetEx(ctx context.Context, key, value string, expireSeconds int64) error
}

// 数据库模块的抽象接口定义
type DB interface {
	// 数据写入数据库
	Put(ctx context.Context, obj Object) error
	// 从数据库读取数据(通过查询条件)
	Query(ctx context.Context, obj Object, params map[string]interface{}) error
	// 删除
	Delete(ctx context.Context, obj Object, params map[string]interface{}) error

	Update(ctx context.Context, obj Object) error

	Exist(ctx context.Context, obj Object, params map[string]interface{}) (bool, error)
}

// 每次读写操作时，操作的一笔数据记录
type Object interface {
	// 获取 key 对应的字段名
	KeyColumn() string
	// 获取 key 对应的值
	Key() interface{}

	// 将 object 序列化成字符串
	Write() (string, error)
	// 读取字符串内容，反序列化到 object 实例中
	Read(body string) error
}

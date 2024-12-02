package DB

import (
	"context"
	"errors"
	"fmt"
	"github.com/TiktokCommence/userService/internal/foundation/common"
	"gorm.io/gorm"
)

var (
	ErrorDBMiss           = errors.New("DB miss")
	ErrorDBLocateTable    = errors.New("the obj don't implement TableName method")
	ErrorDBDuplicateEntry = errors.New("DB duplicate entry")
	ErrorDBUpdate         = errors.New("DB update failed")
)

type tabler interface {
	TableName() string
}

// 数据库模块的抽象接口定义
type DB struct {
	db  *gorm.DB
	opt Options
}
type Config struct {
	Dsn    string
	Tables []interface{}
}

func NewDB(c *Config, opts ...Option) (*DB, error) {
	defaultOpts := Options{
		DuplicateEntry: true,
	}
	for _, opt := range opts {
		opt(&defaultOpts)
	}
	db, err := getDB(c.Dsn, c.Tables)
	if err != nil {
		return nil, err
	}
	return &DB{db: db, opt: defaultOpts}, nil
}

// 数据写入数据库
func (d *DB) Put(ctx context.Context, obj common.Object) error {
	db := d.db
	tabler, ok := obj.(tabler)
	if !ok {
		return ErrorDBLocateTable
	}
	db = db.Table(tabler.TableName())

	// 此处通过两个非原子性动作实现 upsert 效果：
	// 1 尝试创建记录
	// 2 倘若发生唯一键冲突，则改为执行更新操作
	err := db.WithContext(ctx).Create(obj).Error
	if err == nil {
		return nil
	}

	// 判断是否为唯一键冲突
	dup := IsDuplicateEntryErr(err)
	if dup {
		err = ErrorDBDuplicateEntry
	}

	if dup && d.opt.DuplicateEntry {
		res := db.WithContext(ctx).Debug().Where(fmt.Sprintf("`%s` = ?", obj.KeyColumn()), obj.Key()).Updates(obj)
		if res.RowsAffected == 0 {
			return ErrorDBUpdate
		}
	}
	// 其他错误直接返回
	return err
}

func (d *DB) Query(ctx context.Context, obj common.Object, params map[string]interface{}) error {
	db := d.db
	tabler, ok := obj.(tabler)
	if !ok {
		return ErrorDBLocateTable
	}
	db = db.Table(tabler.TableName())
	if ok, err := d.checkParams(params); !ok {
		return err
	}
	err := db.WithContext(ctx).Where(params).First(obj).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorDBMiss
	}
	return err
}

func (d *DB) Delete(ctx context.Context, obj common.Object, params map[string]interface{}) error {
	db := d.db
	tabler, ok := obj.(tabler)
	if !ok {
		return ErrorDBLocateTable
	}
	db = db.Table(tabler.TableName())
	if ok, err := d.checkParams(params); !ok {
		return err
	}
	err := db.WithContext(ctx).Where(params).Delete(obj).Error
	return err
}

func (d *DB) Update(ctx context.Context, obj common.Object) error {
	db := d.db
	tabler, ok := obj.(tabler)
	if !ok {
		return ErrorDBLocateTable
	}
	db = db.Table(tabler.TableName())
	res := db.WithContext(ctx).Updates(obj)
	if res.RowsAffected == 0 {
		return ErrorDBUpdate
	}
	return nil
}
func (d *DB) Exist(ctx context.Context, obj common.Object, params map[string]interface{}) (bool, error) {
	db := d.db
	tabler, ok := obj.(tabler)
	if !ok {
		return false, ErrorDBLocateTable
	}
	db = db.Table(tabler.TableName())
	if ok, err := d.checkParams(params); !ok {
		return false, err
	}
	var cnt int64
	err := db.WithContext(ctx).Where(params).Count(&cnt).Error
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (d *DB) checkParams(params map[string]interface{}) (bool, error) {
	if params == nil {
		return false, errors.New("the map is nil and considered empty")
	}
	if len(params) == 0 {
		return false, errors.New("the map is empty")
	}
	for _, v := range params {
		if v == nil {
			return false, errors.New("the map contains a nil value")
		}
	}
	return true, nil
}

package DB

import (
	"errors"
	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 唯一键冲突错误
const DuplicateEntryErrCode = 1062

type Options struct {
	// 是否开启唯一键冲突后更新的操作
	DuplicateEntry bool
}

type Option func(*Options)

func WithDuplicateEntry(duplicateEntry bool) Option {
	return func(o *Options) {
		o.DuplicateEntry = duplicateEntry
	}
}

// GetClient 获取一个数据库客户端
func getDB(dsn string, tables []interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("connect mysql failed")
	}
	for _, table := range tables {
		if err = db.AutoMigrate(table); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// 是否为唯一键冲突错误
func IsDuplicateEntryErr(err error) bool {
	var mysqlErr *mysql2.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == DuplicateEntryErrCode {
		return true
	}
	return false
}

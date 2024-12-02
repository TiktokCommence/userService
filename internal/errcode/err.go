package errcode

import "errors"

var (
	UserAlreadyExists = errors.New("user already exists")
	UserNotFound      = errors.New("user not found in db")
	CacheMiss         = errors.New("cache miss")
	CacheNullValue    = errors.New("cache null value")
)

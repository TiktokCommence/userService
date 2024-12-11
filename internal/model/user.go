package model

import (
	"encoding/json"
	"time"
)

const (
	UserTableName = "users"
)

type User struct {
	ID        uint64  `gorm:"primaryKey;column:id"`
	Password  string  `gorm:"column:password"`
	Name      *string `gorm:"column:username;type:varchar(200)"`
	Email     string  `gorm:"column:email;type:varchar(100);uniqueIndex"`
	Age       *int32  `gorm:"column:age"`
	Addr1     *string `gorm:"column:addr1;type:varchar(100)"`
	Addr2     *string `gorm:"column:addr2;type:varchar(100)"`
	Phone     *string `gorm:"column:phone;type:varchar(30)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) KeyColumn() string {
	return "id"
}

func (u *User) Key() interface{} {
	return u.ID
}

func (u *User) Write() (string, error) {
	body, err := json.Marshal(u)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (u *User) Read(body string) error {
	return json.Unmarshal([]byte(body), u)
}

func (u *User) TableName() string {
	return UserTableName
}

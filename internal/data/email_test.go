package data

import (
	"context"
	"github.com/TiktokCommence/userService/internal/conf"
	cache3 "github.com/TiktokCommence/userService/internal/foundation/cache"
	"testing"
)

var exampleEmail = "example@qq.com"

func initEmailWorker() *EmailWorker {
	client := cache3.NewRClient(&cache3.Config{
		Address:            "localhost:16379",
		Password:           "",
		MaxIdle:            10,
		IdleTimeoutSeconds: 2,
		MaxActive:          15,
		Wait:               true,
	})
	cache := cache3.NewCache(client)
	ew := NewEmailWorker(cache, &conf.EmailConf{
		Sender:            exampleEmail,
		Secret:            "gmoyxtvrxqsfdhca",
		ExpirationSeconds: 5 * 60,
	})
	return ew
}

func TestEmailWorker_SendEmailCode(t *testing.T) {
	ew := initEmailWorker()
	code, err := ew.SendEmailCode(context.Background(), exampleEmail)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)
}
func TestEmailWorker_VerifyEmailCode(t *testing.T) {
	ew := initEmailWorker()
	ok := ew.VerifyEmailCode(context.Background(), exampleEmail, "4DH0")
	if ok {
		t.Log("success\n")
	} else {
		t.Log("failed\n")
	}
}

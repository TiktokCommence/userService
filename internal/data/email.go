package data

import (
	"context"
	"fmt"
	"github.com/TiktokCommence/userService/internal/biz"
	"github.com/TiktokCommence/userService/internal/conf"
	"github.com/TiktokCommence/userService/internal/foundation/common"
	email2 "github.com/jordan-wright/email"
	"math/rand"
	"net/smtp"
	"time"
)

var _ biz.EmailWorker = (*EmailWorker)(nil)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type EmailWorker struct {
	c  common.Cache
	cf *conf.EmailConf
}

func NewEmailWorker(c common.Cache, cf *conf.EmailConf) *EmailWorker {
	return &EmailWorker{c: c, cf: cf}
}

func (e *EmailWorker) VerifyEmailCode(ctx context.Context, email, code string) bool {
	value, err := e.c.Get(ctx, e.generateKey(email))
	if err != nil {
		return false
	}
	if value == code {
		defer func() {
			e.c.Del(ctx, e.generateKey(email))
		}()
		return true
	}
	return false
}

func (e *EmailWorker) SendEmailCode(ctx context.Context, email string) (string, error) {
	em := email2.NewEmail()
	em.From = e.cf.Sender
	em.To = []string{email}
	code := e.generateCode()
	minutes := fmt.Sprintf("%d", e.cf.ExpirationSeconds/60)
	// 设置邮件的HTML内容
	em.HTML = []byte(`
		<h1>Verification Code</h1>
		<p>你的验证码是: <strong>` + code + `</strong>,该验证码将在` + minutes + `分钟后失效</p>
	`)
	em.Send("smtp.qq.com:587", smtp.PlainAuth("", e.cf.Sender, e.cf.Secret, "smtp.qq.com"))
	err := e.c.SetEx(ctx, e.generateKey(email), code, e.cf.ExpirationSeconds)
	if err != nil {
		return "", err
	}
	return code, nil
}
func (e *EmailWorker) generateCode() string {
	rand.Seed(time.Now().UnixNano())
	// 四位大写英文字母与数字混合验证码
	var result []byte
	for i := 0; i < 4; i++ {
		result = append(result, charset[rand.Intn(len(charset))])
	}
	return string(result)
}
func (e *EmailWorker) generateKey(email string) string {
	return fmt.Sprintf("email:%s", email)
}

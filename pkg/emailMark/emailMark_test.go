package emailMark

import (
	"im/utils"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/XYYSWK/Lutils/pkg/email"
)

func TestSendCode(t *testing.T) {
	mark := New(Config{
		UserMarkDuration: 10 * time.Second,
		CodeMarkDuration: 10 * time.Second,
		SMTPInfo: email.SMTPInfo{
			Port:     465,
			IsSSL:    true,
			Host:     "smtp.qq.com",
			UserName: "2221866404@qq.com",
			Password: "pgytthqtrgmndjcg",
			From:     "2221866404@qq.com",
		},
		AppName: "chat",
	})

	code := utils.RandomString(6)
	emailStr := "2221866404@qq.com"
	// 测试发送
	log.Println("Send01:")
	_ = mark.SendMark(emailStr, code)
	//require.Error(t, err)
	log.Println("Check01:")
	require.True(t, mark.CheckCode(emailStr, code))

	//// 测试频繁请求，不会发送邮件
	//log.Println("Send02:")
	//require.ErrorIs(t, mark.SendMark(emailStr, code), ErrSendTooMany)

	// 测试用户时间间隔后再次请求验证码
	//time.Sleep(mark.config.UserMarkDuration)
	//code = utils.RandomString(6)
	//log.Println("Send04:")
	//require.NoError(t, mark.SendMark(emailStr, code))
	//log.Println("Check04:")
	//require.True(t, mark.CheckCode(emailStr, code))

	// 测试验证码过期
	<-time.After(mark.config.CodeMarkDuration) // 等待一段时间模拟验证码过期
	log.Println("Send05:")
	require.False(t, mark.CheckCode(emailStr, code))
}

package emailMark

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/XYYSWK/Lutils/pkg/email"
)

/*
邮箱验证码标记
*/

type EmailMark struct {
	config   Config
	userMark sync.Map //标记用户
	codeMark sync.Map //记录 code
}

type Config struct {
	UserMarkDuration time.Duration  // 用户标记时长
	CodeMarkDuration time.Duration  // 验证码标记时长
	SMTPInfo         email.SMTPInfo //邮箱配置
	AppName          string         //应用名称
}

func New(config Config) *EmailMark {
	return &EmailMark{
		config:   config,
		userMark: sync.Map{},
		codeMark: sync.Map{},
	}
}

var ErrSendTooMany = errors.New("发送过于频繁")

// CheckUserExist 判断邮箱是否已经被记录
func (m *EmailMark) CheckUserExist(emailStr string) bool {
	//判断该键是否存在
	_, ok := m.userMark.Load(emailStr)
	return ok
}

// SendMark 发送验证码
func (m *EmailMark) SendMark(emailStr, code string) error {
	// 发送频率限制
	if m.CheckUserExist(emailStr) {
		return ErrSendTooMany
	}
	// 记录邮箱
	m.userMark.Store(emailStr, struct{}{})
	sendMark := email.NewEmail(&m.config.SMTPInfo)
	// 发送邮件
	err := sendMark.SendMail([]string{emailStr}, fmt.Sprintf("%s邮箱验证码", m.config.AppName), fmt.Sprintf(m.getHtml(m.config.AppName, code, int(m.config.CodeMarkDuration.Minutes()))))
	if err != nil {
		// 发送失败删除标记
		m.userMark.Delete(emailStr)
		return err
	}
	// 记录验证码
	m.codeMark.Store(emailStr, code)
	// 验证码过期或用户一定间隔时间后被允许再次请求验证码，延时删除标记
	m.DeleteMarkDelay(emailStr)
	return nil
}

// DeleteMarkDelay 验证码过期，延时删除标记
func (m *EmailMark) DeleteMarkDelay(emailStr string) {
	time.AfterFunc(m.config.CodeMarkDuration, func() {
		m.codeMark.Delete(emailStr)
	})
	time.AfterFunc(m.config.UserMarkDuration, func() {
		m.userMark.Delete(emailStr)
	})
}

// CheckCode 校验验证码
func (m *EmailMark) CheckCode(emailStr, code string) bool {
	// 获取验证码
	myCode, ok := m.codeMark.Load(emailStr)
	// 验证成功删除标记
	if ok && myCode == code {
		m.codeMark.Delete(emailStr)
		return true
	}
	return false
}

func (*EmailMark) getHtml(appName, code string, expireMinutes int) string {
	// 构建HTML邮件内容
	htmlContent := "<!DOCTYPE html>" +
		"<html lang=\"zh-CN\">" +
		"<head>" +
		"<meta charset=\"UTF-8\">" +
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">" +
		"<title>验证码验证</title>" +
		"<style>" +
		"    body { margin: 0; padding: 0; background-color: #f7f9fc; font-family: 'Microsoft YaHei', sans-serif; }" +
		"    .email-container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 12px; overflow: hidden; box-shadow: 0 4px 20px rgba(0,0,0,0.08); }" +
		"    .email-header { background: linear-gradient(135deg, #4a90e2 0%, #5c6bc0 100%); padding: 30px 20px; text-align: center; }" +
		"    .header-title { color: #ffffff; margin: 0; font-size: 24px; font-weight: 600; }" +
		"    .email-body { padding: 30px 40px; }" +
		"    .greeting { color: #333333; font-size: 18px; margin-bottom: 20px; }" +
		"    .content-text { color: #666666; line-height: 1.8; font-size: 16px; margin-bottom: 30px; }" +
		"    .code-container { background-color: #f8f9fa; border-radius: 8px; padding: 25px 20px; text-align: center; margin-bottom: 30px; }" +
		"    .verification-code { font-size: 32px; font-weight: bold; color: #4a90e2; letter-spacing: 8px; margin: 0; padding: 10px 0; }" +
		"    .expire-info { color: #999999; font-size: 14px; margin-top: 15px; }" +
		"    .note-text { color: #ff5252; font-size: 14px; line-height: 1.6; padding: 15px; background-color: #fff8f8; border-left: 4px solid #ff5252; border-radius: 4px; margin-bottom: 30px; }" +
		"    .email-footer { padding: 20px 40px; background-color: #f8f9fa; border-top: 1px solid #eee; }" +
		"    .footer-text { color: #999999; font-size: 12px; text-align: center; margin: 0; line-height: 1.6; }" +
		"    @media (max-width: 600px) {" +
		"        .email-container { border-radius: 0; }" +
		"        .email-body { padding: 20px 15px; }" +
		"        .verification-code { font-size: 26px; letter-spacing: 5px; }" +
		"    }" +
		"</style>" +
		"</head>" +
		"<body>" +
		"    <div class=\"email-container\">" +
		"        <div class=\"email-header\">" +
		"            <h1 class=\"header-title\">" + appName + "安全验证</h1>" +
		"        </div>" +
		"        <div class=\"email-body\">" +
		"            <p class=\"greeting\">您好！</p>" +
		"            <p class=\"content-text\">" +
		"                您正在进行身份验证操作，以下是您的验证码。请在验证页面输入此验证码，完成操作。" +
		"            </p>" +
		"            <div class=\"code-container\">" +
		"                <p class=\"verification-code\">" + code + "</p>" +
		"                <p class=\"expire-info\">验证码 " + strconv.Itoa(expireMinutes) + " 分钟内有效，请勿向他人泄露</p>" +
		"            </div>" +
		"            <p class=\"note-text\">" +
		"                注意：如非本人操作，请忽略此邮件。切勿将验证码透露给他人，以保障您的账号安全。" +
		"            </p>" +
		"        </div>" +
		"        <div class=\"email-footer\">" +
		"            <p class=\"footer-text\">" +
		"                此邮件为系统自动发送，请勿直接回复<br>" +
		"                © 2025 " + appName + " 版权所有" +
		"            </p>" +
		"        </div>" +
		"    </div>" +
		"</body>" +
		"</html>"
	return htmlContent
}

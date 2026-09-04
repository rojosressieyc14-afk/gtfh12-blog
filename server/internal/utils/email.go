package utils

import (
	"fmt"
	"net/smtp"
)

func SendVerificationCode(to, code, host, port, user, password, from string) error {
	addr := host + ":" + port

	subject := "PulseBlog - 邮箱验证码"
	body := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; padding: 40px 20px; background: #fafafa;">
  <div style="max-width: 480px; margin: 0 auto; background: #fff; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,0.08); padding: 40px 32px; text-align: center;">
    <h2 style="color: #1d1d1f; margin: 0 0 8px;">PulseBlog</h2>
    <p style="color: #86868b; margin: 0 0 32px; font-size: 14px;">邮箱验证码</p>
    <div style="font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #ff8a4c; margin: 24px 0; padding: 16px; background: #fff8f0; border-radius: 8px;">%s</div>
    <p style="color: #86868b; font-size: 13px; margin: 24px 0 0;">验证码 5 分钟内有效，请勿泄露给他人。</p>
  </div>
</body>
</html>`, code)

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=utf-8\r\n"+
		"\r\n"+
		"%s", from, to, subject, body)

	auth := smtp.PlainAuth("", user, password, host)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}

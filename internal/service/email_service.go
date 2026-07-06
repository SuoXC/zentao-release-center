package service

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/yi-nology/zentao-release-center/internal/config"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config    config.EmailConfig
	enabled   bool
}

func NewEmailService(cfg config.EmailConfig) *EmailService {
	return &EmailService{
		config:  cfg,
		enabled: cfg.Enabled,
	}
}

func (s *EmailService) IsEnabled() bool {
	return s.enabled && s.config.SMTPHost != "" && s.config.Username != ""
}

func (s *EmailService) GetRecipients() []string {
	if s.config.Recipients == "" {
		return nil
	}
	var recipients []string
	for _, r := range strings.Split(s.config.Recipients, ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			recipients = append(recipients, r)
		}
	}
	return recipients
}

func (s *EmailService) Send(to, subject, htmlBody string) error {
	if !s.IsEnabled() {
		return nil
	}

	m := gomail.NewMessage(gomail.SetCharset("UTF-8"))

	if s.config.SenderName != "" {
		m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.SenderName, s.config.Username))
	} else {
		m.SetHeader("From", s.config.Username)
	}
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.Username, s.config.Password)
	d.TLSConfig = &tls.Config{ServerName: s.config.SMTPHost}
	if s.config.UseSSL {
		d.SSL = true
	}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send email: %w", err)
	}
	return nil
}

func (s *EmailService) SendToAll(subject, htmlBody string) []error {
	recipients := s.GetRecipients()
	if len(recipients) == 0 {
		return nil
	}

	var errs []error
	for _, to := range recipients {
		if err := s.Send(to, subject, htmlBody); err != nil {
			errs = append(errs, fmt.Errorf("send to %s: %w", to, err))
		}
	}
	return errs
}

func (s *EmailService) BuildReleaseHTML(release *model.Release, items []*model.ReleaseItem, projectName, version string) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; }
.container { max-width: 800px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
.header { background-color: #4F6BF6; color: #ffffff; padding: 20px 30px; }
.header h1 { margin: 0; font-size: 20px; }
.body { padding: 30px; color: #333333; line-height: 1.8; font-size: 15px; }
.footer { padding: 15px 30px; background-color: #fafafa; color: #999999; font-size: 12px; text-align: center; border-top: 1px solid #eeeeee; }
h2 { color: #1E293B; margin-top: 24px; border-bottom: 1px solid #e2e8f0; padding-bottom: 8px; }
table { border-collapse: collapse; width: 100%; margin: 16px 0; }
th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }
th { background: #f5f5f5; }
.tag { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 12px; }
.tag-bug { background: #fff1f0; color: #cf1322; }
.tag-task { background: #f6ffed; color: #389e0d; }
.tag-note { background: #e6f7ff; color: #0958d9; }
.summary-box { background: #f0f5ff; border-left: 4px solid #4F6BF6; padding: 12px 16px; margin: 16px 0; border-radius: 0 4px 4px 0; }
</style>
</head>
<body>
<div class="container">
	<div class="header"><h1>📦 发布单通知</h1></div>
	<div class="body">`)

	sb.WriteString(fmt.Sprintf(`<div class="summary-box">
		<strong>📋 %s</strong><br>
		版本: %s | 项目: %s | 时间: %s
	</div>`, release.Name, version, projectName, time.Now().Format("2006-01-02 15:04:05")))

	if release.Summary != "" {
		sb.WriteString(fmt.Sprintf("<h2>概述</h2><p>%s</p>", release.Summary))
	}

	var bugs, tasks, notes []*model.ReleaseItem
	for _, item := range items {
		switch item.ItemType {
		case "bug":
			bugs = append(bugs, item)
		case "task":
			tasks = append(tasks, item)
		case "note":
			notes = append(notes, item)
		}
	}

	sb.WriteString(fmt.Sprintf("<h2>📊 统计</h2><p>共 %d 个条目</p><ul>", len(items)))
	if len(bugs) > 0 {
		sb.WriteString(fmt.Sprintf("<li>🐛 Bug 修复: %d 个</li>", len(bugs)))
	}
	if len(tasks) > 0 {
		sb.WriteString(fmt.Sprintf("<li>✅ 任务完成: %d 个</li>", len(tasks)))
	}
	if len(notes) > 0 {
		sb.WriteString(fmt.Sprintf("<li>📝 备注: %d 个</li>", len(notes)))
	}
	sb.WriteString("</ul>")

	if len(bugs) > 0 {
		sb.WriteString(fmt.Sprintf("<h2>🐛 Bug 修复（%d）</h2>", len(bugs)))
		sb.WriteString("<table><tr><th>#</th><th>标题</th><th>严重程度</th><th>优先级</th><th>状态</th><th>指派给</th></tr>")
		for _, b := range bugs {
			sb.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				b.ZentaoID, b.Title, b.Severity, b.Priority, b.Status, b.AssignedTo))
		}
		sb.WriteString("</table>")
	}

	if len(tasks) > 0 {
		sb.WriteString(fmt.Sprintf("<h2>✅ 任务完成（%d）</h2>", len(tasks)))
		sb.WriteString("<table><tr><th>#</th><th>标题</th><th>优先级</th><th>状态</th><th>指派给</th></tr>")
		for _, t := range tasks {
			sb.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				t.ZentaoID, t.Title, t.Priority, t.Status, t.AssignedTo))
		}
		sb.WriteString("</table>")
	}

	if len(notes) > 0 {
		sb.WriteString(fmt.Sprintf("<h2>📝 备注（%d）</h2>", len(notes)))
		for _, n := range notes {
			sb.WriteString(fmt.Sprintf("<h3>%s</h3><p>%s</p><hr>", n.NoteTitle, n.NoteContent))
		}
	}

	sb.WriteString(`</div>
	<div class="footer">此邮件由 zentao-release-center 自动发送，请勿直接回复</div>
</div>
</body>
</html>`)

	return sb.String()
}

func (s *EmailService) BuildReleaseSubject(release *model.Release, projectName, version string) string {
	subject := fmt.Sprintf("[发布通知] %s", release.Name)
	if version != "" {
		subject += fmt.Sprintf(" v%s", version)
	}
	if projectName != "" {
		subject += fmt.Sprintf(" - %s", projectName)
	}
	return subject
}

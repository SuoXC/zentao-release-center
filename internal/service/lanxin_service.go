package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yi-nology/zentao-release-center/internal/config"
	"github.com/yi-nology/zentao-release-center/internal/model"
)

type LanxinService struct {
	url     string
	secret  string
	skipSSL bool
	enabled bool
}

func NewLanxinService(cfg config.LanxinConfig) *LanxinService {
	return &LanxinService{
		url:     cfg.URL,
		secret:  cfg.Secret,
		skipSSL: cfg.SkipSSL,
		enabled: cfg.Enabled,
	}
}

func (s *LanxinService) IsEnabled() bool {
	return s.enabled && s.url != ""
}

func genLanxinSign(secret string) (string, string) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	stringToSign := timestamp + "@" + secret
	h := hmac.New(sha256.New, []byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, timestamp
}

func (s *LanxinService) Send(message string) error {
	if !s.IsEnabled() {
		return nil
	}

	msg := map[string]interface{}{
		"msgType": "text",
		"msgData": map[string]interface{}{
			"text": map[string]string{
				"content": message,
			},
		},
	}

	if s.secret != "" {
		sign, timestamp := genLanxinSign(s.secret)
		msg["sign"] = sign
		msg["timestamp"] = timestamp
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal lanxin message: %w", err)
	}

	client := &http.Client{Timeout: 15 * time.Second}
	if s.skipSSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	req, err := http.NewRequest("POST", s.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err == nil {
		if errCode, ok := respBody["errCode"]; ok {
			if code, ok := errCode.(float64); ok && code == 0 {
				return nil
			}
			return fmt.Errorf("lanxin error: errCode=%v errMsg=%v", respBody["errCode"], respBody["errMsg"])
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("lanxin http error: status %d", resp.StatusCode)
	}

	return nil
}

func (s *LanxinService) BuildReleaseMessage(release *model.Release, items []*model.ReleaseItem, projectName, version string) string {
	var msg string
	msg += "📦 发布单通知\n"
	msg += "━━━━━━━━━━━━━━━━━━━━\n"
	msg += fmt.Sprintf("📋 项目: %s\n", projectName)
	msg += fmt.Sprintf("📌 发布单: %s\n", release.Name)
	if version != "" {
		msg += fmt.Sprintf("🏷️ 版本: %s\n", version)
	}
	msg += fmt.Sprintf("📅 时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	msg += "━━━━━━━━━━━━━━━━━━━━\n\n"

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

	msg += fmt.Sprintf("📊 统计: 共 %d 个条目\n", len(items))
	if len(bugs) > 0 {
		msg += fmt.Sprintf("  🐛 Bug 修复: %d 个\n", len(bugs))
	}
	if len(tasks) > 0 {
		msg += fmt.Sprintf("  ✅ 任务完成: %d 个\n", len(tasks))
	}
	if len(notes) > 0 {
		msg += fmt.Sprintf("  📝 备注: %d 个\n", len(notes))
	}

	if len(bugs) > 0 {
		msg += "\n🐛 Bug 列表:\n"
		for i, b := range bugs {
			if i >= 10 {
				msg += fmt.Sprintf("  ... 还有 %d 个\n", len(bugs)-10)
				break
			}
			msg += fmt.Sprintf("  • [%s] %s\n", b.Priority, b.Title)
		}
	}

	if len(tasks) > 0 {
		msg += "\n✅ 任务列表:\n"
		for i, t := range tasks {
			if i >= 10 {
				msg += fmt.Sprintf("  ... 还有 %d 个\n", len(tasks)-10)
				break
			}
			msg += fmt.Sprintf("  • [%s] %s\n", t.Priority, t.Title)
		}
	}

	if release.Summary != "" {
		msg += fmt.Sprintf("\n📝 概述:\n%s\n", release.Summary)
	}

	return msg
}

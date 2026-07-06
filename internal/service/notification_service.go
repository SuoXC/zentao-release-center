package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/yi-nology/zentao-release-center/internal/model"
)

type NotificationService struct {
	lanxinSvc *LanxinService
	emailSvc  *EmailService
}

type NotifyPreview struct {
	LanxinEnabled bool     `json:"lanxinEnabled"`
	LanxinMessage string   `json:"lanxinMessage,omitempty"`
	EmailEnabled  bool     `json:"emailEnabled"`
	EmailSubject  string   `json:"emailSubject,omitempty"`
	EmailHTML     string   `json:"emailHtml,omitempty"`
	EmailTo       []string `json:"emailTo,omitempty"`
}

type NotifyResult struct {
	LanxinSuccess bool   `json:"lanxinSuccess"`
	LanxinError   string `json:"lanxinError,omitempty"`
	EmailSuccess  bool   `json:"emailSuccess"`
	EmailError    string `json:"emailError,omitempty"`
}

func NewNotificationService(lanxinSvc *LanxinService, emailSvc *EmailService) *NotificationService {
	return &NotificationService{
		lanxinSvc: lanxinSvc,
		emailSvc:  emailSvc,
	}
}

func (s *NotificationService) BuildPreview(
	release *model.Release,
	items []*model.ReleaseItem,
	projectName, version string,
) *NotifyPreview {
	preview := &NotifyPreview{
		LanxinEnabled: s.lanxinSvc != nil && s.lanxinSvc.IsEnabled(),
		EmailEnabled:  s.emailSvc != nil && s.emailSvc.IsEnabled(),
	}
	if preview.LanxinEnabled {
		preview.LanxinMessage = s.lanxinSvc.BuildReleaseMessage(release, items, projectName, version)
	}
	if preview.EmailEnabled {
		preview.EmailSubject = s.emailSvc.BuildReleaseSubject(release, projectName, version)
		preview.EmailHTML = s.emailSvc.BuildReleaseHTML(release, items, projectName, version)
		preview.EmailTo = s.emailSvc.GetRecipients()
	}
	return preview
}

func (s *NotificationService) SendNow(
	release *model.Release,
	items []*model.ReleaseItem,
	projectName, version string,
) *NotifyResult {
	result := &NotifyResult{}
	var wg sync.WaitGroup

	if s.lanxinSvc != nil && s.lanxinSvc.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg := s.lanxinSvc.BuildReleaseMessage(release, items, projectName, version)
			if err := s.lanxinSvc.Send(msg); err != nil {
				result.LanxinError = err.Error()
				log.Printf("[Notification] lanxin send failed: %v", err)
			} else {
				result.LanxinSuccess = true
				log.Printf("[Notification] lanxin sent successfully for release %s", release.Keyword)
			}
		}()
	}

	if s.emailSvc != nil && s.emailSvc.IsEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subject := s.emailSvc.BuildReleaseSubject(release, projectName, version)
			htmlBody := s.emailSvc.BuildReleaseHTML(release, items, projectName, version)
			if errs := s.emailSvc.SendToAll(subject, htmlBody); len(errs) > 0 {
				for _, err := range errs {
					result.EmailError += fmt.Sprintf("%v; ", err)
					log.Printf("[Notification] email send failed: %v", err)
				}
			} else {
				result.EmailSuccess = true
				log.Printf("[Notification] email sent successfully for release %s", release.Keyword)
			}
		}()
	}

	wg.Wait()
	return result
}

func (s *NotificationService) NotifyReleasePublished(
	release *model.Release,
	items []*model.ReleaseItem,
	projectName string,
	version string,
) {
	s.SendNow(release, items, projectName, version)
}

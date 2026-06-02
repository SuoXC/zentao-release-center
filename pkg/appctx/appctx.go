package appctx

import (
	"github.com/yi-nology/zentao-release-center/internal/service"
)

var (
	ProjectSvc          *service.ProjectService
	ReleaseSvc          *service.ReleaseService
	ZentaoProxy         *service.ZentaoProxyService
	DeploymentSvc       *service.DeploymentService
	GitLabSvc           *service.GitLabService
	DockerImageSvc      *service.DockerImageService
	FeatureSvc          *service.FeatureService
	NotificationSvc     *service.NotificationService
	ZentaoBaseURL       string
	GitLabWebhookSecret string
)

package main

import (
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/yi-nology/zentao-release-center/internal/config"
	"github.com/yi-nology/zentao-release-center/internal/gitlab"
	"github.com/yi-nology/zentao-release-center/internal/service"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"github.com/yi-nology/zentao-release-center/internal/zentao"
	"github.com/yi-nology/zentao-release-center/pkg/appctx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dbStore, err := store.NewStore(cfg.Database.Path)
	if err != nil {
		log.Fatalf("init store: %v", err)
	}
	defer dbStore.Close()

	db := dbStore.DB()
	zentaoClient := zentao.NewClient(cfg.ZentaoMini.BaseURL, cfg.ZentaoMini.Timeout)
	gitlabClient := gitlab.NewClient(cfg.GitLab.BaseURL, cfg.GitLab.Token)

	lanxinSvc := service.NewLanxinService(cfg.Lanxin)
	emailSvc := service.NewEmailService(cfg.Email)
	notificationSvc := service.NewNotificationService(lanxinSvc, emailSvc)

	appctx.ProjectSvc = service.NewProjectService(db)
	appctx.ReleaseSvc = service.NewReleaseService(db, zentaoClient, notificationSvc)
	appctx.ZentaoProxy = service.NewZentaoProxyService(zentaoClient)
	appctx.DeploymentSvc = service.NewDeploymentService(db)
	appctx.GitLabSvc = service.NewGitLabService(db, gitlabClient)
	appctx.DockerImageSvc = service.NewDockerImageService(db, gitlabClient)
	appctx.FeatureSvc = service.NewFeatureService(db)
	appctx.NotificationSvc = notificationSvc
	appctx.ZentaoBaseURL = cfg.ZentaoMini.ZentaoBaseURL
	appctx.GitLabWebhookSecret = cfg.GitLab.WebhookSecret

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	h := server.Default(server.WithHostPorts(addr))

	register(h)
	_ = mountFrontend(h)
	h.Spin()
}

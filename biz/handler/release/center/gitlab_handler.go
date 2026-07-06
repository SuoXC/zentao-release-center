package center

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/gitlab"
	"github.com/yi-nology/zentao-release-center/pkg/appctx"
)

func AddRepo(ctx context.Context, c *app.RequestContext) {
	var req center.AddRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.RepoResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	repo, err := appctx.GitLabSvc.AddRepo(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.RepoResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.RepoResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: repo})
}

func DeleteRepo(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteRepoReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.GitLabSvc.DeleteRepo(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListRepos(ctx context.Context, c *app.RequestContext) {
	var req center.ListReposReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.RepoListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.GitLabSvc.ListRepos(req.ProjectId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.RepoListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.RepoListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

// CreateReleaseBranch 已废弃。
// 发布单与发布分支一对一强绑定，分支在创建发布单时（POST /api/releases）同步创建，
// 不再单独提供手工创建分支的入口。
func CreateReleaseBranch(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusGone, &center.BranchResp{
		Base: &center.BaseResp{
			Code:    410,
			Message: "请通过创建发布单（POST /api/releases）时同步创建发布分支；本接口已废弃",
		},
	})
}

func CreateFeatureBranch(ctx context.Context, c *app.RequestContext) {
	var req center.CreateFeatureBranchReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BranchResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	branch, err := appctx.GitLabSvc.CreateFeatureBranch(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BranchResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BranchResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: branch})
}

func DeleteBranch(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteBranchReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.GitLabSvc.DeleteBranch(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func UpdateBranch(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateBranchReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.GitLabSvc.UpdateBranch(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListBranches(ctx context.Context, c *app.RequestContext) {
	var req center.ListBranchesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BranchListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.GitLabSvc.ListBranchesByRelease(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BranchListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BranchListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func AddDockerImage(ctx context.Context, c *app.RequestContext) {
	var req center.AddDockerImageReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DockerImageResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	image, err := appctx.DockerImageSvc.AddManual(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DockerImageResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DockerImageResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: image})
}

func DeleteDockerImage(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteDockerImageReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.DockerImageSvc.Delete(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func UpdateDockerImage(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateDockerImageReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DockerImageResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	image, err := appctx.DockerImageSvc.Update(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DockerImageResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DockerImageResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: image})
}

func ListDockerImages(ctx context.Context, c *app.RequestContext) {
	var req center.ListDockerImagesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DockerImageListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.DockerImageSvc.List(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DockerImageListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DockerImageListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func SearchGitlabProjects(ctx context.Context, c *app.RequestContext) {
	var req center.SearchGitlabProjectsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.GitlabProjectListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.GitLabSvc.SearchProjects(req.Query)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.GitlabProjectListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.GitlabProjectListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func ListGitlabBranches(ctx context.Context, c *app.RequestContext) {
	var req center.ListGitlabBranchesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.GitlabBranchListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.GitLabSvc.ListBranches(int(req.GitlabProjectId))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.GitlabBranchListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.GitlabBranchListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func GitlabWebhook(ctx context.Context, c *app.RequestContext) {
	token := c.GetHeader("X-Gitlab-Token")
	if appctx.GitLabWebhookSecret != "" && string(token) != appctx.GitLabWebhookSecret {
		c.JSON(consts.StatusUnauthorized, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 401, Message: "invalid webhook secret"}})
		return
	}

	var event center.GitlabWebhookEvent
	if err := c.BindJSON(&event); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}

	if event.ObjectKind == "pipeline" {
		pipelineEvent := &gitlab.PipelineEvent{
			ObjectKind: event.ObjectKind,
			ObjectAttributes: gitlab.PipelineEventAttr{
				ID:        event.ObjectAttributes.ID,
				IID:       event.ObjectAttributes.IID,
				ProjectID: event.ObjectAttributes.ProjectID,
				Status:    event.ObjectAttributes.Status,
				Source:    event.ObjectAttributes.Source,
				Ref:       event.ObjectAttributes.Ref,
				SHA:       event.ObjectAttributes.SHA,
				WebURL:    event.ObjectAttributes.WebURL,
			},
			Project: gitlab.PipelineProject{
				ID:                event.Project.ID,
				Name:              event.Project.Name,
				PathWithNamespace: event.Project.PathWithNamespace,
				WebURL:            event.Project.WebURL,
			},
			Commit: gitlab.PipelineCommit{
				ID:      event.Commit.ID,
				Message: event.Commit.Message,
			},
		}
		if err := appctx.DockerImageSvc.AddFromWebhook(pipelineEvent); err != nil {
			c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
			return
		}
	}

	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func CIBuild(ctx context.Context, c *app.RequestContext) {
	var req center.CIBuildReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.DockerImageSvc.AddFromCIBuild(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListDockerImagePool(ctx context.Context, c *app.RequestContext) {
	var req center.ListDockerImagePoolReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DockerImagePoolListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.DockerImageSvc.ListPool(int(req.GetGitlabProjectId()))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DockerImagePoolListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DockerImagePoolListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

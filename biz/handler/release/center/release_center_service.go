package center

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/pkg/appctx"
)

func paging(page, pageSize int32) (int, int) {
	p, ps := int(page), int(pageSize)
	if p == 0 {
		p = 1
	}
	if ps == 0 {
		ps = 20
	}
	return p, ps
}

func CreateProject(ctx context.Context, c *app.RequestContext) {
	var req center.CreateProjectReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ProjectResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	project, err := appctx.ProjectSvc.Create(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ProjectResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ProjectResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: project})
}

func UpdateProject(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateProjectReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ProjectResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ProjectSvc.Update(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ProjectResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	project, err := appctx.ProjectSvc.Get(req.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ProjectResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ProjectResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: project})
}

func ListProjects(ctx context.Context, c *app.RequestContext) {
	var req center.ListProjectsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ListProjectsResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	page, pageSize := paging(req.GetPage(), req.GetPageSize())
	list, total, err := appctx.ProjectSvc.List(req.GetStatus(), page, pageSize)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ListProjectsResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ListProjectsResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list, Total: int32(total)})
}

func GetProject(ctx context.Context, c *app.RequestContext) {
	var req center.GetProjectReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ProjectResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	project, err := appctx.ProjectSvc.Get(req.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ProjectResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ProjectResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: project})
}

func DeleteProject(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteProjectReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ProjectSvc.Delete(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func CreateRelease(ctx context.Context, c *app.RequestContext) {
	var req center.CreateReleaseReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ReleaseResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	rel, err := appctx.ReleaseSvc.Create(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ReleaseResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ReleaseResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: rel})
}

func UpdateRelease(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateReleaseReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ReleaseResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.Update(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ReleaseResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ReleaseResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListReleases(ctx context.Context, c *app.RequestContext) {
	var req center.ListReleasesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ListReleasesResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	page, pageSize := paging(req.GetPage(), req.GetPageSize())
	list, total, err := appctx.ReleaseSvc.List(req.ProjectId, req.GetStatus(), page, pageSize)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ListReleasesResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ListReleasesResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list, Total: int32(total)})
}

func GetRelease(ctx context.Context, c *app.RequestContext) {
	var req center.GetReleaseReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ReleaseResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	rel, err := appctx.ReleaseSvc.Get(req.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ReleaseResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ReleaseResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: rel})
}

func DeleteRelease(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteReleaseReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.Delete(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func AddItem(ctx context.Context, c *app.RequestContext) {
	var req center.AddItemReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if _, err := appctx.ReleaseSvc.AddItem(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func BatchAddItems(ctx context.Context, c *app.RequestContext) {
	var req center.BatchAddItemsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if _, err := appctx.ReleaseSvc.BatchAddItems(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func UpdateItem(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateItemReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.UpdateItem(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func DeleteItem(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteItemReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.DeleteItem(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListItems(ctx context.Context, c *app.RequestContext) {
	var req center.ListItemsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ReleaseItemListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.ReleaseSvc.ListItems(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ReleaseItemListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.ReleaseItemListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func ReorderItems(ctx context.Context, c *app.RequestContext) {
	var req center.ReorderItemsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.ReorderItems(req.ReleaseId, req.Items); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func RefreshItems(ctx context.Context, c *app.RequestContext) {
	var req center.RefreshItemsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.ReleaseSvc.RefreshItems(req.ReleaseId); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func PublishRelease(ctx context.Context, c *app.RequestContext) {
	var req center.PublishReleaseReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.SnapshotResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	snap, err := appctx.ReleaseSvc.Publish(req.ReleaseId, req.GetVersion())
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.SnapshotResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.SnapshotResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: snap})
}

func ListSnapshots(ctx context.Context, c *app.RequestContext) {
	var req center.ListSnapshotsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.SnapshotListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.ReleaseSvc.ListSnapshots(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.SnapshotListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.SnapshotListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func GetSnapshot(ctx context.Context, c *app.RequestContext) {
	var req center.GetSnapshotReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.SnapshotResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	snap, err := appctx.ReleaseSvc.GetSnapshot(req.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.SnapshotResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.SnapshotResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: snap})
}

func ExportRelease(ctx context.Context, c *app.RequestContext) {
	var req center.ExportReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ExportResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	format := req.GetFormat()
	if format == "" {
		format = "markdown"
	}
	content, version, err := appctx.ReleaseSvc.Export(req.ReleaseId, req.GetSnapshotId(), format)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ExportResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ext := ".md"
	if format == "html" {
		ext = ".html"
	}
	filename := "release-v" + version + ext
	c.JSON(consts.StatusOK, &center.ExportResp{
		Base:     &center.BaseResp{Code: 0, Message: "ok"},
		Content:  &content,
		Filename: &filename,
		Format:   &format,
	})
}

func GetZentaoBugs(ctx context.Context, c *app.RequestContext) {
	var req center.ZentaoBugsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ZentaoPaginatedResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	page, pageSize := paging(req.GetPage(), req.GetPageSize())
	list, total, pg, ps, err := appctx.ZentaoProxy.GetBugs(int(req.GetProductId()), int(req.GetProjectId()), req.GetStatus(), page, pageSize)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ZentaoPaginatedResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ls := ""
	if list != nil {
		ls = string(list)
	}
	c.JSON(consts.StatusOK, &center.ZentaoPaginatedResp{
		Base:     &center.BaseResp{Code: 0, Message: "ok"},
		List:     &ls,
		Total:    int32(total),
		Page:     int32(pg),
		PageSize: int32(ps),
	})
}

func GetZentaoTasks(ctx context.Context, c *app.RequestContext) {
	var req center.ZentaoTasksReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.ZentaoPaginatedResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	page, pageSize := paging(req.GetPage(), req.GetPageSize())
	list, total, pg, ps, err := appctx.ZentaoProxy.GetTasks(int(req.GetExecutionId()), int(req.GetProductId()), req.GetStatus(), page, pageSize)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ZentaoPaginatedResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ls := ""
	if list != nil {
		ls = string(list)
	}
	c.JSON(consts.StatusOK, &center.ZentaoPaginatedResp{
		Base:     &center.BaseResp{Code: 0, Message: "ok"},
		List:     &ls,
		Total:    int32(total),
		Page:     int32(pg),
		PageSize: int32(ps),
	})
}

func GetZentaoProducts(ctx context.Context, c *app.RequestContext) {
	var req center.ZentaoProductsReq
	c.BindAndValidate(&req)
	data, err := appctx.ZentaoProxy.GetProducts()
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ds := ""
	if data != nil {
		ds = string(data)
	}
	c.JSON(consts.StatusOK, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: &ds})
}

func GetZentaoProjects(ctx context.Context, c *app.RequestContext) {
	var req center.ZentaoProjectsReq
	c.BindAndValidate(&req)
	data, err := appctx.ZentaoProxy.GetProjects(int(req.GetProductId()))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ds := ""
	if data != nil {
		ds = string(data)
	}
	c.JSON(consts.StatusOK, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: &ds})
}

func GetZentaoExecutions(ctx context.Context, c *app.RequestContext) {
	var req center.ZentaoExecutionsReq
	c.BindAndValidate(&req)
	data, err := appctx.ZentaoProxy.GetExecutions(int(req.GetProjectId()))
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	ds := ""
	if data != nil {
		ds = string(data)
	}
	c.JSON(consts.StatusOK, &center.ZentaoDataResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: &ds})
}

func Health(ctx context.Context, c *app.RequestContext) {
	status := "ok"
	zentaoStatus := "unknown"
	if _, err := appctx.ZentaoProxy.GetProducts(); err == nil {
		zentaoStatus = "connected"
	} else {
		zentaoStatus = "disconnected"
	}
	zb := appctx.ZentaoBaseURL
	c.JSON(consts.StatusOK, &center.HealthResp{
		Base:             &center.BaseResp{Code: 0, Message: "ok"},
		Status:           &status,
		ZentaoMiniStatus: &zentaoStatus,
		ZentaoBaseUrl:    &zb,
	})
}

// AddDeployment .
// @router /api/deployments [POST]
func AddDeployment(ctx context.Context, c *app.RequestContext) {
	var req center.AddDeploymentReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DeploymentResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	dep, err := appctx.DeploymentSvc.Add(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DeploymentResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DeploymentResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: dep})
}

func UpdateDeployment(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateDeploymentReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DeploymentResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.DeploymentSvc.Update(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DeploymentResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	dep, err := appctx.DeploymentSvc.Get(req.ID)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DeploymentResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DeploymentResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: dep})
}

func DeleteDeployment(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteDeploymentReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.DeploymentSvc.Delete(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListDeployments(ctx context.Context, c *app.RequestContext) {
	var req center.ListDeploymentsReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.DeploymentListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.DeploymentSvc.List(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.DeploymentListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.DeploymentListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

// ==================== 功能说明 ====================

func AddFeature(ctx context.Context, c *app.RequestContext) {
	var req center.AddFeatureReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.FeatureResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	feature, err := appctx.FeatureSvc.Add(&req)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.FeatureResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.FeatureResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, Data: feature})
}

func UpdateFeature(ctx context.Context, c *app.RequestContext) {
	var req center.UpdateFeatureReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.FeatureSvc.Update(&req); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func DeleteFeature(ctx context.Context, c *app.RequestContext) {
	var req center.DeleteFeatureReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.FeatureSvc.Delete(req.ID); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

func ListFeatures(ctx context.Context, c *app.RequestContext) {
	var req center.ListFeaturesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.FeatureListResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	list, err := appctx.FeatureSvc.List(req.ReleaseId)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, &center.FeatureListResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.FeatureListResp{Base: &center.BaseResp{Code: 0, Message: "ok"}, List: list})
}

func ReorderFeatures(ctx context.Context, c *app.RequestContext) {
	var req center.ReorderFeaturesReq
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 400, Message: err.Error()}})
		return
	}
	if err := appctx.FeatureSvc.Reorder(req.ReleaseId, req.Items); err != nil {
		c.JSON(consts.StatusInternalServerError, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 500, Message: err.Error()}})
		return
	}
	c.JSON(consts.StatusOK, &center.BaseOnlyResp{Base: &center.BaseResp{Code: 0, Message: "ok"}})
}

package center

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/yi-nology/zentao-release-center/pkg/appctx"
)

func NotifyPreview(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ReleaseId string `json:"releaseId"`
		Version   string `json:"version"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"base": map[string]interface{}{"code": 400, "message": err.Error()}})
		return
	}
	if req.ReleaseId == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"base": map[string]interface{}{"code": 400, "message": "releaseId is required"}})
		return
	}

	svc := appctx.NotificationSvc
	if svc == nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": "notification service not initialized"}})
		return
	}

	release, err := appctx.ReleaseSvc.GetRaw(req.ReleaseId)
	if err != nil || release == nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("release not found: %v", err)}})
		return
	}

	items, err := appctx.ReleaseSvc.GetRawItems(req.ReleaseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("load items failed: %v", err)}})
		return
	}

	deployments, err := appctx.ReleaseSvc.GetDeployments(req.ReleaseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("load deployments failed: %v", err)}})
		return
	}

	projectName := ""
	project, _ := appctx.ProjectSvc.GetRaw(release.ProjectKeyword)
	if project != nil {
		projectName = project.Name
	}

	preview := svc.BuildPreview(release, items, deployments, projectName, req.Version)

	respBytes, _ := json.Marshal(preview)
	var respData interface{}
	json.Unmarshal(respBytes, &respData)

	c.JSON(http.StatusOK, map[string]interface{}{
		"base": map[string]interface{}{"code": 0, "message": "ok"},
		"data": respData,
	})
}

func NotifySend(ctx context.Context, c *app.RequestContext) {
	var req struct {
		ReleaseId string `json:"releaseId"`
		Version   string `json:"version"`
		Channel   string `json:"channel"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"base": map[string]interface{}{"code": 400, "message": err.Error()}})
		return
	}
	if req.ReleaseId == "" {
		c.JSON(http.StatusBadRequest, map[string]interface{}{"base": map[string]interface{}{"code": 400, "message": "releaseId is required"}})
		return
	}

	svc := appctx.NotificationSvc
	if svc == nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": "notification service not initialized"}})
		return
	}

	release, err := appctx.ReleaseSvc.GetRaw(req.ReleaseId)
	if err != nil || release == nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("release not found: %v", err)}})
		return
	}

	items, err := appctx.ReleaseSvc.GetRawItems(req.ReleaseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("load items failed: %v", err)}})
		return
	}

	deployments, err := appctx.ReleaseSvc.GetDeployments(req.ReleaseId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{"base": map[string]interface{}{"code": 500, "message": fmt.Sprintf("load deployments failed: %v", err)}})
		return
	}

	projectName := ""
	project, _ := appctx.ProjectSvc.GetRaw(release.ProjectKeyword)
	if project != nil {
		projectName = project.Name
	}

	result := svc.SendNow(release, items, deployments, projectName, req.Version)

	respBytes, _ := json.Marshal(result)
	var respData interface{}
	json.Unmarshal(respBytes, &respData)

	c.JSON(http.StatusOK, map[string]interface{}{
		"base": map[string]interface{}{"code": 0, "message": "ok"},
		"data": respData,
	})
}

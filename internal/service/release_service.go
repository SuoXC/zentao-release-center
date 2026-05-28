package service

import (
	"encoding/json"
	"fmt"
	"strings"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"github.com/yi-nology/zentao-release-center/internal/zentao"
)

type ReleaseService struct {
	releaseStore    *store.ReleaseStore
	itemStore       *store.ItemStore
	snapshotStore   *store.SnapshotStore
	projectStore    *store.ProjectStore
	deploymentStore *store.DeploymentStore
	zentaoClient    *zentao.Client
}

func NewReleaseService(s *store.Store, zc *zentao.Client) *ReleaseService {
	return &ReleaseService{
		releaseStore:    store.NewReleaseStore(s),
		itemStore:       store.NewItemStore(s),
		snapshotStore:   store.NewSnapshotStore(s),
		projectStore:    store.NewProjectStore(s),
		deploymentStore: store.NewDeploymentStore(s),
		zentaoClient:    zc,
	}
}

func (rs *ReleaseService) Create(req *center.CreateReleaseReq) (*center.Release, error) {
	if req.ProjectId == "" {
		return nil, fmt.Errorf("projectId is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	project, err := rs.projectStore.GetByID(req.ProjectId)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}
	rel, err := rs.releaseStore.Create(req.ProjectId, req.Name, req.GetVersion(), req.GetSummary())
	if err != nil {
		return nil, err
	}
	rs.enrichRelease(rel, project)
	return rel, nil
}

func (rs *ReleaseService) Get(id string) (*center.Release, error) {
	rel, err := rs.releaseStore.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, fmt.Errorf("release not found")
	}
	rs.fillCounts(rel)
	project, _ := rs.projectStore.GetByID(rel.ProjectId)
	if project != nil {
		rs.enrichRelease(rel, project)
	}
	return rel, nil
}

func (rs *ReleaseService) List(projectID, status string, page, pageSize int) ([]*center.Release, int, error) {
	list, total, err := rs.releaseStore.List(projectID, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	for _, rel := range list {
		rs.fillCounts(rel)
		project, _ := rs.projectStore.GetByID(rel.ProjectId)
		if project != nil {
			rs.enrichRelease(rel, project)
		}
	}
	return list, total, nil
}

func (rs *ReleaseService) Update(req *center.UpdateReleaseReq) error {
	if req.ID == "" {
		return fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.IsSetName() {
		fields["name"] = req.Name
	}
	if req.IsSetVersion() {
		fields["version"] = req.Version
	}
	if req.IsSetSummary() {
		fields["summary"] = req.Summary
	}
	if req.IsSetStatus() {
		fields["status"] = req.Status
	}
	return rs.releaseStore.Update(req.ID, fields)
}

func (rs *ReleaseService) Delete(id string) error {
	return rs.releaseStore.Delete(id)
}

func (rs *ReleaseService) AddItem(req *center.AddItemReq) (*center.ReleaseItem, error) {
	return rs.itemStore.Add(req.ReleaseId, req.ItemType,
		int(req.GetZentaoId()), req.GetZentaoType(), req.GetTitle(),
		req.GetSeverity(), req.GetPriority(), req.GetStatus(),
		req.GetAssignedTo(), req.GetResolvedBy(), req.GetZentaoUrl(),
		req.GetSteps(), req.GetNoteTitle(), req.GetNoteContent())
}

func (rs *ReleaseService) BatchAddItems(req *center.BatchAddItemsReq) ([]*center.ReleaseItem, error) {
	var toAdd []struct {
		ItemType, ZentaoType, Title, Severity, Priority, Status, AssignedTo, ResolvedBy, ZentaoURL, Steps, NoteTitle, NoteContent string
		ZentaoID                                                                                                                  int
	}

	for _, item := range req.Items {
		zentaoID := 0
		if item.ZentaoId != nil {
			zentaoID = int(*item.ZentaoId)
		}
		if zentaoID > 0 {
			exists, _ := rs.itemStore.ExistsByZentaoID(req.ReleaseId, zentaoID)
			if exists {
				continue
			}
		}
		toAdd = append(toAdd, struct {
			ItemType, ZentaoType, Title, Severity, Priority, Status, AssignedTo, ResolvedBy, ZentaoURL, Steps, NoteTitle, NoteContent string
			ZentaoID                                                                                                                  int
		}{
			ItemType:    item.GetItemType(),
			ZentaoType:  item.GetZentaoType(),
			Title:       item.GetTitle(),
			Severity:    item.GetSeverity(),
			Priority:    item.GetPriority(),
			Status:      item.GetStatus(),
			AssignedTo:  item.GetAssignedTo(),
			ResolvedBy:  item.GetResolvedBy(),
			ZentaoURL:   item.GetZentaoUrl(),
			Steps:       item.GetSteps(),
			NoteTitle:   item.GetNoteTitle(),
			NoteContent: item.GetNoteContent(),
			ZentaoID:    zentaoID,
		})
	}

	if len(toAdd) == 0 {
		return nil, nil
	}

	tx, err := rs.itemStore.DB().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	items, err := rs.itemStore.AddBatch(tx, req.ReleaseId, toAdd)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return items, nil
}

func (rs *ReleaseService) ListItems(releaseID string) ([]*center.ReleaseItem, error) {
	return rs.itemStore.ListByRelease(releaseID)
}

func (rs *ReleaseService) UpdateItem(req *center.UpdateItemReq) error {
	fields := map[string]interface{}{}
	if req.IsSetNoteTitle() {
		fields["note_title"] = req.NoteTitle
	}
	if req.IsSetNoteContent() {
		fields["note_content"] = req.NoteContent
	}
	if req.IsSetSortOrder() {
		fields["sort_order"] = req.SortOrder
	}
	return rs.itemStore.Update(req.ID, fields)
}

func (rs *ReleaseService) DeleteItem(id string) error {
	return rs.itemStore.Delete(id)
}

func (rs *ReleaseService) ReorderItems(releaseID string, items []*center.SortItem) error {
	sortItems := make([]struct{ ID string; SortOrder int }, len(items))
	for i, item := range items {
		sortItems[i] = struct{ ID string; SortOrder int }{item.ID, int(item.SortOrder)}
	}
	return rs.itemStore.Reorder(sortItems)
}

func (rs *ReleaseService) RefreshItems(releaseID string) error {
	items, err := rs.itemStore.ListByRelease(releaseID)
	if err != nil {
		return err
	}

	rel, err := rs.releaseStore.GetByID(releaseID)
	if err != nil || rel == nil {
		return fmt.Errorf("release not found")
	}

	project, _ := rs.projectStore.GetByID(rel.ProjectId)
	if project == nil || project.ZentaoProductId == 0 {
		return fmt.Errorf("project has no zentao product configured")
	}

	var bugs, tasks []*center.ReleaseItem
	for _, item := range items {
		switch item.ItemType {
		case "bug":
			bugs = append(bugs, item)
		case "task":
			tasks = append(tasks, item)
		}
	}

	if len(bugs) > 0 {
		bugMap := rs.fetchBugMap(int(project.ZentaoProductId))
		for _, item := range bugs {
			if item.ZentaoId == nil {
				continue
			}
			if b, ok := bugMap[int(*item.ZentaoId)]; ok {
				rs.itemStore.Update(item.ID, map[string]interface{}{
					"title":       b.Title,
					"severity":    fmt.Sprintf("%v", b.Severity),
					"priority":    fmt.Sprintf("%v", b.Pri),
					"status":      b.Status,
					"assigned_to": b.AssignedTo.Realname,
					"resolved_by": fmt.Sprintf("%v", b.ResolvedBy),
					"steps":       b.Steps,
				})
			}
		}
	}

	if len(tasks) > 0 {
		taskMap := rs.fetchTaskMap(int(project.ZentaoProductId))
		for _, item := range tasks {
			if item.ZentaoId == nil {
				continue
			}
			if t, ok := taskMap[int(*item.ZentaoId)]; ok {
				rs.itemStore.Update(item.ID, map[string]interface{}{
					"title":       t.Name,
					"priority":    fmt.Sprintf("%v", t.Pri),
					"status":      t.Status,
					"assigned_to": t.AssignedTo.Realname,
				})
			}
		}
	}

	return nil
}

type zentaoBug struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Severity interface{} `json:"severity"`
	Pri      interface{} `json:"pri"`
	Status   string `json:"status"`
	AssignedTo struct {
		Realname string `json:"realname"`
	} `json:"assignedTo"`
	ResolvedBy interface{} `json:"resolvedBy"`
	Steps      string `json:"steps"`
}

type zentaoTask struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Pri   interface{} `json:"pri"`
	Status string `json:"status"`
	AssignedTo struct {
		Realname string `json:"realname"`
	} `json:"assignedTo"`
}

func (rs *ReleaseService) fetchBugMap(productID int) map[int]zentaoBug {
	data, err := rs.zentaoClient.GetBugs(productID, 0, "", 1, 500)
	if err != nil {
		return nil
	}
	var result struct {
		List []zentaoBug `json:"list"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	m := make(map[int]zentaoBug, len(result.List))
	for _, b := range result.List {
		m[b.ID] = b
	}
	return m
}

func (rs *ReleaseService) fetchTaskMap(productID int) map[int]zentaoTask {
	data, err := rs.zentaoClient.GetTasks(0, productID, "", 1, 500)
	if err != nil {
		return nil
	}
	var result struct {
		List []zentaoTask `json:"list"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil
	}
	m := make(map[int]zentaoTask, len(result.List))
	for _, t := range result.List {
		m[t.ID] = t
	}
	return m
}

func (rs *ReleaseService) Publish(releaseID, version string) (*center.ReleaseSnapshot, error) {
	rel, err := rs.releaseStore.GetByID(releaseID)
	if err != nil || rel == nil {
		return nil, fmt.Errorf("release not found")
	}

	items, err := rs.itemStore.ListByRelease(releaseID)
	if err != nil {
		return nil, err
	}

	total, bugs, tasks, notes, err := rs.itemStore.CountByType(releaseID)
	if err != nil {
		return nil, err
	}

	deployments, err := rs.deploymentStore.ListByRelease(releaseID)
	if err != nil {
		return nil, err
	}

	content := rs.generateMarkdown(rel, items, deployments, version)

	snap, err := rs.snapshotStore.Create(releaseID, version, content, total, bugs, tasks, notes)
	if err != nil {
		return nil, err
	}

	if err := rs.releaseStore.IncrementPublish(releaseID); err != nil {
		return nil, err
	}

	return snap, nil
}

func (rs *ReleaseService) ListSnapshots(releaseID string) ([]*center.ReleaseSnapshot, error) {
	return rs.snapshotStore.ListByRelease(releaseID)
}

func (rs *ReleaseService) GetSnapshot(id string) (*center.ReleaseSnapshot, error) {
	snap, err := rs.snapshotStore.GetByID(id)
	if err != nil {
		return nil, err
	}
	if snap == nil {
		return nil, fmt.Errorf("snapshot not found")
	}
	return snap, nil
}

func (rs *ReleaseService) Export(releaseID, snapshotID, format string) (string, string, error) {
	var content string
	var version string

	if snapshotID != "" {
		snap, err := rs.snapshotStore.GetByID(snapshotID)
		if err != nil || snap == nil {
			return "", "", fmt.Errorf("snapshot not found")
		}
		content = snap.Content
		version = snap.Version
	} else {
		rel, err := rs.releaseStore.GetByID(releaseID)
		if err != nil || rel == nil {
			return "", "", fmt.Errorf("release not found")
		}
		items, err := rs.itemStore.ListByRelease(releaseID)
		if err != nil {
			return "", "", err
		}
		deployments, err := rs.deploymentStore.ListByRelease(releaseID)
		if err != nil {
			return "", "", err
		}
		version = rel.Version
		content = rs.generateMarkdown(rel, items, deployments, rel.Version)
	}

	if format == "html" {
		content = rs.markdownToHTML(content)
	}

	return content, version, nil
}

func (rs *ReleaseService) generateMarkdown(rel *center.Release, items []*center.ReleaseItem, deployments []*center.Deployment, version string) string {
	var sb strings.Builder

	sb.WriteString("# ")
	sb.WriteString(rel.Name)
	if version != "" {
		sb.WriteString(" v")
		sb.WriteString(version)
	}
	sb.WriteString("\n\n")

	if rel.Summary != "" {
		sb.WriteString("## 概述\n\n")
		sb.WriteString(rel.Summary)
		sb.WriteString("\n\n")
	}

	if len(deployments) > 0 {
		sb.WriteString("## 部署地址\n\n")
		sb.WriteString("| 功能模块 | 地址 | 说明 |\n")
		sb.WriteString("|---------|------|------|\n")
		for _, d := range deployments {
			desc := d.Description
			if desc == "" {
				desc = "-"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", d.ModuleName, d.Address, desc))
		}
		sb.WriteString("\n")
	}

	var bugs, tasks, notes []*center.ReleaseItem
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

	if len(bugs) > 0 {
		sb.WriteString(fmt.Sprintf("## Bug 修复（%d）\n\n", len(bugs)))
		sb.WriteString("| # | 标题 | 严重程度 | 优先级 | 状态 | 指派给 |\n")
		sb.WriteString("|---|------|---------|--------|------|--------|\n")
		for _, b := range bugs {
		link := "-"
		zentaoURL := ""
		zentaoID := 0
		if b.ZentaoUrl != nil {
			zentaoURL = *b.ZentaoUrl
		}
		if b.ZentaoId != nil {
			zentaoID = int(*b.ZentaoId)
		}
		if zentaoURL != "" {
			link = fmt.Sprintf("[%d](%s)", zentaoID, zentaoURL)
		} else if zentaoID > 0 {
			link = fmt.Sprintf("%d", zentaoID)
		}
		title := b.GetTitle()
		if title == "" {
			title = "-"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			link, title, b.GetSeverity(), b.GetPriority(), b.GetStatus(), b.GetAssignedTo()))
	}
	sb.WriteString("\n")
}

if len(tasks) > 0 {
	sb.WriteString(fmt.Sprintf("## 任务完成（%d）\n\n", len(tasks)))
	sb.WriteString("| # | 标题 | 优先级 | 状态 | 指派给 |\n")
	sb.WriteString("|---|------|--------|------|--------|\n")
	for _, t := range tasks {
		link := "-"
		zentaoURL := ""
		zentaoID := 0
		if t.ZentaoUrl != nil {
			zentaoURL = *t.ZentaoUrl
		}
		if t.ZentaoId != nil {
			zentaoID = int(*t.ZentaoId)
		}
		if zentaoURL != "" {
			link = fmt.Sprintf("[%d](%s)", zentaoID, zentaoURL)
		} else if zentaoID > 0 {
			link = fmt.Sprintf("%d", zentaoID)
		}
			title := t.GetTitle()
			if title == "" {
				title = "-"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				link, title, t.GetPriority(), t.GetStatus(), t.GetAssignedTo()))
		}
		sb.WriteString("\n")
	}

	if len(notes) > 0 {
		sb.WriteString(fmt.Sprintf("## 备注（%d）\n\n", len(notes)))
		for _, n := range notes {
			sb.WriteString(fmt.Sprintf("### %s\n\n%s\n\n---\n\n", n.GetNoteTitle(), n.GetNoteContent()))
		}
	}

	sb.WriteString("---\n*Generated by zentao-release-center*\n")
	return sb.String()
}

func (rs *ReleaseService) markdownToHTML(md string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Release Note</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif;max-width:800px;margin:40px auto;padding:0 20px;color:#333;line-height:1.6}
table{border-collapse:collapse;width:100%%;margin:16px 0}
th,td{border:1px solid #ddd;padding:8px 12px;text-align:left}
th{background:#f5f5f5}
a{color:#4F6BF6}
h1{border-bottom:2px solid #4F6BF6;padding-bottom:8px}
h2{color:#1E293B;margin-top:24px}
hr{border:none;border-top:1px solid #e2e8f0;margin:16px 0}
</style></head><body>%s</body></html>`, rs.simpleMDToHTML(md))
}

func (rs *ReleaseService) simpleMDToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var html strings.Builder
	inTable := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			html.WriteString("<h1>" + strings.TrimPrefix(trimmed, "# ") + "</h1>\n")
		} else if strings.HasPrefix(trimmed, "## ") {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			html.WriteString("<h2>" + strings.TrimPrefix(trimmed, "## ") + "</h2>\n")
		} else if strings.HasPrefix(trimmed, "### ") {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			html.WriteString("<h3>" + strings.TrimPrefix(trimmed, "### ") + "</h3>\n")
		} else if strings.HasPrefix(trimmed, "|") {
			if !inTable {
				html.WriteString("<table>\n")
				inTable = true
			}
			if strings.HasPrefix(trimmed, "|---") {
				continue
			}
			html.WriteString("<tr>")
			cells := strings.Split(trimmed, "|")
			for _, cell := range cells {
				cell = strings.TrimSpace(cell)
				if cell == "" {
					continue
				}
				html.WriteString("<td>" + cell + "</td>")
			}
			html.WriteString("</tr>\n")
		} else if trimmed == "---" {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			html.WriteString("<hr>\n")
		} else if trimmed != "" {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			html.WriteString("<p>" + trimmed + "</p>\n")
		}
	}
	if inTable {
		html.WriteString("</table>\n")
	}
	return html.String()
}

func (rs *ReleaseService) enrichRelease(rel *center.Release, project *center.Project) {
	rel.ProjectName = project.Name
}

func (rs *ReleaseService) fillCounts(rel *center.Release) {
	total, bugs, tasks, notes, err := rs.itemStore.CountByType(rel.ID)
	if err != nil {
		return
	}
	rel.ItemCount = int32(total)
	rel.BugCount = int32(bugs)
	rel.TaskCount = int32(tasks)
	rel.NoteCount = int32(notes)
}

type ZentaoProxyService struct {
	client *zentao.Client
}

func NewZentaoProxyService(zc *zentao.Client) *ZentaoProxyService {
	return &ZentaoProxyService{client: zc}
}

func (zs *ZentaoProxyService) GetBugs(productID, projectID int, status string, page, pageSize int) (json.RawMessage, int, int, int, error) {
	data, err := zs.client.GetBugs(productID, projectID, status, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	var result struct {
		List     json.RawMessage `json:"list"`
		Total    int             `json:"total"`
		Page     int             `json:"page"`
		PageSize int             `json:"pageSize"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, 0, 0, 0, err
	}
	return result.List, result.Total, result.Page, result.PageSize, nil
}

func (zs *ZentaoProxyService) GetTasks(executionID, productID int, status string, page, pageSize int) (json.RawMessage, int, int, int, error) {
	data, err := zs.client.GetTasks(executionID, productID, status, page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	var result struct {
		List     json.RawMessage `json:"list"`
		Total    int             `json:"total"`
		Page     int             `json:"page"`
		PageSize int             `json:"pageSize"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, 0, 0, 0, err
	}
	return result.List, result.Total, result.Page, result.PageSize, nil
}

func (zs *ZentaoProxyService) GetProducts() (json.RawMessage, error) {
	return zs.client.GetProducts()
}

func (zs *ZentaoProxyService) GetProjects(productID int) (json.RawMessage, error) {
	return zs.client.GetProjects(productID)
}

func (zs *ZentaoProxyService) GetExecutions(projectID int) (json.RawMessage, error) {
	return zs.client.GetExecutions(projectID)
}

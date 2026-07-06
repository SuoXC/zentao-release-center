package service

import (
	"encoding/json"
	"fmt"
	"strings"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/gitlab"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"github.com/yi-nology/zentao-release-center/internal/zentao"
	"gorm.io/gorm"
)

type ReleaseService struct {
	releaseStore     *store.ReleaseStore
	itemStore        *store.ItemStore
	snapshotStore    *store.SnapshotStore
	projectStore     *store.ProjectStore
	repoStore        *store.RepoStore
	branchStore      *store.BranchStore
	dockerImageStore *store.DockerImageStore
	featureStore     *store.FeatureStore
	zentaoClient     *zentao.Client
	gitlabClient     *gitlab.Client
	notificationSvc  *NotificationService
}

func NewReleaseService(db *gorm.DB, zc *zentao.Client, gc *gitlab.Client, ns *NotificationService) *ReleaseService {
	return &ReleaseService{
		releaseStore:     store.NewReleaseStore(db),
		itemStore:        store.NewItemStore(db),
		snapshotStore:    store.NewSnapshotStore(db),
		projectStore:     store.NewProjectStore(db),
		repoStore:        store.NewRepoStore(db),
		branchStore:      store.NewBranchStore(db),
		dockerImageStore: store.NewDockerImageStore(db),
		featureStore:     store.NewFeatureStore(db),
		zentaoClient:     zc,
		gitlabClient:     gc,
		notificationSvc:  ns,
	}
}

func (rs *ReleaseService) Create(req *center.CreateReleaseReq) (*center.Release, error) {
	if req.ProjectId == "" {
		return nil, fmt.Errorf("projectId is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if req.RepoId == "" {
		return nil, fmt.Errorf("repoId is required")
	}

	project, err := rs.projectStore.GetByID(req.ProjectId)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	repo, err := rs.repoStore.GetByKeyword(req.RepoId)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repo not found")
	}
	if repo.ProjectKeyword != req.ProjectId {
		return nil, fmt.Errorf("repo 不属于该项目")
	}

	parentBranch := req.GetParentBranch()
	if parentBranch == "" {
		parentBranch = repo.DefaultBranch
	}
	branchName := buildReleaseBranchName(req.GetVersion(), req.Name)
	description := fmt.Sprintf("发布单 %s 自动创建", req.Name)
	if req.GetSummary() != "" {
		description = req.GetSummary()
	}

	gitlabBranchURL, err := rs.gitlabClient.CreateBranch(repo.GitlabProjectID, branchName, parentBranch)
	if err != nil {
		return nil, fmt.Errorf("创建 GitLab 分支失败：%w", err)
	}
	if gitlabBranchURL == "" {
		gitlabBranchURL = repo.RepoURL + "/-/tree/" + branchName
	}

	rel, err := rs.releaseStore.Create(req.ProjectId, req.Name, req.GetVersion(), req.GetSummary(), parentBranch)
	if err != nil {
		// 回滚 GitLab 分支
		_ = rs.gitlabClient.DeleteBranch(repo.GitlabProjectID, branchName)
		return nil, fmt.Errorf("create release failed: %w", err)
	}

	if _, err := rs.branchStore.Create(rel.Keyword, repo.Keyword, branchName, "release", parentBranch, gitlabBranchURL, description); err != nil {
		_ = rs.gitlabClient.DeleteBranch(repo.GitlabProjectID, branchName)
		_ = rs.releaseStore.Delete(rel.Keyword)
		return nil, fmt.Errorf("create release branch failed: %w", err)
	}

	return mapper.ReleaseToThrift(rel, project.Name, 0, 0, 0, 0), nil
}

// buildReleaseBranchName 根据发布单名称/版本生成 release 分支名。
// 默认形式：release/<version>（去除非合法字符）；若 version 为空，回退到 release/<name slug>-<keyword 前 6 位>。
func buildReleaseBranchName(version, name string) string {
	safe := func(s string) string {
		out := make([]rune, 0, len(s))
		for _, r := range s {
			switch {
			case r >= 'a' && r <= 'z',
				r >= 'A' && r <= 'Z',
				r >= '0' && r <= '9',
				r == '-', r == '_', r == '.':
				out = append(out, r)
			default:
				out = append(out, '-')
			}
		}
		return strings.Trim(string(out), "-")
	}
	if v := safe(version); v != "" {
		return "release/" + v
	}
	return "release/" + safe(name)
}

func (rs *ReleaseService) Get(keyword string) (*center.Release, error) {
	rel, err := rs.releaseStore.GetByID(keyword)
	if err != nil {
		return nil, err
	}
	if rel == nil {
		return nil, fmt.Errorf("release not found")
	}
	total, bugs, tasks, notes, _ := rs.itemStore.CountByType(keyword)
	project, _ := rs.projectStore.GetByID(rel.ProjectKeyword)
	projectName := ""
	if project != nil {
		projectName = project.Name
	}
	return mapper.ReleaseToThrift(rel, projectName, total, bugs, tasks, notes), nil
}

func (rs *ReleaseService) List(projectKeyword, status string, page, pageSize int) ([]*center.Release, int, error) {
	releases, total, err := rs.releaseStore.List(projectKeyword, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*center.Release, len(releases))
	for i, rel := range releases {
		total, bugs, tasks, notes, _ := rs.itemStore.CountByType(rel.Keyword)
		project, _ := rs.projectStore.GetByID(rel.ProjectKeyword)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}
		result[i] = mapper.ReleaseToThrift(rel, projectName, total, bugs, tasks, notes)
	}
	return result, total, nil
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
	if req.IsSetParentBranch() {
		fields["parent_branch"] = req.ParentBranch
	}
	return rs.releaseStore.Update(req.ID, fields)
}

func (rs *ReleaseService) Delete(keyword string) error {
	return rs.releaseStore.Delete(keyword)
}

func (rs *ReleaseService) AddItem(req *center.AddItemReq) (*center.ReleaseItem, error) {
	item, err := rs.itemStore.Add(req.ReleaseId, req.ItemType,
		int(req.GetZentaoId()), req.GetZentaoType(), req.GetTitle(),
		req.GetSeverity(), req.GetPriority(), req.GetStatus(),
		req.GetAssignedTo(), req.GetResolvedBy(), req.GetZentaoUrl(),
		req.GetSteps(), req.GetNoteTitle(), req.GetNoteContent())
	if err != nil {
		return nil, err
	}
	return mapper.ItemToThrift(item), nil
}

func (rs *ReleaseService) BatchAddItems(req *center.BatchAddItemsReq) ([]*center.ReleaseItem, error) {
	var toAdd []struct {
		ItemType, ZentaoType, Title, Severity, Priority, Status, AssignedTo, ResolvedBy, ZentaoURL, Steps, NoteTitle, NoteContent string
		ZentaoID int
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
			ZentaoID int
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

	items, err := rs.itemStore.AddBatch(req.ReleaseId, toAdd)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ReleaseItem, len(items))
	for i, item := range items {
		result[i] = mapper.ItemToThrift(item)
	}
	return result, nil
}

func (rs *ReleaseService) ListItems(releaseKeyword string) ([]*center.ReleaseItem, error) {
	items, err := rs.itemStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ReleaseItem, len(items))
	for i, item := range items {
		result[i] = mapper.ItemToThrift(item)
	}
	return result, nil
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

func (rs *ReleaseService) DeleteItem(keyword string) error {
	return rs.itemStore.Delete(keyword)
}

func (rs *ReleaseService) ReorderItems(releaseKeyword string, items []*center.SortItem) error {
	sortItems := make([]struct{ Keyword string; SortOrder int }, len(items))
	for i, item := range items {
		sortItems[i] = struct{ Keyword string; SortOrder int }{item.ID, int(item.SortOrder)}
	}
	return rs.itemStore.Reorder(sortItems)
}

func (rs *ReleaseService) RefreshItems(releaseKeyword string) error {
	items, err := rs.itemStore.ListByRelease(releaseKeyword)
	if err != nil {
		return err
	}

	rel, err := rs.releaseStore.GetByID(releaseKeyword)
	if err != nil || rel == nil {
		return fmt.Errorf("release not found")
	}

	project, _ := rs.projectStore.GetByID(rel.ProjectKeyword)
	if project == nil || project.ZentaoProductID == 0 {
		return fmt.Errorf("project has no zentao product configured")
	}

	var bugs, tasks []*center.ReleaseItem
	for _, item := range items {
		t := mapper.ItemToThrift(item)
		switch t.ItemType {
		case "bug":
			bugs = append(bugs, t)
		case "task":
			tasks = append(tasks, t)
		}
	}

	if len(bugs) > 0 {
		bugMap := rs.fetchBugMap(project.ZentaoProductID)
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
		taskMap := rs.fetchTaskMap(project.ZentaoProductID)
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

func (rs *ReleaseService) Publish(releaseKeyword, version string) (*center.ReleaseSnapshot, error) {
	rel, err := rs.releaseStore.GetByID(releaseKeyword)
	if err != nil || rel == nil {
		return nil, fmt.Errorf("release not found")
	}

	items, err := rs.itemStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}

	total, bugs, tasks, notes, err := rs.itemStore.CountByType(releaseKeyword)
	if err != nil {
		return nil, err
	}

	branches, err := rs.branchStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}

	dockerImages, err := rs.dockerImageStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}

	features, err := rs.featureStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}

	thriftItems := make([]*center.ReleaseItem, len(items))
	for i, item := range items {
		thriftItems[i] = mapper.ItemToThrift(item)
	}
	thriftBranches := make([]*center.ReleaseBranch, len(branches))
	for i, b := range branches {
		thriftBranches[i] = mapper.ReleaseBranchToThrift(b)
	}
	thriftImages := make([]*center.DockerImage, len(dockerImages))
	for i, img := range dockerImages {
		thriftImages[i] = mapper.DockerImageToThrift(img)
	}
	thriftFeatures := make([]*center.ReleaseFeature, len(features))
	for i, f := range features {
		thriftFeatures[i] = mapper.FeatureToThrift(f)
	}
	thriftRel := mapper.ReleaseToThrift(rel, "", total, bugs, tasks, notes)

	content := rs.generateMarkdown(thriftRel, thriftItems, thriftBranches, thriftImages, thriftFeatures, version)

	snap, err := rs.snapshotStore.Create(releaseKeyword, version, content, total, bugs, tasks, notes)
	if err != nil {
		return nil, err
	}

	if err := rs.releaseStore.IncrementPublish(releaseKeyword); err != nil {
		return nil, err
	}

	if rs.notificationSvc != nil {
		project, _ := rs.projectStore.GetByID(rel.ProjectKeyword)
		projectName := ""
		if project != nil {
			projectName = project.Name
		}
		go rs.notificationSvc.NotifyReleasePublished(rel, items, projectName, version)
	}

	return mapper.SnapshotToThrift(snap), nil
}

func (rs *ReleaseService) ListSnapshots(releaseKeyword string) ([]*center.ReleaseSnapshot, error) {
	snaps, err := rs.snapshotStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ReleaseSnapshot, len(snaps))
	for i, snap := range snaps {
		result[i] = mapper.SnapshotToThrift(snap)
	}
	return result, nil
}

func (rs *ReleaseService) GetSnapshot(keyword string) (*center.ReleaseSnapshot, error) {
	snap, err := rs.snapshotStore.GetByID(keyword)
	if err != nil {
		return nil, err
	}
	if snap == nil {
		return nil, fmt.Errorf("snapshot not found")
	}
	return mapper.SnapshotToThrift(snap), nil
}

func (rs *ReleaseService) Export(releaseKeyword, snapshotKeyword, format string) (string, string, error) {
	var content string
	var version string

	if snapshotKeyword != "" {
		snap, err := rs.snapshotStore.GetByID(snapshotKeyword)
		if err != nil || snap == nil {
			return "", "", fmt.Errorf("snapshot not found")
		}
		content = snap.Content
		version = snap.Version
	} else {
		rel, err := rs.releaseStore.GetByID(releaseKeyword)
		if err != nil || rel == nil {
			return "", "", fmt.Errorf("release not found")
		}
		items, err := rs.itemStore.ListByRelease(releaseKeyword)
		if err != nil {
			return "", "", err
		}
		branches, err := rs.branchStore.ListByRelease(releaseKeyword)
		if err != nil {
			return "", "", err
		}
		dockerImages, err := rs.dockerImageStore.ListByRelease(releaseKeyword)
		if err != nil {
			return "", "", err
		}
		features, err := rs.featureStore.ListByRelease(releaseKeyword)
		if err != nil {
			return "", "", err
		}
		version = rel.Version

		thriftItems := make([]*center.ReleaseItem, len(items))
		for i, item := range items {
			thriftItems[i] = mapper.ItemToThrift(item)
		}
		thriftBranches := make([]*center.ReleaseBranch, len(branches))
		for i, b := range branches {
			thriftBranches[i] = mapper.ReleaseBranchToThrift(b)
		}
		thriftImages := make([]*center.DockerImage, len(dockerImages))
		for i, img := range dockerImages {
			thriftImages[i] = mapper.DockerImageToThrift(img)
		}
		thriftFeatures := make([]*center.ReleaseFeature, len(features))
		for i, f := range features {
			thriftFeatures[i] = mapper.FeatureToThrift(f)
		}
		thriftRel := mapper.ReleaseToThrift(rel, "", 0, 0, 0, 0)

		content = rs.generateMarkdown(thriftRel, thriftItems, thriftBranches, thriftImages, thriftFeatures, rel.Version)
	}

	if format == "html" {
		content = rs.markdownToHTML(content)
	}

	return content, version, nil
}

func (rs *ReleaseService) generateMarkdown(rel *center.Release, items []*center.ReleaseItem, branches []*center.ReleaseBranch, images []*center.DockerImage, features []*center.ReleaseFeature, version string) string {
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

	if len(features) > 0 {
		sb.WriteString(fmt.Sprintf("## 功能说明（%d）\n\n", len(features)))
		for i, f := range features {
			sb.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, f.Title))
			if f.Content != "" {
				sb.WriteString(f.Content)
				sb.WriteString("\n\n")
			}
		}
	}

	if len(branches) > 0 {
		sb.WriteString("## 分支信息\n\n")
		sb.WriteString("| 分支名 | 类型 | 父分支 |\n")
		sb.WriteString("|--------|------|--------|\n")
		for _, b := range branches {
			parent := b.ParentBranch
			if parent == "" {
				parent = "-"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", b.BranchName, b.BranchType, parent))
		}
		sb.WriteString("\n")
	}

	if len(images) > 0 {
		sb.WriteString("## Docker 镜像\n\n")
		sb.WriteString("| 镜像 | Commit | 来源 | 已测 |\n")
		sb.WriteString("|------|--------|------|------|\n")
		for _, img := range images {
			tested := "否"
			if img.Tested {
				tested = "是"
			}
			url := img.ImageUrl
			if url == "" {
				url = "-"
			}
			commit := img.CommitSha
			if commit == "" {
				commit = "-"
			}
			source := img.Source
			if source == "" {
				source = "-"
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", url, commit, source, tested))
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

func (rs *ReleaseService) GetRaw(keyword string) (*model.Release, error) {
	return rs.releaseStore.GetByID(keyword)
}

func (rs *ReleaseService) GetRawItems(releaseKeyword string) ([]*model.ReleaseItem, error) {
	return rs.itemStore.ListByRelease(releaseKeyword)
}

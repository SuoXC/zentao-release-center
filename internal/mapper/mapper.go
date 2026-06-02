package mapper

import (
	"fmt"
	"time"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/model"
)

func timeStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.DateTime)
}

func timeStrVal(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.DateTime)
}

func int32Ptr(v int) *int32 {
	r := int32(v)
	return &r
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func ProjectToThrift(p *model.Project) *center.Project {
	return &center.Project{
		ID:                p.Keyword,
		Name:              p.Name,
		Description:       p.Description,
		ZentaoProductId:   int32(p.ZentaoProductID),
		ZentaoProjectId:   int32(p.ZentaoProjectID),
		ZentaoProductName: p.ZentaoProductName,
		ZentaoProjectName: p.ZentaoProjectName,
		ZentaoServer:      p.ZentaoServer,
		Status:            p.Status,
		CreatedAt:         timeStrVal(p.CreatedAt),
		UpdatedAt:         timeStrVal(p.UpdatedAt),
	}
}

func ReleaseToThrift(r *model.Release, projectName string, itemCount, bugCount, taskCount, noteCount int) *center.Release {
	return &center.Release{
		ID:                r.Keyword,
		ProjectId:         r.ProjectKeyword,
		ProjectName:       projectName,
		Name:              r.Name,
		Version:           r.Version,
		Status:            r.Status,
		Summary:           r.Summary,
		ParentBranch:      r.ParentBranch,
		PublishCount:      int32(r.PublishCount),
		FirstPublishedAt:  timeStr(r.FirstPublishedAt),
		LastPublishedAt:   timeStr(r.LastPublishedAt),
		ItemCount:         int32(itemCount),
		BugCount:          int32(bugCount),
		TaskCount:         int32(taskCount),
		NoteCount:         int32(noteCount),
		CreatedAt:         timeStrVal(r.CreatedAt),
		UpdatedAt:         timeStrVal(r.UpdatedAt),
	}
}

func ItemToThrift(i *model.ReleaseItem) *center.ReleaseItem {
	item := &center.ReleaseItem{
		ID:           i.Keyword,
		ReleaseId:    i.ReleaseKeyword,
		ItemType:     i.ItemType,
		SortOrder:    int32(i.SortOrder),
		CreatedAt:    timeStrVal(i.CreatedAt),
		UpdatedAt:    timeStrVal(i.UpdatedAt),
	}
	if i.ZentaoID > 0 {
		item.ZentaoId = int32Ptr(i.ZentaoID)
	}
	if i.ZentaoType != "" {
		item.ZentaoType = strPtr(i.ZentaoType)
	}
	if i.Title != "" {
		item.Title = strPtr(i.Title)
	}
	if i.Severity != "" {
		item.Severity = strPtr(i.Severity)
	}
	if i.Priority != "" {
		item.Priority = strPtr(i.Priority)
	}
	if i.Status != "" {
		item.Status = strPtr(i.Status)
	}
	if i.AssignedTo != "" {
		item.AssignedTo = strPtr(i.AssignedTo)
	}
	if i.ResolvedBy != "" {
		item.ResolvedBy = strPtr(i.ResolvedBy)
	}
	if i.ZentaoURL != "" {
		item.ZentaoUrl = strPtr(i.ZentaoURL)
	}
	if i.Steps != "" {
		item.Steps = strPtr(i.Steps)
	}
	if i.NoteTitle != "" {
		item.NoteTitle = strPtr(i.NoteTitle)
	}
	if i.NoteContent != "" {
		item.NoteContent = strPtr(i.NoteContent)
	}
	return item
}

func SnapshotToThrift(s *model.ReleaseSnapshot) *center.ReleaseSnapshot {
	return &center.ReleaseSnapshot{
		ID:           s.Keyword,
		ReleaseId:    s.ReleaseKeyword,
		Version:      s.Version,
		Content:      s.Content,
		ItemCount:    int32(s.ItemCount),
		BugCount:     int32(s.BugCount),
		TaskCount:    int32(s.TaskCount),
		NoteCount:    int32(s.NoteCount),
		PublishedAt:  timeStrVal(s.PublishedAt),
	}
}

func DeploymentToThrift(d *model.ReleaseDeployment) *center.Deployment {
	return &center.Deployment{
		ID:           d.Keyword,
		ReleaseId:    d.ReleaseKeyword,
		ModuleName:   d.ModuleName,
		Address:      d.Address,
		Description:  d.Description,
		SortOrder:    int32(d.SortOrder),
		CreatedAt:    timeStrVal(d.CreatedAt),
		UpdatedAt:    timeStrVal(d.UpdatedAt),
	}
}

func ProjectRepoToThrift(r *model.ProjectRepo) *center.ProjectRepo {
	return &center.ProjectRepo{
		ID:              r.Keyword,
		ProjectId:       r.ProjectKeyword,
		GitlabProjectId: int32(r.GitlabProjectID),
		RepoUrl:         r.RepoURL,
		RepoName:        r.RepoName,
		DefaultBranch:   r.DefaultBranch,
		CreatedAt:       timeStrVal(r.CreatedAt),
	}
}

func ReleaseBranchToThrift(b *model.ReleaseBranch) *center.ReleaseBranch {
	return &center.ReleaseBranch{
		ID:              b.Keyword,
		ReleaseId:       b.ReleaseKeyword,
		RepoId:          b.RepoKeyword,
		BranchName:      b.BranchName,
		BranchType:      b.BranchType,
		ParentBranch:    b.ParentBranch,
		GitlabBranchUrl: b.GitlabBranchURL,
		Description:     b.Description,
		CreatedAt:       timeStrVal(b.CreatedAt),
	}
}

func DockerImageToThrift(d *model.DockerImage) *center.DockerImage {
	return &center.DockerImage{
		ID:             d.Keyword,
		ReleaseId:      d.ReleaseKeyword,
		RepoId:         d.RepoKeyword,
		ImageName:      d.ImageName,
		ImageTag:       d.ImageTag,
		ImageDigest:    d.ImageDigest,
		Registry:       d.Registry,
		CiPipelineId:   int32(d.CIPipelineID),
		CiPipelineUrl:  d.CIPipelineURL,
		Branch:         d.Branch,
		CommitSha:      d.CommitSHA,
		CommitMessage:  d.CommitMessage,
		Source:         d.Source,
		CreatedAt:      timeStrVal(d.CreatedAt),
	}
}

func FeatureToThrift(f *model.ReleaseFeature) *center.ReleaseFeature {
	return &center.ReleaseFeature{
		ID:          f.Keyword,
		ReleaseId:   f.ReleaseKeyword,
		Title:       f.Title,
		Content:     f.Content,
		SortOrder:   int32(f.SortOrder),
		CreatedAt:   timeStrVal(f.CreatedAt),
		UpdatedAt:   timeStrVal(f.UpdatedAt),
	}
}

func PtrOrDefault[T any](v *T, def T) T {
	if v == nil {
		return def
	}
	return *v
}

func FormatDockerImageName(registry, name, tag string) string {
	if registry != "" {
		return fmt.Sprintf("%s/%s:%s", registry, name, tag)
	}
	return fmt.Sprintf("%s:%s", name, tag)
}

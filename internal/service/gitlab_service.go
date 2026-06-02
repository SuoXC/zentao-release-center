package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/gitlab"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"gorm.io/gorm"
)

type GitLabService struct {
	repoStore   *store.RepoStore
	branchStore *store.BranchStore
	releaseStore *store.ReleaseStore
	gitlabClient *gitlab.Client
}

func NewGitLabService(db *gorm.DB, gc *gitlab.Client) *GitLabService {
	return &GitLabService{
		repoStore:    store.NewRepoStore(db),
		branchStore:  store.NewBranchStore(db),
		releaseStore: store.NewReleaseStore(db),
		gitlabClient: gc,
	}
}

func (gs *GitLabService) SearchProjects(query string) ([]*center.GitlabProject, error) {
	projects, err := gs.gitlabClient.SearchProjects(query)
	if err != nil {
		return nil, err
	}
	result := make([]*center.GitlabProject, len(projects))
	for i, p := range projects {
		result[i] = &center.GitlabProject{
			ID:                int32(p.ID),
			Name:              p.Name,
			NameWithNamespace: p.NameWithNamespace,
			PathWithNamespace: p.PathWithNamespace,
			WebUrl:            p.WebURL,
			HttpUrlToRepo:    p.HTTPURLToRepo,
			DefaultBranch:     p.DefaultBranch,
		}
	}
	return result, nil
}

func (gs *GitLabService) ListBranches(gitlabProjectID int) ([]*center.GitlabBranch, error) {
	branches, err := gs.gitlabClient.ListBranches(gitlabProjectID)
	if err != nil {
		return nil, err
	}
	result := make([]*center.GitlabBranch, len(branches))
	for i, b := range branches {
		result[i] = &center.GitlabBranch{
			Name:        b.Name,
			IsDefault:   b.Default,
			IsProtected: b.Protected,
			WebUrl:      b.WebURL,
		}
	}
	return result, nil
}

func (gs *GitLabService) AddRepo(req *center.AddRepoReq) (*center.ProjectRepo, error) {
	if req.ProjectId == "" {
		return nil, fmt.Errorf("projectId is required")
	}
	repo, err := gs.repoStore.Create(req.ProjectId, int(req.GitlabProjectId), req.RepoUrl, req.RepoName, req.GetDefaultBranch())
	if err != nil {
		return nil, err
	}
	return mapper.ProjectRepoToThrift(repo), nil
}

func (gs *GitLabService) ListRepos(projectKeyword string) ([]*center.ProjectRepo, error) {
	repos, err := gs.repoStore.ListByProject(projectKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ProjectRepo, len(repos))
	for i, r := range repos {
		result[i] = mapper.ProjectRepoToThrift(r)
	}
	return result, nil
}

func (gs *GitLabService) DeleteRepo(keyword string) error {
	return gs.repoStore.Delete(keyword)
}

func (gs *GitLabService) CreateReleaseBranch(req *center.CreateReleaseBranchReq) (*center.ReleaseBranch, error) {
	if req.ReleaseId == "" {
		return nil, fmt.Errorf("releaseId is required")
	}
	if req.RepoId == "" {
		return nil, fmt.Errorf("repoId is required")
	}

	repo, err := gs.repoStore.GetByKeyword(req.RepoId)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repo not found")
	}

	rel, err := gs.releaseStore.GetByID(req.ReleaseId)
	if err != nil || rel == nil {
		return nil, fmt.Errorf("release not found")
	}

	branchName := req.GetBranchName()
	if branchName == "" {
		branchName = fmt.Sprintf("release/%s", rel.Version)
	}

	parentBranch := req.GetParentBranch()
	if parentBranch == "" {
		parentBranch = repo.DefaultBranch
	}

	if err := gs.gitlabClient.CreateBranch(repo.GitlabProjectID, branchName, parentBranch); err != nil {
		return nil, fmt.Errorf("create gitlab branch: %w", err)
	}

	gitlabBranchURL := fmt.Sprintf("%s/-/tree/%s", repo.RepoURL, branchName)

	b, err := gs.branchStore.Create(req.ReleaseId, req.RepoId, branchName, "release", parentBranch, gitlabBranchURL)
	if err != nil {
		return nil, err
	}
	return mapper.ReleaseBranchToThrift(b), nil
}

func (gs *GitLabService) CreateFeatureBranch(req *center.CreateFeatureBranchReq) (*center.ReleaseBranch, error) {
	if req.ReleaseId == "" {
		return nil, fmt.Errorf("releaseId is required")
	}
	if req.RepoId == "" {
		return nil, fmt.Errorf("repoId is required")
	}
	if req.BranchName == "" {
		return nil, fmt.Errorf("branchName is required")
	}

	repo, err := gs.repoStore.GetByKeyword(req.RepoId)
	if err != nil {
		return nil, err
	}
	if repo == nil {
		return nil, fmt.Errorf("repo not found")
	}

	parentBranch := req.GetParentBranch()
	if parentBranch == "" {
		releaseBranches, _ := gs.branchStore.ListByRelease(req.ReleaseId)
		for _, rb := range releaseBranches {
			if rb.BranchType == "release" && rb.RepoKeyword == req.RepoId {
				parentBranch = rb.BranchName
				break
			}
		}
	}
	if parentBranch == "" {
		parentBranch = repo.DefaultBranch
	}

	if err := gs.gitlabClient.CreateBranch(repo.GitlabProjectID, req.BranchName, parentBranch); err != nil {
		return nil, fmt.Errorf("create gitlab branch: %w", err)
	}

	gitlabBranchURL := fmt.Sprintf("%s/-/tree/%s", repo.RepoURL, req.BranchName)

	b, err := gs.branchStore.Create(req.ReleaseId, req.RepoId, req.BranchName, "feature", parentBranch, gitlabBranchURL)
	if err != nil {
		return nil, err
	}
	return mapper.ReleaseBranchToThrift(b), nil
}

func (gs *GitLabService) ListBranchesByRelease(releaseKeyword string) ([]*center.ReleaseBranch, error) {
	branches, err := gs.branchStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ReleaseBranch, len(branches))
	for i, b := range branches {
		result[i] = mapper.ReleaseBranchToThrift(b)
	}
	return result, nil
}

func (gs *GitLabService) DeleteBranch(keyword string) error {
	return gs.branchStore.Delete(keyword)
}

func (gs *GitLabService) UpdateBranch(req *center.UpdateBranchReq) error {
	if req.ID == "" {
		return fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.IsSetDescription() {
		fields["description"] = req.Description
	}
	return gs.branchStore.Update(req.ID, fields)
}

func (gs *GitLabService) FindRepoByGitlabProject(gitlabProjectID int) (*model.ProjectRepo, error) {
	var repo model.ProjectRepo
	if err := gs.repoStore.DB().Where("gitlab_project_id = ?", gitlabProjectID).First(&repo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &repo, nil
}

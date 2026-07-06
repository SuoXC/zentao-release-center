package center

import "time"

type GitlabWebhookEvent struct {
	ObjectKind       string            `json:"object_kind"`
	ObjectAttributes GitlabPipelineAttr `json:"object_attributes"`
	Project          GitlabWebhookProject `json:"project"`
	Commit           GitlabWebhookCommit  `json:"commit"`
}

type GitlabPipelineAttr struct {
	ID         int       `json:"id"`
	IID        int       `json:"iid"`
	ProjectID  int       `json:"project_id"`
	Status     string    `json:"status"`
	Source     string    `json:"source"`
	Ref        string    `json:"ref"`
	SHA        string    `json:"sha"`
	WebURL     string    `json:"web_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt time.Time `json:"finished_at"`
	Duration   int       `json:"duration"`
}

type GitlabWebhookProject struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
}

type GitlabWebhookCommit struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Author  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
}

type ProjectRepo struct {
	ID              string `thrift:"id,1" form:"id" json:"id" query:"id"`
	ProjectId       string `thrift:"projectId,2" form:"projectId" json:"projectId" query:"projectId"`
	GitlabProjectId int32  `thrift:"gitlabProjectId,3" form:"gitlabProjectId" json:"gitlabProjectId" query:"gitlabProjectId"`
	RepoUrl         string `thrift:"repoUrl,4" form:"repoUrl" json:"repoUrl" query:"repoUrl"`
	RepoName        string `thrift:"repoName,5" form:"repoName" json:"repoName" query:"repoName"`
	DefaultBranch   string `thrift:"defaultBranch,6" form:"defaultBranch" json:"defaultBranch" query:"defaultBranch"`
	CreatedAt       string `thrift:"createdAt,7" form:"createdAt" json:"createdAt" query:"createdAt"`
}

func (p *ProjectRepo) GetID() string           { return p.ID }
func (p *ProjectRepo) GetProjectId() string     { return p.ProjectId }
func (p *ProjectRepo) GetGitlabProjectId() int32 { return p.GitlabProjectId }
func (p *ProjectRepo) GetRepoUrl() string       { return p.RepoUrl }
func (p *ProjectRepo) GetRepoName() string      { return p.RepoName }
func (p *ProjectRepo) GetDefaultBranch() string  { return p.DefaultBranch }
func (p *ProjectRepo) GetCreatedAt() string     { return p.CreatedAt }

type AddRepoReq struct {
	ProjectId       string  `thrift:"projectId,1" form:"projectId" json:"projectId" query:"projectId"`
	GitlabProjectId int32   `thrift:"gitlabProjectId,2" form:"gitlabProjectId" json:"gitlabProjectId" query:"gitlabProjectId"`
	RepoUrl         string  `thrift:"repoUrl,3" form:"repoUrl" json:"repoUrl" query:"repoUrl"`
	RepoName        string  `thrift:"repoName,4" form:"repoName" json:"repoName" query:"repoName"`
	DefaultBranch   *string `thrift:"defaultBranch,5,optional" form:"defaultBranch" json:"defaultBranch,omitempty" query:"defaultBranch"`
}

func (p *AddRepoReq) GetProjectId() string            { return p.ProjectId }
func (p *AddRepoReq) GetGitlabProjectId() int32        { return p.GitlabProjectId }
func (p *AddRepoReq) GetRepoUrl() string               { return p.RepoUrl }
func (p *AddRepoReq) GetRepoName() string              { return p.RepoName }
func (p *AddRepoReq) GetDefaultBranch() string {
	if p.DefaultBranch != nil {
		return *p.DefaultBranch
	}
	return "main"
}
func (p *AddRepoReq) IsSetDefaultBranch() bool { return p.DefaultBranch != nil }

type DeleteRepoReq struct {
	ID string `thrift:"id,1" form:"id" json:"id" query:"id"`
}

func (p *DeleteRepoReq) GetID() string { return p.ID }

type ListReposReq struct {
	ProjectId string `thrift:"projectId,1" form:"projectId" json:"projectId" query:"projectId"`
}

func (p *ListReposReq) GetProjectId() string { return p.ProjectId }

type RepoResp struct {
	Base *BaseResp      `thrift:"base,1" form:"base" json:"base" query:"base"`
	Data *ProjectRepo   `thrift:"data,2,optional" form:"data" json:"data,omitempty" query:"data"`
}

func (p *RepoResp) GetBase() *BaseResp      { return p.Base }
func (p *RepoResp) GetData() *ProjectRepo    { return p.Data }

type RepoListResp struct {
	Base *BaseResp       `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*ProjectRepo  `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *RepoListResp) GetBase() *BaseResp       { return p.Base }
func (p *RepoListResp) GetList() []*ProjectRepo   { return p.List }

type ReleaseBranch struct {
	ID              string `thrift:"id,1" form:"id" json:"id" query:"id"`
	ReleaseId       string `thrift:"releaseId,2" form:"releaseId" json:"releaseId" query:"releaseId"`
	RepoId          string `thrift:"repoId,3" form:"repoId" json:"repoId" query:"repoId"`
	BranchName      string `thrift:"branchName,4" form:"branchName" json:"branchName" query:"branchName"`
	BranchType      string `thrift:"branchType,5" form:"branchType" json:"branchType" query:"branchType"`
	ParentBranch    string `thrift:"parentBranch,6" form:"parentBranch" json:"parentBranch" query:"parentBranch"`
	GitlabBranchUrl string `thrift:"gitlabBranchUrl,7" form:"gitlabBranchUrl" json:"gitlabBranchUrl" query:"gitlabBranchUrl"`
	Description     string `thrift:"description,8" form:"description" json:"description" query:"description"`
	CreatedAt       string `thrift:"createdAt,9" form:"createdAt" json:"createdAt" query:"createdAt"`
}

func (p *ReleaseBranch) GetID() string              { return p.ID }
func (p *ReleaseBranch) GetReleaseId() string        { return p.ReleaseId }
func (p *ReleaseBranch) GetRepoId() string           { return p.RepoId }
func (p *ReleaseBranch) GetBranchName() string       { return p.BranchName }
func (p *ReleaseBranch) GetBranchType() string       { return p.BranchType }
func (p *ReleaseBranch) GetParentBranch() string     { return p.ParentBranch }
func (p *ReleaseBranch) GetGitlabBranchUrl() string  { return p.GitlabBranchUrl }
func (p *ReleaseBranch) GetDescription() string      { return p.Description }
func (p *ReleaseBranch) GetCreatedAt() string        { return p.CreatedAt }

type CreateReleaseBranchReq struct {
	ReleaseId   string  `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
	RepoId      string  `thrift:"repoId,2" form:"repoId" json:"repoId" query:"repoId"`
	BranchName  *string `thrift:"branchName,3,optional" form:"branchName" json:"branchName,omitempty" query:"branchName"`
	ParentBranch *string `thrift:"parentBranch,4,optional" form:"parentBranch" json:"parentBranch,omitempty" query:"parentBranch"`
}

func (p *CreateReleaseBranchReq) GetReleaseId() string { return p.ReleaseId }
func (p *CreateReleaseBranchReq) GetRepoId() string    { return p.RepoId }
func (p *CreateReleaseBranchReq) GetBranchName() string {
	if p.BranchName != nil {
		return *p.BranchName
	}
	return ""
}
func (p *CreateReleaseBranchReq) IsSetBranchName() bool { return p.BranchName != nil }
func (p *CreateReleaseBranchReq) GetParentBranch() string {
	if p.ParentBranch != nil {
		return *p.ParentBranch
	}
	return ""
}
func (p *CreateReleaseBranchReq) IsSetParentBranch() bool { return p.ParentBranch != nil }

type CreateFeatureBranchReq struct {
	ReleaseId    string  `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
	RepoId       string  `thrift:"repoId,2" form:"repoId" json:"repoId" query:"repoId"`
	BranchName   string  `thrift:"branchName,3" form:"branchName" json:"branchName" query:"branchName"`
	ParentBranch *string `thrift:"parentBranch,4,optional" form:"parentBranch" json:"parentBranch,omitempty" query:"parentBranch"`
}

func (p *CreateFeatureBranchReq) GetReleaseId() string    { return p.ReleaseId }
func (p *CreateFeatureBranchReq) GetRepoId() string       { return p.RepoId }
func (p *CreateFeatureBranchReq) GetBranchName() string   { return p.BranchName }
func (p *CreateFeatureBranchReq) GetParentBranch() string {
	if p.ParentBranch != nil {
		return *p.ParentBranch
	}
	return ""
}
func (p *CreateFeatureBranchReq) IsSetParentBranch() bool { return p.ParentBranch != nil }

type ListBranchesReq struct {
	ReleaseId string `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
}

func (p *ListBranchesReq) GetReleaseId() string { return p.ReleaseId }

type DeleteBranchReq struct {
	ID string `thrift:"id,1" form:"id" json:"id" query:"id"`
}

func (p *DeleteBranchReq) GetID() string { return p.ID }

type UpdateBranchReq struct {
	ID          string  `thrift:"id,1" form:"id" json:"id"`
	Description *string `thrift:"description,2,optional" form:"description" json:"description,omitempty"`
}

func (p *UpdateBranchReq) GetID() string { return p.ID }
func (p *UpdateBranchReq) GetDescription() string {
	if p.Description != nil {
		return *p.Description
	}
	return ""
}
func (p *UpdateBranchReq) IsSetDescription() bool { return p.Description != nil }

type BranchResp struct {
	Base *BaseResp       `thrift:"base,1" form:"base" json:"base" query:"base"`
	Data *ReleaseBranch  `thrift:"data,2,optional" form:"data" json:"data,omitempty" query:"data"`
}

func (p *BranchResp) GetBase() *BaseResp       { return p.Base }
func (p *BranchResp) GetData() *ReleaseBranch   { return p.Data }

type BranchListResp struct {
	Base *BaseResp        `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*ReleaseBranch `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *BranchListResp) GetBase() *BaseResp        { return p.Base }
func (p *BranchListResp) GetList() []*ReleaseBranch  { return p.List }

type DockerImage struct {
	ID            string `thrift:"id,1" form:"id" json:"id" query:"id"`
	ReleaseId     string `thrift:"releaseId,2" form:"releaseId" json:"releaseId" query:"releaseId"`
	RepoId        string `thrift:"repoId,3" form:"repoId" json:"repoId" query:"repoId"`
	ImageUrl      string `thrift:"imageUrl,4" form:"imageUrl" json:"imageUrl" query:"imageUrl"`
	ImageDigest   string `thrift:"imageDigest,5" form:"imageDigest" json:"imageDigest" query:"imageDigest"`
	CiPipelineId  int32  `thrift:"ciPipelineId,6" form:"ciPipelineId" json:"ciPipelineId" query:"ciPipelineId"`
	CiPipelineUrl string `thrift:"ciPipelineUrl,7" form:"ciPipelineUrl" json:"ciPipelineUrl" query:"ciPipelineUrl"`
	CommitSha     string `thrift:"commitSha,8" form:"commitSha" json:"commitSha" query:"commitSha"`
	CommitMessage string `thrift:"commitMessage,9" form:"commitMessage" json:"commitMessage" query:"commitMessage"`
	Source        string `thrift:"source,10" form:"source" json:"source" query:"source"`
	Tested        bool   `thrift:"tested,11" form:"tested" json:"tested" query:"tested"`
	CreatedAt     string `thrift:"createdAt,12" form:"createdAt" json:"createdAt" query:"createdAt"`
}

func (p *DockerImage) GetID() string             { return p.ID }
func (p *DockerImage) GetReleaseId() string     { return p.ReleaseId }
func (p *DockerImage) GetRepoId() string        { return p.RepoId }
func (p *DockerImage) GetImageUrl() string       { return p.ImageUrl }
func (p *DockerImage) GetImageDigest() string    { return p.ImageDigest }
func (p *DockerImage) GetCiPipelineId() int32   { return p.CiPipelineId }
func (p *DockerImage) GetCiPipelineUrl() string { return p.CiPipelineUrl }
func (p *DockerImage) GetCommitSha() string     { return p.CommitSha }
func (p *DockerImage) GetCommitMessage() string { return p.CommitMessage }
func (p *DockerImage) GetSource() string        { return p.Source }
func (p *DockerImage) GetTested() bool          { return p.Tested }
func (p *DockerImage) GetCreatedAt() string     { return p.CreatedAt }

type AddDockerImageReq struct {
	ReleaseId     string  `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
	RepoId        *string `thrift:"repoId,2,optional" form:"repoId" json:"repoId,omitempty" query:"repoId"`
	ImageUrl      string  `thrift:"imageUrl,3" form:"imageUrl" json:"imageUrl" query:"imageUrl"`
	ImageDigest   *string `thrift:"imageDigest,4,optional" form:"imageDigest" json:"imageDigest,omitempty" query:"imageDigest"`
	CommitSha     *string `thrift:"commitSha,5,optional" form:"commitSha" json:"commitSha,omitempty" query:"commitSha"`
	CommitMessage *string `thrift:"commitMessage,6,optional" form:"commitMessage" json:"commitMessage,omitempty" query:"commitMessage"`
}

func (p *AddDockerImageReq) GetReleaseId() string { return p.ReleaseId }
func (p *AddDockerImageReq) GetRepoId() string {
	if p.RepoId != nil {
		return *p.RepoId
	}
	return ""
}
func (p *AddDockerImageReq) GetImageUrl() string { return p.ImageUrl }
func (p *AddDockerImageReq) GetImageDigest() string {
	if p.ImageDigest != nil {
		return *p.ImageDigest
	}
	return ""
}
func (p *AddDockerImageReq) GetCommitSha() string {
	if p.CommitSha != nil {
		return *p.CommitSha
	}
	return ""
}
func (p *AddDockerImageReq) GetCommitMessage() string {
	if p.CommitMessage != nil {
		return *p.CommitMessage
	}
	return ""
}
func (p *AddDockerImageReq) IsSetRepoId() bool        { return p.RepoId != nil }
func (p *AddDockerImageReq) IsSetImageDigest() bool    { return p.ImageDigest != nil }
func (p *AddDockerImageReq) IsSetCommitSha() bool      { return p.CommitSha != nil }
func (p *AddDockerImageReq) IsSetCommitMessage() bool { return p.CommitMessage != nil }

type UpdateDockerImageReq struct {
	ID            string  `thrift:"id,1" form:"id" json:"id"`
	ImageUrl      *string `thrift:"imageUrl,2,optional" form:"imageUrl" json:"imageUrl,omitempty"`
	Tested        *bool   `thrift:"tested,3,optional" form:"tested" json:"tested,omitempty"`
	CommitSha     *string `thrift:"commitSha,4,optional" form:"commitSha" json:"commitSha,omitempty"`
	CommitMessage *string `thrift:"commitMessage,5,optional" form:"commitMessage" json:"commitMessage,omitempty"`
}

func (p *UpdateDockerImageReq) GetID() string { return p.ID }
func (p *UpdateDockerImageReq) GetImageUrl() string {
	if p.ImageUrl != nil {
		return *p.ImageUrl
	}
	return ""
}
func (p *UpdateDockerImageReq) GetTested() bool {
	if p.Tested != nil {
		return *p.Tested
	}
	return false
}
func (p *UpdateDockerImageReq) GetCommitSha() string {
	if p.CommitSha != nil {
		return *p.CommitSha
	}
	return ""
}
func (p *UpdateDockerImageReq) GetCommitMessage() string {
	if p.CommitMessage != nil {
		return *p.CommitMessage
	}
	return ""
}
func (p *UpdateDockerImageReq) IsSetImageUrl() bool      { return p.ImageUrl != nil }
func (p *UpdateDockerImageReq) IsSetTested() bool        { return p.Tested != nil }
func (p *UpdateDockerImageReq) IsSetCommitSha() bool     { return p.CommitSha != nil }
func (p *UpdateDockerImageReq) IsSetCommitMessage() bool { return p.CommitMessage != nil }

type DeleteDockerImageReq struct {
	ID string `thrift:"id,1" form:"id" json:"id" query:"id"`
}

func (p *DeleteDockerImageReq) GetID() string { return p.ID }

type ListDockerImagesReq struct {
	ReleaseId string `thrift:"releaseId,1" form:"releaseId" json:"releaseId" query:"releaseId"`
}

func (p *ListDockerImagesReq) GetReleaseId() string { return p.ReleaseId }

type DockerImageResp struct {
	Base *BaseResp      `thrift:"base,1" form:"base" json:"base" query:"base"`
	Data *DockerImage   `thrift:"data,2,optional" form:"data" json:"data,omitempty" query:"data"`
}

func (p *DockerImageResp) GetBase() *BaseResp      { return p.Base }
func (p *DockerImageResp) GetData() *DockerImage    { return p.Data }

type DockerImageListResp struct {
	Base *BaseResp       `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*DockerImage  `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *DockerImageListResp) GetBase() *BaseResp       { return p.Base }
func (p *DockerImageListResp) GetList() []*DockerImage   { return p.List }

type GitlabProject struct {
	ID                int32  `thrift:"id,1" form:"id" json:"id" query:"id"`
	Name              string `thrift:"name,2" form:"name" json:"name" query:"name"`
	NameWithNamespace string `thrift:"nameWithNamespace,3" form:"nameWithNamespace" json:"nameWithNamespace" query:"nameWithNamespace"`
	PathWithNamespace string `thrift:"pathWithNamespace,4" form:"pathWithNamespace" json:"pathWithNamespace" query:"pathWithNamespace"`
	WebUrl            string `thrift:"webUrl,5" form:"webUrl" json:"webUrl" query:"webUrl"`
	HttpUrlToRepo    string `thrift:"httpUrlToRepo,6" form:"httpUrlToRepo" json:"httpUrlToRepo" query:"httpUrlToRepo"`
	DefaultBranch     string `thrift:"defaultBranch,7" form:"defaultBranch" json:"defaultBranch" query:"defaultBranch"`
}

func (p *GitlabProject) GetID() int32                { return p.ID }
func (p *GitlabProject) GetName() string              { return p.Name }
func (p *GitlabProject) GetNameWithNamespace() string { return p.NameWithNamespace }
func (p *GitlabProject) GetPathWithNamespace() string { return p.PathWithNamespace }
func (p *GitlabProject) GetWebUrl() string            { return p.WebUrl }
func (p *GitlabProject) GetHttpUrlToRepo() string    { return p.HttpUrlToRepo }
func (p *GitlabProject) GetDefaultBranch() string     { return p.DefaultBranch }

type GitlabBranch struct {
	Name        string `thrift:"name,1" form:"name" json:"name" query:"name"`
	IsDefault   bool   `thrift:"isDefault,2" form:"isDefault" json:"isDefault" query:"isDefault"`
	IsProtected bool   `thrift:"isProtected,3" form:"isProtected" json:"isProtected" query:"isProtected"`
	WebUrl      string `thrift:"webUrl,4" form:"webUrl" json:"webUrl" query:"webUrl"`
}

func (p *GitlabBranch) GetName() string        { return p.Name }
func (p *GitlabBranch) GetIsDefault() bool     { return p.IsDefault }
func (p *GitlabBranch) GetIsProtected() bool   { return p.IsProtected }
func (p *GitlabBranch) GetWebUrl() string      { return p.WebUrl }

type SearchGitlabProjectsReq struct {
	Query string `thrift:"query,1" form:"query" json:"query" query:"query"`
}

func (p *SearchGitlabProjectsReq) GetQuery() string { return p.Query }

type ListGitlabBranchesReq struct {
	GitlabProjectId int32 `thrift:"gitlabProjectId,1" form:"gitlabProjectId" json:"gitlabProjectId" query:"gitlabProjectId"`
}

func (p *ListGitlabBranchesReq) GetGitlabProjectId() int32 { return p.GitlabProjectId }

type GitlabProjectListResp struct {
	Base *BaseResp        `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*GitlabProject `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *GitlabProjectListResp) GetBase() *BaseResp        { return p.Base }
func (p *GitlabProjectListResp) GetList() []*GitlabProject  { return p.List }

type GitlabBranchListResp struct {
	Base *BaseResp       `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*GitlabBranch `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *GitlabBranchListResp) GetBase() *BaseResp       { return p.Base }
func (p *GitlabBranchListResp) GetList() []*GitlabBranch  { return p.List }

type CIBuildReq struct {
	ReleaseId     string  `thrift:"releaseId,1" form:"releaseId" json:"releaseId"`
	RepoId        *string `thrift:"repoId,2,optional" form:"repoId" json:"repoId,omitempty"`
	ImageUrl      string  `thrift:"imageUrl,3" form:"imageUrl" json:"imageUrl"`
	ImageDigest   *string `thrift:"imageDigest,4,optional" form:"imageDigest" json:"imageDigest,omitempty"`
	CiPipelineId  *int32  `thrift:"ciPipelineId,5,optional" form:"ciPipelineId" json:"ciPipelineId,omitempty"`
	CiPipelineUrl *string `thrift:"ciPipelineUrl,6,optional" form:"ciPipelineUrl" json:"ciPipelineUrl,omitempty"`
	CommitSha     *string `thrift:"commitSha,7,optional" form:"commitSha" json:"commitSha,omitempty"`
	CommitMessage *string `thrift:"commitMessage,8,optional" form:"commitMessage" json:"commitMessage,omitempty"`
}

func (p *CIBuildReq) GetReleaseId() string { return p.ReleaseId }
func (p *CIBuildReq) GetRepoId() string {
	if p.RepoId != nil {
		return *p.RepoId
	}
	return ""
}
func (p *CIBuildReq) GetImageUrl() string { return p.ImageUrl }
func (p *CIBuildReq) GetImageDigest() string {
	if p.ImageDigest != nil {
		return *p.ImageDigest
	}
	return ""
}
func (p *CIBuildReq) GetCiPipelineId() int32 {
	if p.CiPipelineId != nil {
		return *p.CiPipelineId
	}
	return 0
}
func (p *CIBuildReq) GetCiPipelineUrl() string {
	if p.CiPipelineUrl != nil {
		return *p.CiPipelineUrl
	}
	return ""
}
func (p *CIBuildReq) GetCommitSha() string {
	if p.CommitSha != nil {
		return *p.CommitSha
	}
	return ""
}
func (p *CIBuildReq) GetCommitMessage() string {
	if p.CommitMessage != nil {
		return *p.CommitMessage
	}
	return ""
}
func (p *CIBuildReq) IsSetRepoId() bool        { return p.RepoId != nil }
func (p *CIBuildReq) IsSetImageDigest() bool    { return p.ImageDigest != nil }
func (p *CIBuildReq) IsSetCiPipelineId() bool   { return p.CiPipelineId != nil }
func (p *CIBuildReq) IsSetCiPipelineUrl() bool  { return p.CiPipelineUrl != nil }
func (p *CIBuildReq) IsSetCommitSha() bool      { return p.CommitSha != nil }
func (p *CIBuildReq) IsSetCommitMessage() bool { return p.CommitMessage != nil }

type DockerImagePoolItem struct {
	ID            string `thrift:"id,1" form:"id" json:"id"`
	ImageUrl      string `thrift:"imageUrl,2" form:"imageUrl" json:"imageUrl"`
	ImageDigest   string `thrift:"imageDigest,3" form:"imageDigest" json:"imageDigest"`
	CommitSha     string `thrift:"commitSha,4" form:"commitSha" json:"commitSha"`
	CommitMessage string `thrift:"commitMessage,5" form:"commitMessage" json:"commitMessage"`
	CiPipelineId  int32  `thrift:"ciPipelineId,6" form:"ciPipelineId" json:"ciPipelineId"`
	CiPipelineUrl string `thrift:"ciPipelineUrl,7" form:"ciPipelineUrl" json:"ciPipelineUrl"`
	CreatedAt     string `thrift:"createdAt,8" form:"createdAt" json:"createdAt"`
}

type ListDockerImagePoolReq struct {
	GitlabProjectId *int32 `thrift:"gitlabProjectId,1,optional" form:"gitlabProjectId" json:"gitlabProjectId,omitempty" query:"gitlabProjectId"`
}

func (p *ListDockerImagePoolReq) GetGitlabProjectId() int32 {
	if p.GitlabProjectId != nil {
		return *p.GitlabProjectId
	}
	return 0
}
func (p *ListDockerImagePoolReq) IsSetGitlabProjectId() bool { return p.GitlabProjectId != nil }

type DockerImagePoolListResp struct {
	Base *BaseResp               `thrift:"base,1" form:"base" json:"base" query:"base"`
	List []*DockerImagePoolItem  `thrift:"list,2" form:"list" json:"list" query:"list"`
}

func (p *DockerImagePoolListResp) GetBase() *BaseResp               { return p.Base }
func (p *DockerImagePoolListResp) GetList() []*DockerImagePoolItem  { return p.List }

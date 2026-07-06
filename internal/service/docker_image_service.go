package service

import (
	"fmt"
	"strings"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/gitlab"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"gorm.io/gorm"
)

type DockerImageService struct {
	imageStore   *store.DockerImageStore
	poolStore    *store.DockerImagePoolStore
	repoStore    *store.RepoStore
	branchStore  *store.BranchStore
	releaseStore *store.ReleaseStore
	db           *gorm.DB
	gitlabClient *gitlab.Client
}

func NewDockerImageService(db *gorm.DB, gc *gitlab.Client) *DockerImageService {
	return &DockerImageService{
		imageStore:   store.NewDockerImageStore(db),
		poolStore:    store.NewDockerImagePoolStore(db),
		repoStore:    store.NewRepoStore(db),
		branchStore:  store.NewBranchStore(db),
		releaseStore: store.NewReleaseStore(db),
		db:           db,
		gitlabClient: gc,
	}
}

func (ds *DockerImageService) AddManual(req *center.AddDockerImageReq) (*center.DockerImage, error) {
	if req.ReleaseId == "" {
		return nil, fmt.Errorf("releaseId is required")
	}
	if req.ImageUrl == "" {
		return nil, fmt.Errorf("imageUrl is required")
	}

	image, err := ds.imageStore.Create(
		req.ReleaseId,
		req.GetRepoId(),
		req.ImageUrl,
		req.GetImageDigest(),
		req.GetCommitSha(),
		req.GetCommitMessage(),
		"manual",
		0, "",
	)
	if err != nil {
		return nil, err
	}
	return mapper.DockerImageToThrift(image), nil
}

func (ds *DockerImageService) AddFromWebhook(event *gitlab.PipelineEvent) error {
	if event.ObjectAttributes.Status != "success" {
		return nil
	}

	gitlabProjectID := event.Project.ID

	var repo model.ProjectRepo
	if err := ds.db.Where("gitlab_project_id = ?", gitlabProjectID).First(&repo).Error; err != nil {
		// 没有关联仓库，仍然保存到池中
		ds.poolStore.Create(gitlabProjectID, fmt.Sprintf("pipeline-%d:no-repo", event.ObjectAttributes.ID),
			"", event.Commit.ID, event.Commit.Message,
			event.ObjectAttributes.ID, event.ObjectAttributes.WebURL)
		return nil
	}

	commitSHA := event.Commit.ID
	commitMessage := event.Commit.Message
	pipelineID := event.ObjectAttributes.ID
	pipelineURL := event.ObjectAttributes.WebURL
	imageURL := buildImageURL(event, repo.RepoName)
	if imageURL == "" {
		imageURL = fmt.Sprintf("%s-pipeline-%d:%s", repo.RepoName, pipelineID, commitSHA)
	}

	// 保存到全局镜像池
	ds.poolStore.Create(gitlabProjectID, imageURL, "", commitSHA, commitMessage, pipelineID, pipelineURL)

	// 查找所有 draft 状态的发布单并添加镜像
	var releases []model.Release
	ds.db.Where("status = ? AND project_keyword = ?", "draft", repo.ProjectKeyword).Find(&releases)

	for _, rel := range releases {
		ds.imageStore.Create(
			rel.Keyword,
			repo.Keyword,
			imageURL,
			"",
			commitSHA,
			commitMessage,
			"webhook",
			pipelineID,
			pipelineURL,
		)
	}
	return nil
}

func (ds *DockerImageService) AddFromCIBuild(req *center.CIBuildReq) error {
	if req.ImageUrl == "" {
		return fmt.Errorf("imageUrl is required")
	}

	// 尝试找到仓库
	var repo model.ProjectRepo
	ds.db.Where("gitlab_project_id = ?", 0).First(&repo) // 占位，下面用 release 查

	gitlabProjectID := 0
	if req.ReleaseId != "" {
		var release model.Release
		if err := ds.db.Where("keyword = ?", req.ReleaseId).First(&release).Error; err == nil {
			// 找到项目关联的仓库
			var repos []model.ProjectRepo
			ds.db.Where("project_keyword = ?", release.ProjectKeyword).Find(&repos)
			if len(repos) > 0 {
				gitlabProjectID = repos[0].GitlabProjectID
			}
		}
	}

	// 保存到全局镜像池
	ds.poolStore.Create(gitlabProjectID, req.ImageUrl, req.GetImageDigest(),
		req.GetCommitSha(), req.GetCommitMessage(),
		int(req.GetCiPipelineId()), req.GetCiPipelineUrl())

	// 如果指定了发布单，直接添加
	if req.ReleaseId != "" {
		_, err := ds.imageStore.Create(
			req.ReleaseId,
			req.GetRepoId(),
			req.ImageUrl,
			req.GetImageDigest(),
			req.GetCommitSha(),
			req.GetCommitMessage(),
			"ci",
			int(req.GetCiPipelineId()),
			req.GetCiPipelineUrl(),
		)
		return err
	}
	return nil
}

func (ds *DockerImageService) List(releaseKeyword string) ([]*center.DockerImage, error) {
	images, err := ds.imageStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.DockerImage, len(images))
	for i, img := range images {
		result[i] = mapper.DockerImageToThrift(img)
	}
	return result, nil
}

func (ds *DockerImageService) ListPool(gitlabProjectID int) ([]*center.DockerImagePoolItem, error) {
	var images []*model.DockerImagePool
	var err error
	if gitlabProjectID > 0 {
		images, err = ds.poolStore.ListByProject(gitlabProjectID)
	} else {
		images, err = ds.poolStore.ListAll()
	}
	if err != nil {
		return nil, err
	}
	result := make([]*center.DockerImagePoolItem, len(images))
	for i, img := range images {
		result[i] = &center.DockerImagePoolItem{
			ID:            img.Keyword,
			ImageUrl:      img.ImageURL,
			ImageDigest:   img.ImageDigest,
			CommitSha:     img.CommitSHA,
			CommitMessage: img.CommitMessage,
			CiPipelineId:  int32(img.CIPipelineID),
			CiPipelineUrl: img.CIPipelineURL,
			CreatedAt:     img.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return result, nil
}

func (ds *DockerImageService) Delete(keyword string) error {
	return ds.imageStore.Delete(keyword)
}

func (ds *DockerImageService) Update(req *center.UpdateDockerImageReq) (*center.DockerImage, error) {
	if req.ID == "" {
		return nil, fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.IsSetImageUrl() {
		fields["image_url"] = req.GetImageUrl()
	}
	if req.IsSetTested() {
		fields["tested"] = req.GetTested()
	}
	if req.IsSetCommitSha() {
		fields["commit_sha"] = req.GetCommitSha()
	}
	if req.IsSetCommitMessage() {
		fields["commit_message"] = req.GetCommitMessage()
	}
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}
	image, err := ds.imageStore.Update(req.ID, fields)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, fmt.Errorf("image not found")
	}
	return mapper.DockerImageToThrift(image), nil
}

// buildImageURL 从 webhook 中提取完整镜像 URL（若 CI 携带了 IMAGE_NAME/CI_REGISTRY_IMAGE 等变量）。
func buildImageURL(event *gitlab.PipelineEvent, fallbackName string) string {
	registry := ""
	name := ""
	tag := ""
	if event.Variables != nil {
		for _, v := range event.Variables {
			switch v.Key {
			case "CI_REGISTRY_IMAGE":
				parts := strings.SplitN(v.Value, "/", 2)
				if len(parts) == 2 {
					registry = parts[0]
					name = parts[1]
				} else {
					name = v.Value
				}
			case "IMAGE_NAME":
				name = v.Value
			case "CI_REGISTRY":
				registry = v.Value
			case "IMAGE_TAG":
				tag = v.Value
			}
		}
	}
	if name == "" {
		name = fallbackName
		if name == "" {
			return ""
		}
	}
	if tag == "" {
		tag = event.ObjectAttributes.Ref
		if tag == "" {
			tag = "latest"
		}
	}
	name = strings.TrimPrefix(name, "/")
	if registry != "" {
		return fmt.Sprintf("%s/%s:%s", registry, name, tag)
	}
	return fmt.Sprintf("%s:%s", name, tag)
}

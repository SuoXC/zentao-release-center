package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type DockerImagePoolStore struct {
	db *gorm.DB
}

func NewDockerImagePoolStore(db *gorm.DB) *DockerImagePoolStore {
	return &DockerImagePoolStore{db: db}
}

func (ps *DockerImagePoolStore) Create(gitlabProjectID int, imageName, imageTag, imageDigest, registry string, ciPipelineID int, ciPipelineURL, branch, commitSHA, commitMessage string) (*model.DockerImagePool, error) {
	// 检查是否已存在（去重）
	var existing model.DockerImagePool
	err := ps.db.Where("gitlab_project_id = ? AND image_name = ? AND image_tag = ? AND commit_sha = ?",
		gitlabProjectID, imageName, imageTag, commitSHA).First(&existing).Error
	if err == nil {
		return &existing, nil // 已存在，返回现有记录
	}

	p := &model.DockerImagePool{
		Keyword:         uuid.New().String(),
		GitlabProjectID: gitlabProjectID,
		ImageName:       imageName,
		ImageTag:        imageTag,
		ImageDigest:     imageDigest,
		Registry:        registry,
		CIPipelineID:    ciPipelineID,
		CIPipelineURL:   ciPipelineURL,
		Branch:          branch,
		CommitSHA:       commitSHA,
		CommitMessage:   commitMessage,
	}
	if err := ps.db.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (ps *DockerImagePoolStore) ListByProject(gitlabProjectID int) ([]*model.DockerImagePool, error) {
	var images []*model.DockerImagePool
	if err := ps.db.Where("gitlab_project_id = ?", gitlabProjectID).Order("created_at DESC").Limit(50).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (ps *DockerImagePoolStore) ListAll() ([]*model.DockerImagePool, error) {
	var images []*model.DockerImagePool
	if err := ps.db.Order("created_at DESC").Limit(100).Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

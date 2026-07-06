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

func (ps *DockerImagePoolStore) Create(gitlabProjectID int, imageURL, imageDigest, commitSHA, commitMessage string, ciPipelineID int, ciPipelineURL string) (*model.DockerImagePool, error) {
	// 检查是否已存在（去重）
	var existing model.DockerImagePool
	err := ps.db.Where("gitlab_project_id = ? AND image_url = ? AND commit_sha = ?",
		gitlabProjectID, imageURL, commitSHA).First(&existing).Error
	if err == nil {
		return &existing, nil // 已存在，返回现有记录
	}

	p := &model.DockerImagePool{
		Keyword:         uuid.New().String(),
		GitlabProjectID: gitlabProjectID,
		ImageURL:        imageURL,
		ImageDigest:     imageDigest,
		CIPipelineID:    ciPipelineID,
		CIPipelineURL:   ciPipelineURL,
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

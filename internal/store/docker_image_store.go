package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type DockerImageStore struct {
	db *gorm.DB
}

func NewDockerImageStore(db *gorm.DB) *DockerImageStore {
	return &DockerImageStore{db: db}
}

func (ds *DockerImageStore) Create(releaseKeyword, repoKeyword, imageName, imageTag, imageDigest, registry string, ciPipelineID int, ciPipelineURL, branch, commitSHA, commitMessage, source string) (*model.DockerImage, error) {
	d := &model.DockerImage{
		Keyword:        uuid.New().String(),
		ReleaseKeyword: releaseKeyword,
		RepoKeyword:    repoKeyword,
		ImageName:      imageName,
		ImageTag:       imageTag,
		ImageDigest:    imageDigest,
		Registry:       registry,
		CIPipelineID:   ciPipelineID,
		CIPipelineURL:  ciPipelineURL,
		Branch:         branch,
		CommitSHA:      commitSHA,
		CommitMessage:  commitMessage,
		Source:         source,
	}
	if err := ds.db.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (ds *DockerImageStore) GetByKeyword(keyword string) (*model.DockerImage, error) {
	var d model.DockerImage
	if err := ds.db.Where("keyword = ?", keyword).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (ds *DockerImageStore) ListByRelease(releaseKeyword string) ([]*model.DockerImage, error) {
	var images []*model.DockerImage
	if err := ds.db.Where("release_keyword = ?", releaseKeyword).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, err
	}
	return images, nil
}

func (ds *DockerImageStore) Delete(keyword string) error {
	return ds.db.Where("keyword = ?", keyword).Delete(&model.DockerImage{}).Error
}

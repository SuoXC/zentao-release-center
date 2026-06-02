package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type RepoStore struct {
	db *gorm.DB
}

func NewRepoStore(db *gorm.DB) *RepoStore {
	return &RepoStore{db: db}
}

func (rs *RepoStore) DB() *gorm.DB {
	return rs.db
}

func (rs *RepoStore) Create(projectKeyword string, gitlabProjectID int, repoURL, repoName, defaultBranch string) (*model.ProjectRepo, error) {
	r := &model.ProjectRepo{
		Keyword:         uuid.New().String(),
		ProjectKeyword:  projectKeyword,
		GitlabProjectID: gitlabProjectID,
		RepoURL:         repoURL,
		RepoName:        repoName,
		DefaultBranch:   defaultBranch,
	}
	if err := rs.db.Create(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *RepoStore) GetByKeyword(keyword string) (*model.ProjectRepo, error) {
	var r model.ProjectRepo
	if err := rs.db.Where("keyword = ?", keyword).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (rs *RepoStore) ListByProject(projectKeyword string) ([]*model.ProjectRepo, error) {
	var repos []*model.ProjectRepo
	if err := rs.db.Where("project_keyword = ?", projectKeyword).Order("created_at DESC").Find(&repos).Error; err != nil {
		return nil, err
	}
	return repos, nil
}

func (rs *RepoStore) Delete(keyword string) error {
	return rs.db.Where("keyword = ?", keyword).Delete(&model.ProjectRepo{}).Error
}

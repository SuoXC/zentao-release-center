package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type BranchStore struct {
	db *gorm.DB
}

func NewBranchStore(db *gorm.DB) *BranchStore {
	return &BranchStore{db: db}
}

func (bs *BranchStore) Create(releaseKeyword, repoKeyword, branchName, branchType, parentBranch, gitlabBranchURL, description string) (*model.ReleaseBranch, error) {
	b := &model.ReleaseBranch{
		Keyword:         uuid.New().String(),
		ReleaseKeyword:  releaseKeyword,
		RepoKeyword:     repoKeyword,
		BranchName:      branchName,
		BranchType:      branchType,
		ParentBranch:    parentBranch,
		GitlabBranchURL: gitlabBranchURL,
		Description:     description,
	}
	if err := bs.db.Create(b).Error; err != nil {
		return nil, err
	}
	return b, nil
}

func (bs *BranchStore) GetByKeyword(keyword string) (*model.ReleaseBranch, error) {
	var b model.ReleaseBranch
	if err := bs.db.Where("keyword = ?", keyword).First(&b).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (bs *BranchStore) ListByRelease(releaseKeyword string) ([]*model.ReleaseBranch, error) {
	var branches []*model.ReleaseBranch
	if err := bs.db.Where("release_keyword = ?", releaseKeyword).Order("created_at DESC").Find(&branches).Error; err != nil {
		return nil, err
	}
	return branches, nil
}

func (bs *BranchStore) Delete(keyword string) error {
	return bs.db.Where("keyword = ?", keyword).Delete(&model.ReleaseBranch{}).Error
}

func (bs *BranchStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return bs.db.Model(&model.ReleaseBranch{}).Where("keyword = ?", keyword).Updates(fields).Error
}

func (bs *BranchStore) FindByBranchAndRepo(releaseKeyword, repoKeyword, branchName string) (*model.ReleaseBranch, error) {
	var b model.ReleaseBranch
	if err := bs.db.Where("release_keyword = ? AND repo_keyword = ? AND branch_name = ?", releaseKeyword, repoKeyword, branchName).First(&b).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

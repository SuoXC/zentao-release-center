package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type DeploymentStore struct {
	db *gorm.DB
}

func NewDeploymentStore(db *gorm.DB) *DeploymentStore {
	return &DeploymentStore{db: db}
}

func (ds *DeploymentStore) Create(releaseKeyword, moduleName, address, description string) (*model.ReleaseDeployment, error) {
	var maxOrder int
	ds.db.Model(&model.ReleaseDeployment{}).Where("release_keyword = ?", releaseKeyword).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	d := &model.ReleaseDeployment{
		Keyword:        uuid.New().String(),
		ReleaseKeyword: releaseKeyword,
		ModuleName:     moduleName,
		Address:        address,
		Description:    description,
		SortOrder:      maxOrder + 1,
	}
	if err := ds.db.Create(d).Error; err != nil {
		return nil, err
	}
	return d, nil
}

func (ds *DeploymentStore) GetByID(keyword string) (*model.ReleaseDeployment, error) {
	var d model.ReleaseDeployment
	if err := ds.db.Where("keyword = ?", keyword).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (ds *DeploymentStore) ListByRelease(releaseKeyword string) ([]*model.ReleaseDeployment, error) {
	var deps []*model.ReleaseDeployment
	if err := ds.db.Where("release_keyword = ?", releaseKeyword).Order("sort_order ASC, created_at ASC").Find(&deps).Error; err != nil {
		return nil, err
	}
	return deps, nil
}

var deploymentAllowedFields = map[string]bool{
	"module_name": true, "address": true, "description": true, "sort_order": true,
}

func (ds *DeploymentStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	allowed := make(map[string]interface{})
	for k, v := range fields {
		if deploymentAllowedFields[k] {
			allowed[k] = v
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	return ds.db.Model(&model.ReleaseDeployment{}).Where("keyword = ?", keyword).Updates(allowed).Error
}

func (ds *DeploymentStore) Delete(keyword string) error {
	return ds.db.Where("keyword = ?", keyword).Delete(&model.ReleaseDeployment{}).Error
}

func (ds *DeploymentStore) CountByRelease(releaseKeyword string) (int, error) {
	var count int64
	err := ds.db.Model(&model.ReleaseDeployment{}).Where("release_keyword = ?", releaseKeyword).Count(&count).Error
	return int(count), err
}

func (ds *DeploymentStore) Reorder(items []struct {
	Keyword   string
	SortOrder int
}) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.ReleaseDeployment{}).Where("keyword = ?", item.Keyword).Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

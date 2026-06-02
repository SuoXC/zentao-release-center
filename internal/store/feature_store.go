package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type FeatureStore struct {
	db *gorm.DB
}

func NewFeatureStore(db *gorm.DB) *FeatureStore {
	return &FeatureStore{db: db}
}

func (fs *FeatureStore) Create(releaseKeyword, title, content string) (*model.ReleaseFeature, error) {
	var maxOrder int
	fs.db.Model(&model.ReleaseFeature{}).Where("release_keyword = ?", releaseKeyword).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	f := &model.ReleaseFeature{
		Keyword:        uuid.New().String(),
		ReleaseKeyword: releaseKeyword,
		Title:          title,
		Content:        content,
		SortOrder:      maxOrder + 1,
	}
	if err := fs.db.Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (fs *FeatureStore) GetByID(keyword string) (*model.ReleaseFeature, error) {
	var f model.ReleaseFeature
	if err := fs.db.Where("keyword = ?", keyword).First(&f).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &f, nil
}

func (fs *FeatureStore) ListByRelease(releaseKeyword string) ([]*model.ReleaseFeature, error) {
	var features []*model.ReleaseFeature
	if err := fs.db.Where("release_keyword = ?", releaseKeyword).Order("sort_order ASC, created_at ASC").Find(&features).Error; err != nil {
		return nil, err
	}
	return features, nil
}

var featureAllowedFields = map[string]bool{
	"title": true, "content": true, "sort_order": true,
}

func (fs *FeatureStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	allowed := make(map[string]interface{})
	for k, v := range fields {
		if featureAllowedFields[k] {
			allowed[k] = v
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	return fs.db.Model(&model.ReleaseFeature{}).Where("keyword = ?", keyword).Updates(allowed).Error
}

func (fs *FeatureStore) Delete(keyword string) error {
	return fs.db.Where("keyword = ?", keyword).Delete(&model.ReleaseFeature{}).Error
}

func (fs *FeatureStore) CountByRelease(releaseKeyword string) (int, error) {
	var count int64
	err := fs.db.Model(&model.ReleaseFeature{}).Where("release_keyword = ?", releaseKeyword).Count(&count).Error
	return int(count), err
}

func (fs *FeatureStore) Reorder(items []struct {
	Keyword   string
	SortOrder int
}) error {
	return fs.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.ReleaseFeature{}).Where("keyword = ?", item.Keyword).Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

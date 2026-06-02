package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"gorm.io/gorm"
)

type FeatureService struct {
	featureStore *store.FeatureStore
}

func NewFeatureService(db *gorm.DB) *FeatureService {
	return &FeatureService{
		featureStore: store.NewFeatureStore(db),
	}
}

func (fs *FeatureService) Add(req *center.AddFeatureReq) (*center.ReleaseFeature, error) {
	if req.ReleaseId == "" {
		return nil, fmt.Errorf("releaseId is required")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	feature, err := fs.featureStore.Create(req.ReleaseId, req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return mapper.FeatureToThrift(feature), nil
}

func (fs *FeatureService) Get(keyword string) (*center.ReleaseFeature, error) {
	feature, err := fs.featureStore.GetByID(keyword)
	if err != nil {
		return nil, err
	}
	if feature == nil {
		return nil, fmt.Errorf("feature not found")
	}
	return mapper.FeatureToThrift(feature), nil
}

func (fs *FeatureService) List(releaseKeyword string) ([]*center.ReleaseFeature, error) {
	features, err := fs.featureStore.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.ReleaseFeature, len(features))
	for i, f := range features {
		result[i] = mapper.FeatureToThrift(f)
	}
	return result, nil
}

func (fs *FeatureService) Update(req *center.UpdateFeatureReq) error {
	if req.ID == "" {
		return fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Content != nil {
		fields["content"] = *req.Content
	}
	return fs.featureStore.Update(req.ID, fields)
}

func (fs *FeatureService) Delete(keyword string) error {
	return fs.featureStore.Delete(keyword)
}

func (fs *FeatureService) Reorder(releaseKeyword string, items []center.SortItem) error {
	sortItems := make([]struct{ Keyword string; SortOrder int }, len(items))
	for i, item := range items {
		sortItems[i] = struct{ Keyword string; SortOrder int }{item.ID, int(item.SortOrder)}
	}
	return fs.featureStore.Reorder(sortItems)
}

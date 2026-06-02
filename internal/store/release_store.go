package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type ReleaseStore struct {
	db *gorm.DB
}

func NewReleaseStore(db *gorm.DB) *ReleaseStore {
	return &ReleaseStore{db: db}
}

func (rs *ReleaseStore) Create(projectKeyword, name, version, summary, parentBranch string) (*model.Release, error) {
	r := &model.Release{
		Keyword:        uuid.New().String(),
		ProjectKeyword: projectKeyword,
		Name:           name,
		Version:        version,
		Status:         "draft",
		Summary:        summary,
		ParentBranch:   parentBranch,
	}
	if err := rs.db.Create(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

func (rs *ReleaseStore) GetByID(keyword string) (*model.Release, error) {
	var r model.Release
	if err := rs.db.Where("keyword = ?", keyword).First(&r).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

func (rs *ReleaseStore) List(projectKeyword, status string, page, pageSize int) ([]*model.Release, int, error) {
	var releases []*model.Release
	query := rs.db.Model(&model.Release{})

	if projectKeyword != "" {
		query = query.Where("project_keyword = ?", projectKeyword)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&releases).Error; err != nil {
		return nil, 0, err
	}

	return releases, int(total), nil
}

var releaseAllowedFields = map[string]bool{
	"name": true, "version": true, "summary": true, "status": true, "parent_branch": true,
}

func (rs *ReleaseStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	allowed := make(map[string]interface{})
	for k, v := range fields {
		if releaseAllowedFields[k] {
			allowed[k] = v
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	return rs.db.Model(&model.Release{}).Where("keyword = ?", keyword).Updates(allowed).Error
}

func (rs *ReleaseStore) Delete(keyword string) error {
	return rs.db.Where("keyword = ?", keyword).Delete(&model.Release{}).Error
}

func (rs *ReleaseStore) IncrementPublish(keyword string) error {
	now := time.Now()
	return rs.db.Model(&model.Release{}).Where("keyword = ?", keyword).Updates(map[string]interface{}{
		"publish_count":     gorm.Expr("publish_count + 1"),
		"status":            "published",
		"last_published_at": now,
		"first_published_at": gorm.Expr("COALESCE(first_published_at, ?)", now),
	}).Error
}

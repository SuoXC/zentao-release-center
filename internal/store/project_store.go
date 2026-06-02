package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type ProjectStore struct {
	db *gorm.DB
}

func NewProjectStore(db *gorm.DB) *ProjectStore {
	return &ProjectStore{db: db}
}

func (ps *ProjectStore) Create(name, description string, zentaoProductID, zentaoProjectID int, zentaoProductName, zentaoProjectName string) (*model.Project, error) {
	p := &model.Project{
		Keyword:           uuid.New().String(),
		Name:              name,
		Description:       description,
		ZentaoProductID:   zentaoProductID,
		ZentaoProjectID:   zentaoProjectID,
		ZentaoProductName: zentaoProductName,
		ZentaoProjectName: zentaoProjectName,
		Status:            "active",
	}
	if err := ps.db.Create(p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (ps *ProjectStore) GetByID(keyword string) (*model.Project, error) {
	var p model.Project
	if err := ps.db.Where("keyword = ?", keyword).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (ps *ProjectStore) List(status string, page, pageSize int) ([]*model.Project, int, error) {
	var projects []*model.Project
	query := ps.db.Model(&model.Project{})

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

	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&projects).Error; err != nil {
		return nil, 0, err
	}

	return projects, int(total), nil
}

func (ps *ProjectStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	return ps.db.Model(&model.Project{}).Where("keyword = ?", keyword).Updates(fields).Error
}

func (ps *ProjectStore) Delete(keyword string) error {
	return ps.db.Where("keyword = ?", keyword).Delete(&model.Project{}).Error
}

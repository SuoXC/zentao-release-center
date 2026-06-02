package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"gorm.io/gorm"
)

type ProjectService struct {
	store *store.ProjectStore
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{store: store.NewProjectStore(db)}
}

func (ps *ProjectService) Create(req *center.CreateProjectReq) (*center.Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	p, err := ps.store.Create(req.Name, req.GetDescription(), int(req.GetZentaoProductId()), int(req.GetZentaoProjectId()), req.GetZentaoProductName(), req.GetZentaoProjectName())
	if err != nil {
		return nil, err
	}
	return mapper.ProjectToThrift(p), nil
}

func (ps *ProjectService) Get(keyword string) (*center.Project, error) {
	p, err := ps.store.GetByID(keyword)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project not found")
	}
	return mapper.ProjectToThrift(p), nil
}

func (ps *ProjectService) List(status string, page, pageSize int) ([]*center.Project, int, error) {
	projects, total, err := ps.store.List(status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*center.Project, len(projects))
	for i, p := range projects {
		result[i] = mapper.ProjectToThrift(p)
	}
	return result, total, nil
}

func (ps *ProjectService) Update(req *center.UpdateProjectReq) error {
	if req.ID == "" {
		return fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.IsSetName() {
		fields["name"] = req.Name
	}
	if req.IsSetDescription() {
		fields["description"] = req.Description
	}
	if req.IsSetZentaoProductId() {
		fields["zentao_product_id"] = req.ZentaoProductId
	}
	if req.IsSetZentaoProjectId() {
		fields["zentao_project_id"] = req.ZentaoProjectId
	}
	if req.IsSetZentaoProductName() {
		fields["zentao_product_name"] = req.ZentaoProductName
	}
	if req.IsSetZentaoProjectName() {
		fields["zentao_project_name"] = req.ZentaoProjectName
	}
	if req.IsSetStatus() {
		fields["status"] = req.Status
	}
	return ps.store.Update(req.ID, fields)
}

func (ps *ProjectService) Delete(keyword string) error {
	return ps.store.Delete(keyword)
}

func (ps *ProjectService) GetRaw(keyword string) (*model.Project, error) {
	return ps.store.GetByID(keyword)
}

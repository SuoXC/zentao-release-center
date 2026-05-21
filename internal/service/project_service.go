package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/store"
)

type ProjectService struct {
	projectStore *store.ProjectStore
}

func NewProjectService(s *store.Store) *ProjectService {
	return &ProjectService{projectStore: store.NewProjectStore(s)}
}

func (ps *ProjectService) Create(req *center.CreateProjectReq) (*center.Project, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	return ps.projectStore.Create(req.Name, req.GetDescription(), int(req.GetZentaoProductId()), int(req.GetZentaoProjectId()), req.GetZentaoProductName(), req.GetZentaoProjectName())
}

func (ps *ProjectService) Get(id string) (*center.Project, error) {
	p, err := ps.projectStore.GetByID(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("project not found")
	}
	return p, nil
}

func (ps *ProjectService) List(status string, page, pageSize int) ([]*center.Project, int, error) {
	return ps.projectStore.List(status, page, pageSize)
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
	return ps.projectStore.Update(req.ID, fields)
}

func (ps *ProjectService) Delete(id string) error {
	return ps.projectStore.Delete(id)
}

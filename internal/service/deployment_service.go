package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/mapper"
	"github.com/yi-nology/zentao-release-center/internal/store"
	"gorm.io/gorm"
)

type DeploymentService struct {
	store *store.DeploymentStore
}

func NewDeploymentService(db *gorm.DB) *DeploymentService {
	return &DeploymentService{store: store.NewDeploymentStore(db)}
}

func (ds *DeploymentService) Add(req *center.AddDeploymentReq) (*center.Deployment, error) {
	if req.ReleaseId == "" {
		return nil, fmt.Errorf("releaseId is required")
	}
	if req.ModuleName == "" {
		return nil, fmt.Errorf("moduleName is required")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("address is required")
	}
	d, err := ds.store.Create(req.ReleaseId, req.ModuleName, req.Address, req.GetDescription())
	if err != nil {
		return nil, err
	}
	return mapper.DeploymentToThrift(d), nil
}

func (ds *DeploymentService) Get(keyword string) (*center.Deployment, error) {
	d, err := ds.store.GetByID(keyword)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, fmt.Errorf("deployment not found")
	}
	return mapper.DeploymentToThrift(d), nil
}

func (ds *DeploymentService) List(releaseKeyword string) ([]*center.Deployment, error) {
	deps, err := ds.store.ListByRelease(releaseKeyword)
	if err != nil {
		return nil, err
	}
	result := make([]*center.Deployment, len(deps))
	for i, d := range deps {
		result[i] = mapper.DeploymentToThrift(d)
	}
	return result, nil
}

func (ds *DeploymentService) Update(req *center.UpdateDeploymentReq) error {
	if req.ID == "" {
		return fmt.Errorf("id is required")
	}
	fields := map[string]interface{}{}
	if req.IsSetModuleName() {
		fields["module_name"] = req.ModuleName
	}
	if req.IsSetAddress() {
		fields["address"] = req.Address
	}
	if req.IsSetDescription() {
		fields["description"] = req.Description
	}
	return ds.store.Update(req.ID, fields)
}

func (ds *DeploymentService) Delete(keyword string) error {
	return ds.store.Delete(keyword)
}

func (ds *DeploymentService) ListByRelease(releaseKeyword string) ([]*center.Deployment, error) {
	return ds.List(releaseKeyword)
}

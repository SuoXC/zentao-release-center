package service

import (
	"fmt"

	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
	"github.com/yi-nology/zentao-release-center/internal/store"
)

type DeploymentService struct {
	store *store.DeploymentStore
}

func NewDeploymentService(s *store.Store) *DeploymentService {
	return &DeploymentService{store: store.NewDeploymentStore(s)}
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
	return ds.store.Create(req.ReleaseId, req.ModuleName, req.Address, req.GetDescription())
}

func (ds *DeploymentService) Get(id string) (*center.Deployment, error) {
	d, err := ds.store.GetByID(id)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, fmt.Errorf("deployment not found")
	}
	return d, nil
}

func (ds *DeploymentService) List(releaseID string) ([]*center.Deployment, error) {
	return ds.store.ListByRelease(releaseID)
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

func (ds *DeploymentService) Delete(id string) error {
	return ds.store.Delete(id)
}

func (ds *DeploymentService) ListByRelease(releaseID string) ([]*center.Deployment, error) {
	return ds.store.ListByRelease(releaseID)
}

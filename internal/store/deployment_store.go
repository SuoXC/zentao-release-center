package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
)

type DeploymentStore struct {
	db *sql.DB
}

func NewDeploymentStore(s *Store) *DeploymentStore {
	return &DeploymentStore{db: s.DB()}
}

func (ds *DeploymentStore) Create(releaseID, moduleName, address, description string) (*center.Deployment, error) {
	now := time.Now().Format(time.DateTime)
	id := uuid.New().String()

	var maxOrder sql.NullInt32
	ds.db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM release_deployments WHERE release_id = ?", releaseID).Scan(&maxOrder)
	sortOrder := int(maxOrder.Int32) + 1

	_, err := ds.db.Exec(`INSERT INTO release_deployments (id, release_id, module_name, address, description, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		id, releaseID, moduleName, address, description, sortOrder, now, now)
	if err != nil {
		return nil, err
	}
	return ds.GetByID(id)
}

func (ds *DeploymentStore) GetByID(id string) (*center.Deployment, error) {
	d := &center.Deployment{}
	var desc sql.NullString
	err := ds.db.QueryRow(`SELECT id, release_id, module_name, address, description, sort_order, created_at, updated_at FROM release_deployments WHERE id = ?`, id).
		Scan(&d.ID, &d.ReleaseId, &d.ModuleName, &d.Address, &desc, &d.SortOrder, &d.CreatedAt, &d.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.Description = desc.String
	return d, nil
}

func (ds *DeploymentStore) ListByRelease(releaseID string) ([]*center.Deployment, error) {
	rows, err := ds.db.Query(`SELECT id, release_id, module_name, address, description, sort_order, created_at, updated_at FROM release_deployments WHERE release_id = ? ORDER BY sort_order ASC, created_at ASC`, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*center.Deployment
	for rows.Next() {
		d := &center.Deployment{}
		var desc sql.NullString
		if err := rows.Scan(&d.ID, &d.ReleaseId, &d.ModuleName, &d.Address, &desc, &d.SortOrder, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		d.Description = desc.String
		list = append(list, d)
	}
	return list, nil
}

var deploymentAllowedFields = map[string]bool{
	"module_name": true, "address": true, "description": true, "sort_order": true,
}

func (ds *DeploymentStore) Update(id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	fields["updated_at"] = time.Now().Format(time.DateTime)

	setClauses := ""
	args := []interface{}{}
	for k, v := range fields {
		if !deploymentAllowedFields[k] && k != "updated_at" {
			continue
		}
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += k + " = ?"
		args = append(args, v)
	}
	if setClauses == "" {
		return nil
	}
	args = append(args, id)
	_, err := ds.db.Exec("UPDATE release_deployments SET "+setClauses+" WHERE id = ?", args...)
	return err
}

func (ds *DeploymentStore) Delete(id string) error {
	_, err := ds.db.Exec("DELETE FROM release_deployments WHERE id = ?", id)
	return err
}

func (ds *DeploymentStore) CountByRelease(releaseID string) (int, error) {
	var count int
	err := ds.db.QueryRow("SELECT COUNT(*) FROM release_deployments WHERE release_id = ?", releaseID).Scan(&count)
	return count, err
}

func (ds *DeploymentStore) Reorder(items []struct {
	ID        string
	SortOrder int
}) error {
	tx, err := ds.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range items {
		if _, err := tx.Exec("UPDATE release_deployments SET sort_order = ? WHERE id = ?", item.SortOrder, item.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}


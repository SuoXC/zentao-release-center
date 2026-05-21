package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
)

type ProjectStore struct {
	db *sql.DB
}

func NewProjectStore(s *Store) *ProjectStore {
	return &ProjectStore{db: s.DB()}
}

func (ps *ProjectStore) Create(name, description string, zentaoProductID, zentaoProjectID int, zentaoProductName, zentaoProjectName string) (*center.Project, error) {
	now := time.Now().Format(time.DateTime)
	id := uuid.New().String()
	_, err := ps.db.Exec(`INSERT INTO projects (id, name, description, zentao_product_id, zentao_project_id, zentao_product_name, zentao_project_name, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'active', ?, ?)`,
		id, name, description, zentaoProductID, zentaoProjectID, zentaoProductName, zentaoProjectName, now, now)
	if err != nil {
		return nil, err
	}
	return ps.GetByID(id)
}

func (ps *ProjectStore) GetByID(id string) (*center.Project, error) {
	p := &center.Project{}
	var desc, zentaoServer sql.NullString
	var zentaoProdID, zentaoProjID sql.NullInt32
	var zentaoProdName, zentaoProjName sql.NullString

	err := ps.db.QueryRow(`SELECT id, name, description, zentao_product_id, zentao_project_id, zentao_product_name, zentao_project_name, zentao_server, status, created_at, updated_at FROM projects WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &desc, &zentaoProdID, &zentaoProjID, &zentaoProdName, &zentaoProjName, &zentaoServer, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Description = desc.String
	p.ZentaoProductId = zentaoProdID.Int32
	p.ZentaoProjectId = zentaoProjID.Int32
	p.ZentaoProductName = zentaoProdName.String
	p.ZentaoProjectName = zentaoProjName.String
	p.ZentaoServer = zentaoServer.String
	return p, nil
}

func (ps *ProjectStore) List(status string, page, pageSize int) ([]*center.Project, int, error) {
	var args []interface{}
	query := "SELECT id, name, description, zentao_product_id, zentao_project_id, zentao_product_name, zentao_project_name, zentao_server, status, created_at, updated_at FROM projects"
	countQuery := "SELECT COUNT(*) FROM projects"

	if status != "" {
		query += " WHERE status = ?"
		countQuery += " WHERE status = ?"
		args = append(args, status)
	}

	var total int
	if err := ps.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", pageSize, offset)

	rows, err := ps.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*center.Project
	for rows.Next() {
		p := &center.Project{}
		var desc, zentaoServer sql.NullString
		var zentaoProdID, zentaoProjID sql.NullInt32
		var zentaoProdName, zentaoProjName sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &desc, &zentaoProdID, &zentaoProjID, &zentaoProdName, &zentaoProjName, &zentaoServer, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		p.Description = desc.String
		p.ZentaoProductId = zentaoProdID.Int32
		p.ZentaoProjectId = zentaoProjID.Int32
		p.ZentaoProductName = zentaoProdName.String
		p.ZentaoProjectName = zentaoProjName.String
		p.ZentaoServer = zentaoServer.String
		list = append(list, p)
	}
	return list, total, nil
}

func (ps *ProjectStore) Update(id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	fields["updated_at"] = time.Now().Format(time.DateTime)

	setClauses := ""
	args := []interface{}{}
	for k, v := range fields {
		if setClauses != "" {
			setClauses += ", "
		}
		setClauses += k + " = ?"
		args = append(args, v)
	}
	args = append(args, id)
	_, err := ps.db.Exec("UPDATE projects SET "+setClauses+" WHERE id = ?", args...)
	return err
}

func (ps *ProjectStore) Delete(id string) error {
	_, err := ps.db.Exec("DELETE FROM projects WHERE id = ?", id)
	return err
}

package store

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
)

type ReleaseStore struct {
	db *sql.DB
}

func NewReleaseStore(s *Store) *ReleaseStore {
	return &ReleaseStore{db: s.DB()}
}

func (rs *ReleaseStore) Create(projectID, name, version, summary string) (*center.Release, error) {
	now := time.Now().Format(time.DateTime)
	id := uuid.New().String()
	_, err := rs.db.Exec(`INSERT INTO releases (id, project_id, name, version, status, summary, publish_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, 'draft', ?, 0, ?, ?)`,
		id, projectID, name, version, summary, now, now)
	if err != nil {
		return nil, err
	}
	return rs.GetByID(id)
}

func scanRelease(row *sql.Row) (*center.Release, error) {
	r := &center.Release{}
	var version, summary, firstPub, lastPub sql.NullString
	var publishCount sql.NullInt32
	err := row.Scan(&r.ID, &r.ProjectId, &r.Name, &version, &r.Status, &summary, &publishCount, &firstPub, &lastPub, &r.CreatedAt, &r.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.Version = version.String
	r.Summary = summary.String
	r.PublishCount = publishCount.Int32
	r.FirstPublishedAt = firstPub.String
	r.LastPublishedAt = lastPub.String
	return r, nil
}

func (rs *ReleaseStore) GetByID(id string) (*center.Release, error) {
	return scanRelease(rs.db.QueryRow(`SELECT id, project_id, name, version, status, summary, publish_count, first_published_at, last_published_at, created_at, updated_at FROM releases WHERE id = ?`, id))
}

func (rs *ReleaseStore) List(projectID, status string, page, pageSize int) ([]*center.Release, int, error) {
	var args []interface{}
	query := "SELECT id, project_id, name, version, status, summary, publish_count, first_published_at, last_published_at, created_at, updated_at FROM releases"
	countQuery := "SELECT COUNT(*) FROM releases"

	where := ""
	if projectID != "" {
		where += "project_id = ?"
		args = append(args, projectID)
	}
	if status != "" {
		if where != "" {
			where += " AND "
		}
		where += "status = ?"
		args = append(args, status)
	}
	if where != "" {
		query += " WHERE " + where
		countQuery += " WHERE " + where
	}

	var total int
	if err := rs.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
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

	rows, err := rs.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*center.Release
	for rows.Next() {
		var version, summary, firstPub, lastPub sql.NullString
		var publishCount sql.NullInt32
		r := &center.Release{}
		if err := rows.Scan(&r.ID, &r.ProjectId, &r.Name, &version, &r.Status, &summary, &publishCount, &firstPub, &lastPub, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, 0, err
		}
		r.Version = version.String
		r.Summary = summary.String
		r.PublishCount = publishCount.Int32
		r.FirstPublishedAt = firstPub.String
		r.LastPublishedAt = lastPub.String
		list = append(list, r)
	}
	return list, total, nil
}

func (rs *ReleaseStore) Update(id string, fields map[string]interface{}) error {
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
	_, err := rs.db.Exec("UPDATE releases SET "+setClauses+" WHERE id = ?", args...)
	return err
}

func (rs *ReleaseStore) Delete(id string) error {
	_, err := rs.db.Exec("DELETE FROM releases WHERE id = ?", id)
	return err
}

func (rs *ReleaseStore) IncrementPublish(id string) error {
	now := time.Now().Format(time.DateTime)
	_, err := rs.db.Exec(`UPDATE releases SET publish_count = publish_count + 1, status = 'published', last_published_at = ?,
		first_published_at = COALESCE(first_published_at, ?), updated_at = ? WHERE id = ?`, now, now, now, id)
	return err
}

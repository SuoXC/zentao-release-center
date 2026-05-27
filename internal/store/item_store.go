package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
)

type ItemStore struct {
	db *sql.DB
}

func NewItemStore(s *Store) *ItemStore {
	return &ItemStore{db: s.DB()}
}

func (is *ItemStore) DB() *sql.DB {
	return is.db
}

func strPtr(s string) *string { return &s }
func int32Ptr(i int32) *int32  { return &i }

func scanItem(row *sql.Row) (*center.ReleaseItem, error) {
	item := &center.ReleaseItem{}
	var zentaoID sql.NullInt32
	var zentaoType, title, severity, priority, status, assignedTo, resolvedBy, zentaoURL, steps, noteTitle, noteContent sql.NullString
	err := row.Scan(&item.ID, &item.ReleaseId, &item.ItemType, &item.SortOrder,
		&zentaoID, &zentaoType, &title, &severity, &priority, &status, &assignedTo, &resolvedBy, &zentaoURL, &steps,
		&noteTitle, &noteContent, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if zentaoID.Valid {
		item.ZentaoId = int32Ptr(zentaoID.Int32)
	}
	if zentaoType.Valid {
		item.ZentaoType = strPtr(zentaoType.String)
	}
	if title.Valid {
		item.Title = strPtr(title.String)
	}
	if severity.Valid {
		item.Severity = strPtr(severity.String)
	}
	if priority.Valid {
		item.Priority = strPtr(priority.String)
	}
	if status.Valid {
		item.Status = strPtr(status.String)
	}
	if assignedTo.Valid {
		item.AssignedTo = strPtr(assignedTo.String)
	}
	if resolvedBy.Valid {
		item.ResolvedBy = strPtr(resolvedBy.String)
	}
	if zentaoURL.Valid {
		item.ZentaoUrl = strPtr(zentaoURL.String)
	}
	if steps.Valid {
		item.Steps = strPtr(steps.String)
	}
	if noteTitle.Valid {
		item.NoteTitle = strPtr(noteTitle.String)
	}
	if noteContent.Valid {
		item.NoteContent = strPtr(noteContent.String)
	}
	return item, nil
}

func (is *ItemStore) Add(releaseID, itemType string, zentaoID int, zentaoType, title, severity, priority, status, assignedTo, resolvedBy, zentaoURL, steps, noteTitle, noteContent string) (*center.ReleaseItem, error) {
	now := time.Now().Format(time.DateTime)
	id := uuid.New().String()

	var maxOrder sql.NullInt32
	is.db.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM release_items WHERE release_id = ?", releaseID).Scan(&maxOrder)
	sortOrder := int(maxOrder.Int32) + 1

	_, err := is.db.Exec(`INSERT INTO release_items (id, release_id, item_type, sort_order, zentao_id, zentao_type, title, severity, priority, status, assigned_to, resolved_by, zentao_url, steps, note_title, note_content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, releaseID, itemType, sortOrder, zentaoID, zentaoType, title, severity, priority, status, assignedTo, resolvedBy, zentaoURL, steps, noteTitle, noteContent, now, now)
	if err != nil {
		return nil, err
	}
	return is.GetByID(id)
}

func (is *ItemStore) ExistsByZentaoID(releaseID string, zentaoID int) (bool, error) {
	var count int
	err := is.db.QueryRow("SELECT COUNT(*) FROM release_items WHERE release_id = ? AND zentao_id = ?", releaseID, zentaoID).Scan(&count)
	return count > 0, err
}

func (is *ItemStore) AddBatch(tx *sql.Tx, releaseID string, items []struct {
	ItemType, ZentaoType, Title, Severity, Priority, Status, AssignedTo, ResolvedBy, ZentaoURL, Steps, NoteTitle, NoteContent string
	ZentaoID int
}) ([]*center.ReleaseItem, error) {
	var maxOrder sql.NullInt32
	tx.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM release_items WHERE release_id = ?", releaseID).Scan(&maxOrder)
	sortOrder := int(maxOrder.Int32)

	var result []*center.ReleaseItem
	now := time.Now().Format(time.DateTime)
	for _, item := range items {
		sortOrder++
		id := uuid.New().String()
		_, err := tx.Exec(`INSERT INTO release_items (id, release_id, item_type, sort_order, zentao_id, zentao_type, title, severity, priority, status, assigned_to, resolved_by, zentao_url, steps, note_title, note_content, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			id, releaseID, item.ItemType, sortOrder, item.ZentaoID, item.ZentaoType, item.Title, item.Severity, item.Priority, item.Status, item.AssignedTo, item.ResolvedBy, item.ZentaoURL, item.Steps, item.NoteTitle, item.NoteContent, now, now)
		if err != nil {
			return nil, err
		}
		created, err := is.GetByID(id)
		if err != nil {
			return nil, err
		}
		result = append(result, created)
	}
	return result, nil
}

func (is *ItemStore) GetByID(id string) (*center.ReleaseItem, error) {
	return scanItem(is.db.QueryRow(`SELECT id, release_id, item_type, sort_order, zentao_id, zentao_type, title, severity, priority, status, assigned_to, resolved_by, zentao_url, steps, note_title, note_content, created_at, updated_at FROM release_items WHERE id = ?`, id))
}

func (is *ItemStore) ListByRelease(releaseID string) ([]*center.ReleaseItem, error) {
	rows, err := is.db.Query(`SELECT id, release_id, item_type, sort_order, zentao_id, zentao_type, title, severity, priority, status, assigned_to, resolved_by, zentao_url, steps, note_title, note_content, created_at, updated_at FROM release_items WHERE release_id = ? ORDER BY sort_order ASC, created_at ASC`, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*center.ReleaseItem
	for rows.Next() {
		item := &center.ReleaseItem{}
		var zentaoID sql.NullInt32
		var zentaoType, title, severity, priority, status, assignedTo, resolvedBy, zentaoURL, steps, noteTitle, noteContent sql.NullString
		if err := rows.Scan(&item.ID, &item.ReleaseId, &item.ItemType, &item.SortOrder,
			&zentaoID, &zentaoType, &title, &severity, &priority, &status, &assignedTo, &resolvedBy, &zentaoURL, &steps,
			&noteTitle, &noteContent, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if zentaoID.Valid {
			item.ZentaoId = int32Ptr(zentaoID.Int32)
		}
		if zentaoType.Valid {
			item.ZentaoType = strPtr(zentaoType.String)
		}
		if title.Valid {
			item.Title = strPtr(title.String)
		}
		if severity.Valid {
			item.Severity = strPtr(severity.String)
		}
		if priority.Valid {
			item.Priority = strPtr(priority.String)
		}
		if status.Valid {
			item.Status = strPtr(status.String)
		}
		if assignedTo.Valid {
			item.AssignedTo = strPtr(assignedTo.String)
		}
		if resolvedBy.Valid {
			item.ResolvedBy = strPtr(resolvedBy.String)
		}
		if zentaoURL.Valid {
			item.ZentaoUrl = strPtr(zentaoURL.String)
		}
		if steps.Valid {
			item.Steps = strPtr(steps.String)
		}
		if noteTitle.Valid {
			item.NoteTitle = strPtr(noteTitle.String)
		}
		if noteContent.Valid {
			item.NoteContent = strPtr(noteContent.String)
		}
		list = append(list, item)
	}
	return list, nil
}

var itemAllowedFields = map[string]bool{
	"title": true, "severity": true, "priority": true, "status": true,
	"assigned_to": true, "resolved_by": true, "steps": true,
	"note_title": true, "note_content": true, "sort_order": true,
}

func (is *ItemStore) Update(id string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	fields["updated_at"] = time.Now().Format(time.DateTime)
	setClauses := ""
	args := []interface{}{}
	for k, v := range fields {
		if !itemAllowedFields[k] && k != "updated_at" {
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
	_, err := is.db.Exec("UPDATE release_items SET "+setClauses+" WHERE id = ?", args...)
	return err
}

func (is *ItemStore) Delete(id string) error {
	_, err := is.db.Exec("DELETE FROM release_items WHERE id = ?", id)
	return err
}

func (is *ItemStore) Reorder(items []struct {
	ID        string
	SortOrder int
}) error {
	tx, err := is.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range items {
		if _, err := tx.Exec("UPDATE release_items SET sort_order = ? WHERE id = ?", item.SortOrder, item.ID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (is *ItemStore) CountByType(releaseID string) (total, bugs, tasks, notes int, err error) {
	rows, err := is.db.Query("SELECT item_type, COUNT(*) FROM release_items WHERE release_id = ? GROUP BY item_type", releaseID)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var itemType string
		var count int
		if err := rows.Scan(&itemType, &count); err != nil {
			return 0, 0, 0, 0, err
		}
		total += count
		switch itemType {
		case "bug":
			bugs = count
		case "task":
			tasks = count
		case "note":
			notes = count
		}
	}
	return
}

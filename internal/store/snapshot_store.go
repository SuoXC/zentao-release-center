package store

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	center "github.com/yi-nology/zentao-release-center/biz/model/release/center"
)

type SnapshotStore struct {
	db *sql.DB
}

func NewSnapshotStore(s *Store) *SnapshotStore {
	return &SnapshotStore{db: s.DB()}
}

func (ss *SnapshotStore) Create(releaseID, version, content string, itemCount, bugCount, taskCount, noteCount int) (*center.ReleaseSnapshot, error) {
	id := uuid.New().String()
	now := time.Now().Format(time.DateTime)
	_, err := ss.db.Exec(`INSERT INTO release_snapshots (id, release_id, version, content, item_count, bug_count, task_count, note_count, published_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		id, releaseID, version, content, itemCount, bugCount, taskCount, noteCount, now)
	if err != nil {
		return nil, err
	}
	return ss.GetByID(id)
}

func (ss *SnapshotStore) GetByID(id string) (*center.ReleaseSnapshot, error) {
	snap := &center.ReleaseSnapshot{}
	err := ss.db.QueryRow(`SELECT id, release_id, version, content, item_count, bug_count, task_count, note_count, published_at FROM release_snapshots WHERE id = ?`, id).
		Scan(&snap.ID, &snap.ReleaseId, &snap.Version, &snap.Content, &snap.ItemCount, &snap.BugCount, &snap.TaskCount, &snap.NoteCount, &snap.PublishedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return snap, nil
}

func (ss *SnapshotStore) ListByRelease(releaseID string) ([]*center.ReleaseSnapshot, error) {
	rows, err := ss.db.Query(`SELECT id, release_id, version, content, item_count, bug_count, task_count, note_count, published_at FROM release_snapshots WHERE release_id = ? ORDER BY published_at DESC`, releaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*center.ReleaseSnapshot
	for rows.Next() {
		snap := &center.ReleaseSnapshot{}
		if err := rows.Scan(&snap.ID, &snap.ReleaseId, &snap.Version, &snap.Content, &snap.ItemCount, &snap.BugCount, &snap.TaskCount, &snap.NoteCount, &snap.PublishedAt); err != nil {
			return nil, err
		}
		list = append(list, snap)
	}
	return list, nil
}

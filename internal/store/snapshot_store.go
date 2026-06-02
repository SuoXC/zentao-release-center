package store

import (
	"time"

	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type SnapshotStore struct {
	db *gorm.DB
}

func NewSnapshotStore(db *gorm.DB) *SnapshotStore {
	return &SnapshotStore{db: db}
}

func (ss *SnapshotStore) Create(releaseKeyword, version, content string, itemCount, bugCount, taskCount, noteCount int) (*model.ReleaseSnapshot, error) {
	snap := &model.ReleaseSnapshot{
		Keyword:        uuid.New().String(),
		ReleaseKeyword: releaseKeyword,
		Version:        version,
		Content:        content,
		ItemCount:      itemCount,
		BugCount:       bugCount,
		TaskCount:      taskCount,
		NoteCount:      noteCount,
		PublishedAt:    time.Now(),
	}
	if err := ss.db.Create(snap).Error; err != nil {
		return nil, err
	}
	return snap, nil
}

func (ss *SnapshotStore) GetByID(keyword string) (*model.ReleaseSnapshot, error) {
	var snap model.ReleaseSnapshot
	if err := ss.db.Where("keyword = ?", keyword).First(&snap).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &snap, nil
}

func (ss *SnapshotStore) ListByRelease(releaseKeyword string) ([]*model.ReleaseSnapshot, error) {
	var snaps []*model.ReleaseSnapshot
	if err := ss.db.Where("release_keyword = ?", releaseKeyword).Order("published_at DESC").Find(&snaps).Error; err != nil {
		return nil, err
	}
	return snaps, nil
}

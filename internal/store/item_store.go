package store

import (
	"github.com/google/uuid"
	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/gorm"
)

type ItemStore struct {
	db *gorm.DB
}

func NewItemStore(db *gorm.DB) *ItemStore {
	return &ItemStore{db: db}
}

func (is *ItemStore) Add(releaseKeyword, itemType string, zentaoID int, zentaoType, title, severity, priority, status, assignedTo, resolvedBy, zentaoURL, steps, noteTitle, noteContent string) (*model.ReleaseItem, error) {
	var maxOrder int
	is.db.Model(&model.ReleaseItem{}).Where("release_keyword = ?", releaseKeyword).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	item := &model.ReleaseItem{
		Keyword:        uuid.New().String(),
		ReleaseKeyword: releaseKeyword,
		ItemType:       itemType,
		SortOrder:      maxOrder + 1,
		ZentaoID:       zentaoID,
		ZentaoType:     zentaoType,
		Title:          title,
		Severity:       severity,
		Priority:       priority,
		Status:         status,
		AssignedTo:     assignedTo,
		ResolvedBy:     resolvedBy,
		ZentaoURL:      zentaoURL,
		Steps:          steps,
		NoteTitle:      noteTitle,
		NoteContent:    noteContent,
	}
	if err := is.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (is *ItemStore) ExistsByZentaoID(releaseKeyword string, zentaoID int) (bool, error) {
	var count int64
	err := is.db.Model(&model.ReleaseItem{}).Where("release_keyword = ? AND zentao_id = ?", releaseKeyword, zentaoID).Count(&count).Error
	return count > 0, err
}

func (is *ItemStore) AddBatch(releaseKeyword string, items []struct {
	ItemType, ZentaoType, Title, Severity, Priority, Status, AssignedTo, ResolvedBy, ZentaoURL, Steps, NoteTitle, NoteContent string
	ZentaoID int
}) ([]*model.ReleaseItem, error) {
	var maxOrder int
	is.db.Model(&model.ReleaseItem{}).Where("release_keyword = ?", releaseKeyword).Select("COALESCE(MAX(sort_order), 0)").Scan(&maxOrder)

	var result []*model.ReleaseItem
	for _, item := range items {
		maxOrder++
		m := &model.ReleaseItem{
			Keyword:        uuid.New().String(),
			ReleaseKeyword: releaseKeyword,
			ItemType:       item.ItemType,
			SortOrder:      maxOrder,
			ZentaoID:       item.ZentaoID,
			ZentaoType:     item.ZentaoType,
			Title:          item.Title,
			Severity:       item.Severity,
			Priority:       item.Priority,
			Status:         item.Status,
			AssignedTo:     item.AssignedTo,
			ResolvedBy:     item.ResolvedBy,
			ZentaoURL:      item.ZentaoURL,
			Steps:          item.Steps,
			NoteTitle:      item.NoteTitle,
			NoteContent:    item.NoteContent,
		}
		if err := is.db.Create(m).Error; err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, nil
}

func (is *ItemStore) GetByID(keyword string) (*model.ReleaseItem, error) {
	var item model.ReleaseItem
	if err := is.db.Where("keyword = ?", keyword).First(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (is *ItemStore) ListByRelease(releaseKeyword string) ([]*model.ReleaseItem, error) {
	var items []*model.ReleaseItem
	if err := is.db.Where("release_keyword = ?", releaseKeyword).Order("sort_order ASC, created_at ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

var itemAllowedFields = map[string]bool{
	"title": true, "severity": true, "priority": true, "status": true,
	"assigned_to": true, "resolved_by": true, "steps": true,
	"note_title": true, "note_content": true, "sort_order": true,
}

func (is *ItemStore) Update(keyword string, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil
	}
	allowed := make(map[string]interface{})
	for k, v := range fields {
		if itemAllowedFields[k] {
			allowed[k] = v
		}
	}
	if len(allowed) == 0 {
		return nil
	}
	return is.db.Model(&model.ReleaseItem{}).Where("keyword = ?", keyword).Updates(allowed).Error
}

func (is *ItemStore) Delete(keyword string) error {
	return is.db.Where("keyword = ?", keyword).Delete(&model.ReleaseItem{}).Error
}

func (is *ItemStore) Reorder(items []struct {
	Keyword   string
	SortOrder int
}) error {
	return is.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.ReleaseItem{}).Where("keyword = ?", item.Keyword).Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (is *ItemStore) CountByType(releaseKeyword string) (total, bugs, tasks, notes int, err error) {
	type result struct {
		ItemType string
		Count    int
	}
	var results []result
	if err = is.db.Model(&model.ReleaseItem{}).Where("release_keyword = ?", releaseKeyword).
		Select("item_type, COUNT(*) as count").Group("item_type").Scan(&results).Error; err != nil {
		return
	}
	for _, r := range results {
		total += r.Count
		switch r.ItemType {
		case "bug":
			bugs = r.Count
		case "task":
			tasks = r.Count
		case "note":
			notes = r.Count
		}
	}
	return
}

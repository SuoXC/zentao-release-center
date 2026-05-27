package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *Store) DB() *sql.DB {
	return s.db
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		zentao_product_id INTEGER DEFAULT 0,
		zentao_project_id INTEGER DEFAULT 0,
		zentao_product_name TEXT DEFAULT '',
		zentao_project_name TEXT DEFAULT '',
		zentao_server TEXT DEFAULT '',
		status TEXT DEFAULT 'active',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS releases (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		name TEXT NOT NULL,
		version TEXT DEFAULT '',
		status TEXT DEFAULT 'draft',
		summary TEXT DEFAULT '',
		publish_count INTEGER DEFAULT 0,
		first_published_at DATETIME,
		last_published_at DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_releases_project ON releases(project_id);

	CREATE TABLE IF NOT EXISTS release_items (
		id TEXT PRIMARY KEY,
		release_id TEXT NOT NULL,
		item_type TEXT NOT NULL,
		sort_order INTEGER DEFAULT 0,
		zentao_id INTEGER DEFAULT 0,
		zentao_type TEXT DEFAULT '',
		title TEXT DEFAULT '',
		severity TEXT DEFAULT '',
		priority TEXT DEFAULT '',
		status TEXT DEFAULT '',
		assigned_to TEXT DEFAULT '',
		resolved_by TEXT DEFAULT '',
		zentao_url TEXT DEFAULT '',
		steps TEXT DEFAULT '',
		note_title TEXT DEFAULT '',
		note_content TEXT DEFAULT '',
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_items_release ON release_items(release_id);

	CREATE TABLE IF NOT EXISTS release_snapshots (
		id TEXT PRIMARY KEY,
		release_id TEXT NOT NULL,
		version TEXT DEFAULT '',
		content TEXT NOT NULL,
		item_count INTEGER DEFAULT 0,
		bug_count INTEGER DEFAULT 0,
		task_count INTEGER DEFAULT 0,
		note_count INTEGER DEFAULT 0,
		published_at DATETIME NOT NULL,
		FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_snapshots_release ON release_snapshots(release_id);

	CREATE TABLE IF NOT EXISTS release_deployments (
		id TEXT PRIMARY KEY,
		release_id TEXT NOT NULL,
		module_name TEXT NOT NULL,
		address TEXT NOT NULL,
		description TEXT DEFAULT '',
		sort_order INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (release_id) REFERENCES releases(id) ON DELETE CASCADE
	);
	CREATE INDEX IF NOT EXISTS idx_deployments_release ON release_deployments(release_id);
	`
	_, err := s.db.Exec(schema)
	return err
}

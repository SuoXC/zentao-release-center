package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yi-nology/zentao-release-center/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Store struct {
	db *gorm.DB
}

func NewStore(dbPath string) (*Store, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying db: %w", err)
	}
	sqlDB.SetMaxOpenConns(5)
	sqlDB.SetMaxIdleConns(2)

	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	if err := db.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		return nil, fmt.Errorf("enable WAL: %w", err)
	}

	if err := db.AutoMigrate(
		&model.Project{},
		&model.Release{},
		&model.ReleaseItem{},
		&model.ReleaseSnapshot{},
		&model.ReleaseDeployment{},
		&model.ProjectRepo{},
		&model.ReleaseBranch{},
		&model.DockerImage{},
		&model.DockerImagePool{},
		&model.ReleaseFeature{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) DB() *gorm.DB {
	return s.db
}

func (s *Store) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

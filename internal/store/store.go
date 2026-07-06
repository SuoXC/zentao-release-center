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
		&model.ProjectRepo{},
		&model.ReleaseBranch{},
		&model.DockerImage{},
		&model.DockerImagePool{},
		&model.ReleaseFeature{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	if db.Migrator().HasTable("release_deployments") {
		if err := db.Migrator().DropTable("release_deployments"); err != nil {
			return nil, fmt.Errorf("drop legacy release_deployments: %w", err)
		}
	}

	if db.Migrator().HasTable(&model.DockerImage{}) {
		for _, col := range []string{"image_name", "image_tag", "registry", "branch"} {
			if db.Migrator().HasColumn(&model.DockerImage{}, col) {
				if err := db.Migrator().DropColumn(&model.DockerImage{}, col); err != nil {
					return nil, fmt.Errorf("drop legacy column %s on docker_images: %w", col, err)
				}
			}
		}
	}

	if db.Migrator().HasTable(&model.DockerImagePool{}) {
		for _, col := range []string{"image_name", "image_tag", "registry", "branch"} {
			if db.Migrator().HasColumn(&model.DockerImagePool{}, col) {
				if err := db.Migrator().DropColumn(&model.DockerImagePool{}, col); err != nil {
					return nil, fmt.Errorf("drop legacy column %s on docker_image_pool: %w", col, err)
				}
			}
		}
	}

	db.Exec(`DROP INDEX IF EXISTS idx_release_branch_unique`)
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_release_branch_unique ON release_branches (release_keyword) WHERE branch_type = 'release'`).Error; err != nil {
		return nil, fmt.Errorf("create release_branch unique index: %w", err)
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

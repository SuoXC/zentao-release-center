package model

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	Name                string
	Description         string `gorm:"default:''"`
	ZentaoProductID     int    `gorm:"default:0"`
	ZentaoProjectID     int    `gorm:"default:0"`
	ZentaoProductName   string `gorm:"default:''"`
	ZentaoProjectName   string `gorm:"default:''"`
	ZentaoServer        string `gorm:"default:''"`
	Status              string `gorm:"default:'active'"`
}

type Release struct {
	gorm.Model          `json:"-"`
	Keyword             string     `json:"-" gorm:"uniqueIndex;size:36"`
	ProjectKeyword      string     `gorm:"index;size:36"`
	Name                string
	Version             string     `gorm:"default:''"`
	Status              string     `gorm:"default:'draft'"`
	Summary             string     `gorm:"default:''"`
	ParentBranch        string     `gorm:"default:''"`
	PublishCount        int        `gorm:"default:0"`
	FirstPublishedAt    *time.Time
	LastPublishedAt     *time.Time
}

type ReleaseItem struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	ReleaseKeyword      string `gorm:"index;size:36"`
	ItemType            string
	SortOrder           int    `gorm:"default:0"`
	ZentaoID            int    `gorm:"default:0"`
	ZentaoType          string `gorm:"default:''"`
	Title               string `gorm:"default:''"`
	Severity            string `gorm:"default:''"`
	Priority            string `gorm:"default:''"`
	Status              string `gorm:"default:''"`
	AssignedTo          string `gorm:"default:''"`
	ResolvedBy          string `gorm:"default:''"`
	ZentaoURL           string `gorm:"default:''"`
	Steps               string `gorm:"default:''"`
	NoteTitle           string `gorm:"default:''"`
	NoteContent         string `gorm:"default:''"`
}

type ReleaseSnapshot struct {
	gorm.Model          `json:"-"`
	Keyword             string    `json:"-" gorm:"uniqueIndex;size:36"`
	ReleaseKeyword      string    `gorm:"index;size:36"`
	Version             string    `gorm:"default:''"`
	Content             string
	ItemCount           int       `gorm:"default:0"`
	BugCount            int       `gorm:"default:0"`
	TaskCount           int       `gorm:"default:0"`
	NoteCount           int       `gorm:"default:0"`
	PublishedAt         time.Time
}

type ProjectRepo struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	ProjectKeyword      string `gorm:"index;size:36"`
	GitlabProjectID     int
	RepoURL             string
	RepoName            string
	DefaultBranch       string `gorm:"default:'main'"`
}

type ReleaseBranch struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	ReleaseKeyword      string `gorm:"index;size:36"`
	RepoKeyword         string `gorm:"index;size:36"`
	BranchName          string
	BranchType          string `gorm:"default:'release'"`
	ParentBranch        string `gorm:"default:''"`
	GitlabBranchURL     string `gorm:"default:''"`
	Description         string `gorm:"default:''"`
}

type DockerImage struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	ReleaseKeyword      string `gorm:"index;size:36"`
	RepoKeyword         string `gorm:"index;size:36"`
	ImageURL            string `gorm:"default:''"`
	ImageDigest         string `gorm:"default:''"`
	CIPipelineID        int    `gorm:"default:0"`
	CIPipelineURL       string `gorm:"default:''"`
	CommitSHA           string `gorm:"default:''"`
	CommitMessage       string `gorm:"default:''"`
	Source              string `gorm:"default:'manual'"`
	Tested              bool   `gorm:"default:false"`
}

type DockerImagePool struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	GitlabProjectID     int    `gorm:"index"`
	ImageURL            string `gorm:"default:''"`
	ImageDigest         string `gorm:"default:''"`
	CIPipelineID        int    `gorm:"default:0"`
	CIPipelineURL       string `gorm:"default:''"`
	CommitSHA           string `gorm:"default:''"`
	CommitMessage       string `gorm:"default:''"`
}

type ReleaseFeature struct {
	gorm.Model          `json:"-"`
	Keyword             string `json:"-" gorm:"uniqueIndex;size:36"`
	ReleaseKeyword      string `gorm:"index;size:36"`
	Title               string
	Content             string `gorm:"type:text"`
	SortOrder           int    `gorm:"default:0"`
}

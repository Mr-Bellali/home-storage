package models

import (
	"gorm.io/gorm"
	"time"
)

// User represents a user in the system
type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Workspaces []Workspace    `json:"workspaces,omitempty" gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`

	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"password"`
}

type Workspace struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Hashkey     *string `json:"hashkey" gorm:"unique;"`
	UserId      uint    `json:"owner_id" gorm:"not null"`
}

type File struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Filename    string         `gorm:"not null" json:"filename"`
	Filepath    string         `gorm:"not null" json:"filepath"` 
	Size        int64          `json:"size"`                     
	MIMEType    string         `json:"mime_type"`                
	UserID      uint           `gorm:"not null" json:"user_id"`  
	WorkspaceID uint           `gorm:"not null" json:"workspace_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Associations
	User      User      `gorm:"foreignKey:UserID"`
	Workspace Workspace `gorm:"foreignKey:WorkspaceID"`
}

// GetAllModels returns a slice of all model structs for migration
func GetAllModels() []interface{} {
	return []interface{}{
		&User{},
		&Workspace{},
		&File{},
	}
}

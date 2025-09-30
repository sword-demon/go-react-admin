// Package model provides GORM data models
package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"` // Soft delete
}

// Status constants
const (
	StatusEnabled  uint8 = 1 // Active/Enabled
	StatusDisabled uint8 = 0 // Inactive/Disabled
)

// Gender constants
const (
	GenderUnknown uint8 = 0
	GenderMale    uint8 = 1
	GenderFemale  uint8 = 2
)

// DataScope constants (data permission scope)
const (
	DataScopeAll          uint8 = 1 // All data
	DataScopeDeptAndChild uint8 = 2 // Current dept + child depts
	DataScopeDeptOnly     uint8 = 3 // Current dept only
	DataScopeSelfOnly     uint8 = 4 // Only own data
)

// MenuType constants
const (
	MenuTypeDirectory uint8 = 1 // Directory (folder)
	MenuTypeMenu      uint8 = 2 // Menu (page)
	MenuTypeButton    uint8 = 3 // Button (permission)
)

// PermissionType constants (pattern type)
const (
	PermissionTypeGlobal string = "global" // *:*
	PermissionTypeModule string = "module" // user:*
	PermissionTypeAction string = "action" // user:read
	PermissionTypePath   string = "path"   // /api/users/*
)

// TableName returns the table name for a model (helper function)
func TableName(name string) string {
	return "sys_" + name
}

package model

// RolePermission represents sys_role_permission table (permission pattern matching)
type RolePermission struct {
	BaseModel
	RoleID            uint64 `gorm:"column:role_id;not null;index:idx_role_id" json:"role_id"`
	PermissionPattern string `gorm:"column:permission_pattern;type:varchar(100);not null" json:"permission_pattern"` // e.g., user:*, /api/users/*
	PermissionType    string `gorm:"column:permission_type;type:varchar(20);not null;index:idx_type" json:"permission_type"`
	Description       string `gorm:"column:description;type:varchar(200)" json:"description"`
	Status            uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Active,0=Disabled'" json:"status"`

	// Relations
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName returns the table name
func (RolePermission) TableName() string {
	return "sys_role_permission"
}

// IsActive checks if permission is active
func (p *RolePermission) IsActive() bool {
	return p.Status == StatusEnabled
}

// IsGlobal checks if permission is global (*:*)
func (p *RolePermission) IsGlobal() bool {
	return p.PermissionType == PermissionTypeGlobal
}

// IsModule checks if permission is module-level (user:*)
func (p *RolePermission) IsModule() bool {
	return p.PermissionType == PermissionTypeModule
}

// IsAction checks if permission is action-level (user:read)
func (p *RolePermission) IsAction() bool {
	return p.PermissionType == PermissionTypeAction
}

// IsPath checks if permission is path-level (/api/users/*)
func (p *RolePermission) IsPath() bool {
	return p.PermissionType == PermissionTypePath
}

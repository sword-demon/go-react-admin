package model

// UserRole represents sys_user_role table (many-to-many: user-role)
type UserRole struct {
	UserID uint64 `gorm:"column:user_id;primaryKey;index:idx_user_id" json:"user_id"`
	RoleID uint64 `gorm:"column:role_id;primaryKey;index:idx_role_id" json:"role_id"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName returns the table name
func (UserRole) TableName() string {
	return "sys_user_role"
}

// RoleMenu represents sys_role_menu table (many-to-many: role-menu)
type RoleMenu struct {
	RoleID uint64 `gorm:"column:role_id;primaryKey;index:idx_role_id" json:"role_id"`
	MenuID uint64 `gorm:"column:menu_id;primaryKey;index:idx_menu_id" json:"menu_id"`

	// Relations
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Menu *Menu `gorm:"foreignKey:MenuID" json:"menu,omitempty"`
}

// TableName returns the table name
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

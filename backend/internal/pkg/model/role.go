package model

// Role represents sys_role table
type Role struct {
	BaseModel
	RoleName  string `gorm:"column:role_name;type:varchar(64);not null" json:"role_name"`
	RoleKey   string `gorm:"column:role_key;type:varchar(64);not null;uniqueIndex:idx_role_key" json:"role_key"`
	RoleSort  int    `gorm:"column:role_sort;not null;default:0" json:"role_sort"`
	DataScope uint8  `gorm:"column:data_scope;type:tinyint;not null;default:4;comment:'1=All,2=DeptAndChild,3=DeptOnly,4=SelfOnly'" json:"data_scope"`
	Status    uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Enabled,0=Disabled'" json:"status"`
	Remark    string `gorm:"column:remark;type:varchar(500)" json:"remark"`

	// Relations
	Users       []*User           `gorm:"many2many:sys_user_role;" json:"users,omitempty"`
	Menus       []*Menu           `gorm:"many2many:sys_role_menu;" json:"menus,omitempty"`
	Permissions []*RolePermission `gorm:"foreignKey:RoleID" json:"permissions,omitempty"`
}

// TableName returns the table name
func (Role) TableName() string {
	return "sys_role"
}

// IsEnabled checks if role is enabled
func (r *Role) IsEnabled() bool {
	return r.Status == StatusEnabled
}

// HasFullDataAccess checks if role has access to all data
func (r *Role) HasFullDataAccess() bool {
	return r.DataScope == DataScopeAll
}

// CanAccessChildDepts checks if role can access child departments
func (r *Role) CanAccessChildDepts() bool {
	return r.DataScope == DataScopeAll || r.DataScope == DataScopeDeptAndChild
}

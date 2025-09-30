package model

// Menu represents sys_menu table (tree structure)
type Menu struct {
	BaseModel
	ParentID  uint64 `gorm:"column:parent_id;not null;default:0;index:idx_parent_id" json:"parent_id"`
	MenuName  string `gorm:"column:menu_name;type:varchar(64);not null" json:"menu_name"`
	MenuType  uint8  `gorm:"column:menu_type;type:tinyint;not null;default:2;comment:'1=Directory,2=Menu,3=Button'" json:"menu_type"`
	Path      string `gorm:"column:path;type:varchar(200)" json:"path"`           // Route path (e.g., /users)
	Component string `gorm:"column:component;type:varchar(200)" json:"component"` // Component path (e.g., @/views/users/index)
	Icon      string `gorm:"column:icon;type:varchar(100)" json:"icon"`
	OrderNum  int    `gorm:"column:order_num;not null;default:0" json:"order_num"`
	PermKey   string `gorm:"column:perm_key;type:varchar(100);index:idx_perm_key" json:"perm_key"` // e.g., user:list, user:create
	Visible   uint8  `gorm:"column:visible;type:tinyint;not null;default:1;comment:'1=Visible,0=Hidden'" json:"visible"`
	Status    uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Enabled,0=Disabled'" json:"status"`
	Remark    string `gorm:"column:remark;type:varchar(500)" json:"remark"`

	// Relations
	Parent   *Menu   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []*Menu `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Roles    []*Role `gorm:"many2many:sys_role_menu;" json:"roles,omitempty"`
}

// TableName returns the table name
func (Menu) TableName() string {
	return "sys_menu"
}

// IsEnabled checks if menu is enabled
func (m *Menu) IsEnabled() bool {
	return m.Status == StatusEnabled
}

// IsVisible checks if menu is visible
func (m *Menu) IsVisible() bool {
	return m.Visible == StatusEnabled
}

// IsDirectory checks if menu is a directory
func (m *Menu) IsDirectory() bool {
	return m.MenuType == MenuTypeDirectory
}

// IsMenu checks if menu is a menu (page)
func (m *Menu) IsMenu() bool {
	return m.MenuType == MenuTypeMenu
}

// IsButton checks if menu is a button (permission)
func (m *Menu) IsButton() bool {
	return m.MenuType == MenuTypeButton
}

// IsRoot checks if menu is root (no parent)
func (m *Menu) IsRoot() bool {
	return m.ParentID == 0
}

// HasChildren checks if menu has children
func (m *Menu) HasChildren() bool {
	return len(m.Children) > 0
}

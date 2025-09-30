package model

// Dept represents sys_dept table (tree structure)
type Dept struct {
	BaseModel
	ParentID  uint64 `gorm:"column:parent_id;not null;default:0;index:idx_parent_id" json:"parent_id"`
	Ancestors string `gorm:"column:ancestors;type:varchar(500);default:'';index:idx_ancestors" json:"ancestors"` // e.g., "0,1,2"
	DeptName  string `gorm:"column:dept_name;type:varchar(64);not null" json:"dept_name"`
	OrderNum  int    `gorm:"column:order_num;not null;default:0" json:"order_num"`
	Leader    string `gorm:"column:leader;type:varchar(64)" json:"leader"`
	Phone     string `gorm:"column:phone;type:varchar(20)" json:"phone"`
	Email     string `gorm:"column:email;type:varchar(128)" json:"email"`
	Status    uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Enabled,0=Disabled'" json:"status"`

	// Relations
	Parent   *Dept   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []*Dept `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []*User `gorm:"foreignKey:DeptID" json:"users,omitempty"`
}

// TableName returns the table name
func (Dept) TableName() string {
	return "sys_dept"
}

// IsEnabled checks if dept is enabled
func (d *Dept) IsEnabled() bool {
	return d.Status == StatusEnabled
}

// IsRoot checks if dept is root (no parent)
func (d *Dept) IsRoot() bool {
	return d.ParentID == 0
}

// HasChildren checks if dept has children
func (d *Dept) HasChildren() bool {
	return len(d.Children) > 0
}

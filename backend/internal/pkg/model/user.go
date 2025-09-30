package model

// User represents sys_user table
type User struct {
	BaseModel
	Username string `gorm:"column:username;type:varchar(64);not null;uniqueIndex:idx_username" json:"username"`
	Password string `gorm:"column:password;type:varchar(128);not null" json:"-"` // bcrypt hash
	NickName string `gorm:"column:nick_name;type:varchar(64)" json:"nick_name"`
	Email    string `gorm:"column:email;type:varchar(128);index:idx_email" json:"email"`
	Phone    string `gorm:"column:phone;type:varchar(20)" json:"phone"`
	Gender   uint8  `gorm:"column:gender;type:tinyint;default:0;comment:'0=Unknown,1=Male,2=Female'" json:"gender"`
	Avatar   string `gorm:"column:avatar;type:varchar(255)" json:"avatar"`
	DeptID   uint64 `gorm:"column:dept_id;index:idx_dept_id" json:"dept_id"`
	Status   uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Enabled,0=Disabled'" json:"status"`
	Remark   string `gorm:"column:remark;type:varchar(500)" json:"remark"`

	// Relations
	Dept  *Dept   `gorm:"foreignKey:DeptID" json:"dept,omitempty"`
	Roles []*Role `gorm:"many2many:sys_user_role;" json:"roles,omitempty"`
}

// TableName returns the table name
func (User) TableName() string {
	return "sys_user"
}

// IsEnabled checks if user is enabled
func (u *User) IsEnabled() bool {
	return u.Status == StatusEnabled
}

// IsMale checks if user is male
func (u *User) IsMale() bool {
	return u.Gender == GenderMale
}

// IsFemale checks if user is female
func (u *User) IsFemale() bool {
	return u.Gender == GenderFemale
}

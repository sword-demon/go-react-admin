package model

// LoginLog represents sys_login_log table
type LoginLog struct {
	BaseModel
	UserID    uint64 `gorm:"column:user_id;index:idx_user_id" json:"user_id"`
	Username  string `gorm:"column:username;type:varchar(64);not null" json:"username"`
	IPAddress string `gorm:"column:ip_address;type:varchar(128)" json:"ip_address"`
	Location  string `gorm:"column:location;type:varchar(255)" json:"location"` // IP location
	Browser   string `gorm:"column:browser;type:varchar(200)" json:"browser"`
	OS        string `gorm:"column:os;type:varchar(200)" json:"os"`
	Status    uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Success,0=Failed'" json:"status"`
	Message   string `gorm:"column:message;type:varchar(500)" json:"message"`
	LoginTime string `gorm:"column:login_time;type:datetime" json:"login_time"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name
func (LoginLog) TableName() string {
	return "sys_login_log"
}

// IsSuccess checks if login was successful
func (l *LoginLog) IsSuccess() bool {
	return l.Status == StatusEnabled
}

// OperationLog represents sys_operation_log table (Phase 2)
type OperationLog struct {
	BaseModel
	UserID        uint64 `gorm:"column:user_id;index:idx_user_id" json:"user_id"`
	Username      string `gorm:"column:username;type:varchar(64);not null" json:"username"`
	Module        string `gorm:"column:module;type:varchar(64)" json:"module"`          // Module name (e.g., User Management)
	Operation     string `gorm:"column:operation;type:varchar(100)" json:"operation"`   // Operation type (e.g., Create, Update, Delete)
	Method        string `gorm:"column:method;type:varchar(200)" json:"method"`         // Request method (e.g., POST /api/v1/users)
	RequestParams string `gorm:"column:request_params;type:text" json:"request_params"` // Request parameters (JSON)
	ResponseData  string `gorm:"column:response_data;type:text" json:"response_data"`   // Response data (JSON, first 2000 chars)
	IPAddress     string `gorm:"column:ip_address;type:varchar(128)" json:"ip_address"`
	Location      string `gorm:"column:location;type:varchar(255)" json:"location"`
	Status        uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Success,0=Failed'" json:"status"`
	ErrorMessage  string `gorm:"column:error_message;type:varchar(500)" json:"error_message"`
	CostTime      int64  `gorm:"column:cost_time;comment:'milliseconds'" json:"cost_time"` // Request duration (ms)
	OperationTime string `gorm:"column:operation_time;type:datetime" json:"operation_time"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name
func (OperationLog) TableName() string {
	return "sys_operation_log"
}

// IsSuccess checks if operation was successful
func (o *OperationLog) IsSuccess() bool {
	return o.Status == StatusEnabled
}

// AuditLog represents sys_audit_log table (Phase 2 - sensitive operations)
type AuditLog struct {
	BaseModel
	UserID         uint64 `gorm:"column:user_id;index:idx_user_id" json:"user_id"`
	Username       string `gorm:"column:username;type:varchar(64);not null" json:"username"`
	AuditType      string `gorm:"column:audit_type;type:varchar(64);index:idx_type" json:"audit_type"` // e.g., RolePermissionChange, UserStatusChange
	TargetType     string `gorm:"column:target_type;type:varchar(64)" json:"target_type"`              // e.g., User, Role
	TargetID       uint64 `gorm:"column:target_id;index:idx_target" json:"target_id"`                  // Target record ID
	TargetName     string `gorm:"column:target_name;type:varchar(128)" json:"target_name"`             // Target record name
	Action         string `gorm:"column:action;type:varchar(100)" json:"action"`                       // e.g., GrantPermission, DisableUser
	BeforeSnapshot string `gorm:"column:before_snapshot;type:text" json:"before_snapshot"`             // JSON snapshot before change
	AfterSnapshot  string `gorm:"column:after_snapshot;type:text" json:"after_snapshot"`               // JSON snapshot after change
	Changes        string `gorm:"column:changes;type:text" json:"changes"`                             // Detailed changes (JSON)
	IPAddress      string `gorm:"column:ip_address;type:varchar(128)" json:"ip_address"`
	Location       string `gorm:"column:location;type:varchar(255)" json:"location"`
	Reason         string `gorm:"column:reason;type:varchar(500)" json:"reason"` // Reason for change
	Status         uint8  `gorm:"column:status;type:tinyint;not null;default:1;comment:'1=Success,0=Failed'" json:"status"`
	ErrorMessage   string `gorm:"column:error_message;type:varchar(500)" json:"error_message"`
	AuditTime      string `gorm:"column:audit_time;type:datetime" json:"audit_time"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName returns the table name
func (AuditLog) TableName() string {
	return "sys_audit_log"
}

// IsSuccess checks if audit action was successful
func (a *AuditLog) IsSuccess() bool {
	return a.Status == StatusEnabled
}

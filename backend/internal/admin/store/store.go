// Package store defines data access layer interfaces
package store

import (
	"context"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
)

// IStore represents the top-level data access interface
type IStore interface {
	Users() IUserStore
	Roles() IRoleStore
	Depts() IDeptStore
	Menus() IMenuStore
	Permissions() IPermissionStore
	Close() error
}

// IUserStore defines user data access operations
type IUserStore interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, opts *ListOptions) ([]*model.User, int64, error)
	GetUserRoles(ctx context.Context, userID uint64) ([]*model.Role, error)
	AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error
}

// IRoleStore defines role data access operations
type IRoleStore interface {
	Create(ctx context.Context, role *model.Role) error
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Role, error)
	GetByKey(ctx context.Context, roleKey string) (*model.Role, error)
	List(ctx context.Context, opts *ListOptions) ([]*model.Role, int64, error)
	AssignMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error
	GetRoleMenus(ctx context.Context, roleID uint64) ([]*model.Menu, error)
}

// IDeptStore defines department data access operations
type IDeptStore interface {
	Create(ctx context.Context, dept *model.Dept) error
	Update(ctx context.Context, dept *model.Dept) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Dept, error)
	List(ctx context.Context) ([]*model.Dept, error)
	GetChildren(ctx context.Context, parentID uint64) ([]*model.Dept, error)
	GetDeptIDs(ctx context.Context, deptID uint64, includeChildren bool) ([]uint64, error)
	BuildTree(depts []*model.Dept) []*model.Dept
}

// IMenuStore defines menu data access operations
type IMenuStore interface {
	Create(ctx context.Context, menu *model.Menu) error
	Update(ctx context.Context, menu *model.Menu) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*model.Menu, error)
	List(ctx context.Context) ([]*model.Menu, error)
	GetByUserID(ctx context.Context, userID uint64) ([]*model.Menu, error)
	GetChildren(ctx context.Context, parentID uint64) ([]*model.Menu, error)
	BuildTree(menus []*model.Menu) []*model.Menu
	GetMenusByRoleID(ctx context.Context, roleID uint64) ([]*model.Menu, error)
}

// IPermissionStore defines permission data access operations
type IPermissionStore interface {
	GetUserPermissions(ctx context.Context, userID uint64) ([]string, error)
	GetRolePermissions(ctx context.Context, roleID uint64) ([]*model.RolePermission, error)
	CreateRolePermission(ctx context.Context, perm *model.RolePermission) error
	DeleteRolePermissions(ctx context.Context, roleID uint64) error
	BatchCreateRolePermissions(ctx context.Context, perms []*model.RolePermission) error
	GetPermissionsByPattern(ctx context.Context, pattern string) ([]*model.RolePermission, error)
	GetAllPermissions(ctx context.Context) (map[uint64][]string, error)
}

// ListOptions defines common list query options
type ListOptions struct {
	Page     int
	PageSize int
	Filters  map[string]interface{}
	OrderBy  string
}

// DefaultListOptions returns default list options
func DefaultListOptions() *ListOptions {
	return &ListOptions{
		Page:     1,
		PageSize: 10,
		Filters:  make(map[string]interface{}),
	}
}

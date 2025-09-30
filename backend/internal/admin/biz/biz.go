// Package biz defines business logic interfaces
package biz

import (
	"github.com/sword-demon/go-react-admin/internal/admin/biz/dept"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/menu"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/permission"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/role"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/user"
)

// IBiz represents the top-level business logic interface
type IBiz interface {
	Users() user.IUserBiz
	Roles() role.IRoleBiz
	Depts() dept.IDeptBiz
	Menus() menu.IMenuBiz
	Permissions() permission.IPermissionBiz
}

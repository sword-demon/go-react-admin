package biz

import (
	"github.com/sword-demon/go-react-admin/internal/admin/biz/dept"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/menu"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/permission"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/role"
	"github.com/sword-demon/go-react-admin/internal/admin/biz/user"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
)

// bizFactory implements IBiz interface
type bizFactory struct {
	store store.IStore
	cache *cache.RedisClient
}

// NewBiz creates a new biz instance
func NewBiz(store store.IStore, cache *cache.RedisClient) IBiz {
	return &bizFactory{
		store: store,
		cache: cache,
	}
}

// Users returns user biz
func (b *bizFactory) Users() user.IUserBiz {
	return user.NewUserBiz(b.store, b.cache)
}

// Roles returns role biz
func (b *bizFactory) Roles() role.IRoleBiz {
	return role.NewRoleBiz(b.store, b.cache)
}

// Depts returns dept biz
func (b *bizFactory) Depts() dept.IDeptBiz {
	return dept.NewDeptBiz(b.store)
}

// Menus returns menu biz
func (b *bizFactory) Menus() menu.IMenuBiz {
	return menu.NewMenuBiz(b.store)
}

// Permissions returns permission biz
func (b *bizFactory) Permissions() permission.IPermissionBiz {
	return permission.NewPermissionBiz(b.store, b.cache)
}
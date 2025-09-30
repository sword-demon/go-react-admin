package store

import (
	"gorm.io/gorm"
)

// datastore implements IStore interface
type datastore struct {
	db *gorm.DB
}

// NewStore creates a new store instance
func NewStore(db *gorm.DB) IStore {
	return &datastore{db: db}
}

// Users returns user store
func (ds *datastore) Users() IUserStore {
	return newUserStore(ds.db)
}

// Roles returns role store
func (ds *datastore) Roles() IRoleStore {
	return newRoleStore(ds.db)
}

// Depts returns dept store
func (ds *datastore) Depts() IDeptStore {
	return newDeptStore(ds.db)
}

// Menus returns menu store
func (ds *datastore) Menus() IMenuStore {
	return newMenuStore(ds.db)
}

// Permissions returns permission store
func (ds *datastore) Permissions() IPermissionStore {
	return newPermissionStore(ds.db)
}

// Close closes database connection
func (ds *datastore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

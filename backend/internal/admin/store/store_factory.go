package store

import (
	"context"

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

// Transaction executes a function within a database transaction
// Usage example:
//
//	err := store.Transaction(ctx, func(txStore IStore) error {
//	    if err := txStore.Users().Create(ctx, user); err != nil {
//	        return err
//	    }
//	    if err := txStore.Roles().AssignRoles(ctx, user.ID, roleIDs); err != nil {
//	        return err
//	    }
//	    return nil
//	})
func (ds *datastore) Transaction(ctx context.Context, fn func(IStore) error) error {
	return ds.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create a new store instance with the transaction DB
		txStore := NewStore(tx)
		return fn(txStore)
	})
}

// Close closes database connection
func (ds *datastore) Close() error {
	sqlDB, err := ds.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

package store

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// permissionStore implements IPermissionStore interface
type permissionStore struct {
	db *gorm.DB
}

// newPermissionStore creates a new permission store
func newPermissionStore(db *gorm.DB) IPermissionStore {
	return &permissionStore{db: db}
}

// GetUserPermissions retrieves all permission patterns for a user (via roles)
func (s *permissionStore) GetUserPermissions(ctx context.Context, userID uint64) ([]string, error) {
	var permissions []string

	// Get all permission patterns from user's roles
	err := s.db.WithContext(ctx).
		Model(&model.RolePermission{}).
		Select("sys_role_permission.permission_pattern").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_permission.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Where("sys_role_permission.status = ?", model.StatusEnabled).
		Pluck("permission_pattern", &permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetRolePermissions retrieves permission records for a role
func (s *permissionStore) GetRolePermissions(ctx context.Context, roleID uint64) ([]*model.RolePermission, error) {
	var permissions []*model.RolePermission
	err := s.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Where("status = ?", model.StatusEnabled).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// CreateRolePermission creates a new role permission
func (s *permissionStore) CreateRolePermission(ctx context.Context, perm *model.RolePermission) error {
	return s.db.WithContext(ctx).Create(perm).Error
}

// DeleteRolePermissions deletes all permissions for a role
func (s *permissionStore) DeleteRolePermissions(ctx context.Context, roleID uint64) error {
	return s.db.WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&model.RolePermission{}).Error
}

// BatchCreateRolePermissions batch creates role permissions (replaces existing)
func (s *permissionStore) BatchCreateRolePermissions(ctx context.Context, perms []*model.RolePermission) error {
	if len(perms) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing permissions for this role
		roleID := perms[0].RoleID
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing permissions: %w", err)
		}

		// Create new permissions
		if err := tx.Create(&perms).Error; err != nil {
			return fmt.Errorf("failed to create permissions: %w", err)
		}

		return nil
	})
}

// GetPermissionsByPattern retrieves permissions by pattern (for debugging)
func (s *permissionStore) GetPermissionsByPattern(ctx context.Context, pattern string) ([]*model.RolePermission, error) {
	var permissions []*model.RolePermission
	err := s.db.WithContext(ctx).
		Where("permission_pattern = ?", pattern).
		Where("status = ?", model.StatusEnabled).
		Preload("Role").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetAllPermissions retrieves all permission patterns (for cache warming)
func (s *permissionStore) GetAllPermissions(ctx context.Context) (map[uint64][]string, error) {
	var records []struct {
		RoleID            uint64
		PermissionPattern string
	}

	err := s.db.WithContext(ctx).
		Model(&model.RolePermission{}).
		Select("role_id, permission_pattern").
		Where("status = ?", model.StatusEnabled).
		Find(&records).Error

	if err != nil {
		return nil, err
	}

	// Group by role ID
	result := make(map[uint64][]string)
	for _, record := range records {
		result[record.RoleID] = append(result[record.RoleID], record.PermissionPattern)
	}

	return result, nil
}

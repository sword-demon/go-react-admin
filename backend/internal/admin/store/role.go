package store

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// roleStore implements IRoleStore interface
type roleStore struct {
	db *gorm.DB
}

// newRoleStore creates a new role store
func newRoleStore(db *gorm.DB) IRoleStore {
	return &roleStore{db: db}
}

// Create creates a new role
func (s *roleStore) Create(ctx context.Context, role *model.Role) error {
	return s.db.WithContext(ctx).Create(role).Error
}

// Update updates role information
func (s *roleStore) Update(ctx context.Context, role *model.Role) error {
	return s.db.WithContext(ctx).Model(role).Updates(role).Error
}

// Delete soft deletes a role
func (s *roleStore) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Delete(&model.Role{}, id).Error
}

// Get retrieves a role by ID
func (s *roleStore) Get(ctx context.Context, id uint64) (*model.Role, error) {
	var role model.Role
	err := s.db.WithContext(ctx).
		Preload("Menus").
		Preload("Permissions").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByKey retrieves a role by role_key
func (s *roleStore) GetByKey(ctx context.Context, roleKey string) (*model.Role, error) {
	var role model.Role
	err := s.db.WithContext(ctx).
		Where("role_key = ?", roleKey).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List retrieves roles with pagination and filters
func (s *roleStore) List(ctx context.Context, opts *ListOptions) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	// Build base query
	query := s.db.WithContext(ctx).Model(&model.Role{})

	// Apply filters
	if roleName, ok := opts.Filters["role_name"].(string); ok && roleName != "" {
		query = query.Where("role_name LIKE ?", "%"+roleName+"%")
	}
	if status, ok := opts.Filters["status"].(int); ok {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (opts.Page - 1) * opts.PageSize
	query = query.Offset(offset).Limit(opts.PageSize)

	// Apply ordering
	if opts.OrderBy != "" {
		query = query.Order(opts.OrderBy)
	} else {
		query = query.Order("role_sort ASC, created_at DESC")
	}

	// Execute query
	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AssignMenus assigns menus to a role (replaces existing menus)
func (s *roleStore) AssignMenus(ctx context.Context, roleID uint64, menuIDs []uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing menu assignments
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing menus: %w", err)
		}

		// Assign new menus
		if len(menuIDs) > 0 {
			roleMenus := make([]model.RoleMenu, 0, len(menuIDs))
			for _, menuID := range menuIDs {
				roleMenus = append(roleMenus, model.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				})
			}
			if err := tx.Create(&roleMenus).Error; err != nil {
				return fmt.Errorf("failed to assign menus: %w", err)
			}
		}

		return nil
	})
}

// GetRoleMenus retrieves menus for a role
func (s *roleStore) GetRoleMenus(ctx context.Context, roleID uint64) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.WithContext(ctx).
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Where("sys_role_menu.role_id = ?", roleID).
		Where("sys_menu.status = ?", model.StatusEnabled).
		Order("sys_menu.order_num ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

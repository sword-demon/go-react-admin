package store

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// menuStore implements IMenuStore interface
type menuStore struct {
	db *gorm.DB
}

// newMenuStore creates a new menu store
func newMenuStore(db *gorm.DB) IMenuStore {
	return &menuStore{db: db}
}

// Create creates a new menu
func (s *menuStore) Create(ctx context.Context, menu *model.Menu) error {
	return s.db.WithContext(ctx).Create(menu).Error
}

// Update updates menu information
func (s *menuStore) Update(ctx context.Context, menu *model.Menu) error {
	return s.db.WithContext(ctx).Model(menu).Updates(menu).Error
}

// Delete soft deletes a menu
func (s *menuStore) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check if menu has children
		var count int64
		if err := tx.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete menu with children")
		}

		// Delete role-menu associations
		if err := tx.Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error; err != nil {
			return fmt.Errorf("failed to delete role-menu associations: %w", err)
		}

		return tx.Delete(&model.Menu{}, id).Error
	})
}

// Get retrieves a menu by ID
func (s *menuStore) Get(ctx context.Context, id uint64) (*model.Menu, error) {
	var menu model.Menu
	err := s.db.WithContext(ctx).
		Preload("Parent").
		First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// List retrieves all menus (for building tree)
func (s *menuStore) List(ctx context.Context) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.WithContext(ctx).
		Where("status = ?", model.StatusEnabled).
		Order("order_num ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// GetByUserID retrieves menus accessible by user (via roles)
func (s *menuStore) GetByUserID(ctx context.Context, userID uint64) ([]*model.Menu, error) {
	var menus []*model.Menu

	// Get menus through user's roles
	err := s.db.WithContext(ctx).
		Distinct("sys_menu.*").
		Table("sys_menu").
		Joins("JOIN sys_role_menu ON sys_role_menu.menu_id = sys_menu.id").
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role_menu.role_id").
		Where("sys_user_role.user_id = ?", userID).
		Where("sys_menu.status = ?", model.StatusEnabled).
		Where("sys_menu.visible = ?", model.StatusEnabled).
		Order("sys_menu.order_num ASC").
		Find(&menus).Error

	if err != nil {
		return nil, err
	}
	return menus, nil
}

// GetChildren retrieves child menus
func (s *menuStore) GetChildren(ctx context.Context, parentID uint64) ([]*model.Menu, error) {
	var menus []*model.Menu
	err := s.db.WithContext(ctx).
		Where("parent_id = ?", parentID).
		Where("status = ?", model.StatusEnabled).
		Order("order_num ASC").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// BuildTree builds menu tree structure
func (s *menuStore) BuildTree(menus []*model.Menu) []*model.Menu {
	// Create map for quick lookup
	menuMap := make(map[uint64]*model.Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	// Build tree
	var roots []*model.Menu
	for _, menu := range menus {
		if menu.ParentID == 0 {
			roots = append(roots, menu)
		} else if parent, exists := menuMap[menu.ParentID]; exists {
			if parent.Children == nil {
				parent.Children = []*model.Menu{}
			}
			parent.Children = append(parent.Children, menu)
		}
	}

	return roots
}

// GetMenusByRoleID retrieves menus by role ID
func (s *menuStore) GetMenusByRoleID(ctx context.Context, roleID uint64) ([]*model.Menu, error) {
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

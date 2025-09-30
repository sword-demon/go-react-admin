package store

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// userStore implements IUserStore interface
type userStore struct {
	db *gorm.DB
}

// newUserStore creates a new user store
func newUserStore(db *gorm.DB) IUserStore {
	return &userStore{db: db}
}

// Create creates a new user
func (s *userStore) Create(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

// Update updates user information
func (s *userStore) Update(ctx context.Context, user *model.User) error {
	return s.db.WithContext(ctx).Model(user).Updates(user).Error
}

// Delete soft deletes a user
func (s *userStore) Delete(ctx context.Context, id uint64) error {
	return s.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

// Get retrieves a user by ID
func (s *userStore) Get(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).
		Preload("Dept").
		Preload("Roles").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (s *userStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).
		Where("username = ?", username).
		Preload("Dept").
		Preload("Roles").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List retrieves users with pagination and filters
func (s *userStore) List(ctx context.Context, opts *ListOptions) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// Build base query
	query := s.db.WithContext(ctx).Model(&model.User{})

	// Apply filters
	if username, ok := opts.Filters["username"].(string); ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status, ok := opts.Filters["status"].(int); ok {
		query = query.Where("status = ?", status)
	}
	if deptID, ok := opts.Filters["dept_id"].(uint64); ok && deptID > 0 {
		query = query.Where("dept_id = ?", deptID)
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
		query = query.Order("created_at DESC")
	}

	// Preload relations
	query = query.Preload("Dept").Preload("Roles")

	// Execute query
	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserRoles retrieves roles for a user
func (s *userStore) GetUserRoles(ctx context.Context, userID uint64) ([]*model.Role, error) {
	var roles []*model.Role
	err := s.db.WithContext(ctx).
		Joins("JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignRoles assigns roles to a user (replaces existing roles)
func (s *userStore) AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete existing role assignments
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return fmt.Errorf("failed to delete existing roles: %w", err)
		}

		// Assign new roles
		if len(roleIDs) > 0 {
			userRoles := make([]model.UserRole, 0, len(roleIDs))
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, model.UserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return fmt.Errorf("failed to assign roles: %w", err)
			}
		}

		return nil
	})
}

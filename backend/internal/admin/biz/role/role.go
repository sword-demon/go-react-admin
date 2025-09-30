package role

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/admin/biz"
	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"gorm.io/gorm"
)

// roleBiz implements IRoleBiz interface
type roleBiz struct {
	store store.IStore
	cache *cache.RedisClient
}

// NewRoleBiz creates a new role biz
func NewRoleBiz(store store.IStore, cache *cache.RedisClient) biz.IRoleBiz {
	return &roleBiz{
		store: store,
		cache: cache,
	}
}

// Create creates a new role
func (b *roleBiz) Create(ctx context.Context, req *biz.CreateRoleRequest) (*biz.RoleResponse, error) {
	// Check if role_key exists
	_, err := b.store.Roles().GetByKey(ctx, req.RoleKey)
	if err == nil {
		return nil, fmt.Errorf("role_key already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create role model
	role := &model.Role{
		RoleName:  req.RoleName,
		RoleKey:   req.RoleKey,
		RoleSort:  req.Sort,
		DataScope: model.DataScopeSelfOnly, // Default to self only
		Status:    model.StatusEnabled,
	}

	// Save to database
	if err := b.store.Roles().Create(ctx, role); err != nil {
		return nil, err
	}

	return b.toRoleResponse(role), nil
}

// Update updates role information
func (b *roleBiz) Update(ctx context.Context, id uint64, req *biz.UpdateRoleRequest) error {
	role, err := b.store.Roles().Get(ctx, id)
	if err != nil {
		return err
	}

	if req.RoleName != "" {
		role.RoleName = req.RoleName
	}
	if req.Sort > 0 {
		role.RoleSort = req.Sort
	}
	if req.Status >= 0 {
		role.Status = uint8(req.Status)
	}

	return b.store.Roles().Update(ctx, role)
}

// Delete soft deletes a role
func (b *roleBiz) Delete(ctx context.Context, id uint64) error {
	return b.store.Roles().Delete(ctx, id)
}

// Get retrieves role by ID
func (b *roleBiz) Get(ctx context.Context, id uint64) (*biz.RoleResponse, error) {
	role, err := b.store.Roles().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return b.toRoleResponse(role), nil
}

// List retrieves roles with pagination
func (b *roleBiz) List(ctx context.Context, req *biz.ListRoleRequest) (*biz.ListRoleResponse, error) {
	opts := &store.ListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
		Filters:  make(map[string]interface{}),
	}

	if req.RoleName != "" {
		opts.Filters["role_name"] = req.RoleName
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}

	roles, total, err := b.store.Roles().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	items := make([]*biz.RoleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, b.toRoleResponse(role))
	}

	return &biz.ListRoleResponse{
		Total: total,
		Items: items,
	}, nil
}

// AssignPermissions assigns permissions to role
func (b *roleBiz) AssignPermissions(ctx context.Context, roleID uint64, permissionIDs []uint64) error {
	// TODO: Implement permission assignment logic
	return fmt.Errorf("not implemented yet")
}

// toRoleResponse converts model to response
func (b *roleBiz) toRoleResponse(role *model.Role) *biz.RoleResponse {
	return &biz.RoleResponse{
		ID:       role.ID,
		RoleName: role.RoleName,
		RoleKey:  role.RoleKey,
		Sort:     role.RoleSort,
		Status:   int8(role.Status),
	}
}

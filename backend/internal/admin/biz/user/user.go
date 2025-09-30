package user

import (
	"context"
	"fmt"

	"github.com/sword-demon/go-react-admin/internal/admin/store"
	"github.com/sword-demon/go-react-admin/internal/pkg/cache"
	"github.com/sword-demon/go-react-admin/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userBiz implements IUserBiz interface
type userBiz struct {
	store store.IStore
	cache *cache.RedisClient
}

// IUserBiz defines user business logic operations
type IUserBiz interface {
	Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
	Update(ctx context.Context, id uint64, req *UpdateUserRequest) error
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*UserResponse, error)
	List(ctx context.Context, req *ListUserRequest) (*ListUserResponse, error)
	ChangePassword(ctx context.Context, id uint64, req *ChangePasswordRequest) error
	AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error
}

// Request/Response structs
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeptID   uint64 `json:"dept_id"`
}

type UpdateUserRequest struct {
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeptID   uint64 `json:"dept_id"`
	Status   int8   `json:"status"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserResponse struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	DeptID   uint64 `json:"dept_id"`
	Status   int8   `json:"status"`
}

type ListUserRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	Status   *int8  `form:"status"`
}

type ListUserResponse struct {
	Total int64           `json:"total"`
	Items []*UserResponse `json:"items"`
}

// NewUserBiz creates a new user biz
func NewUserBiz(store store.IStore, cache *cache.RedisClient) IUserBiz {
	return &userBiz{
		store: store,
		cache: cache,
	}
}

// Create creates a new user with hashed password
func (b *userBiz) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	// Check if username exists
	_, err := b.store.Users().GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		NickName: req.NickName,
		Email:    req.Email,
		Phone:    req.Phone,
		DeptID:   req.DeptID,
		Status:   model.StatusEnabled,
	}

	// Save to database
	if err := b.store.Users().Create(ctx, user); err != nil {
		return nil, err
	}

	return b.toUserResponse(user), nil
}

// Update updates user information
func (b *userBiz) Update(ctx context.Context, id uint64, req *UpdateUserRequest) error {
	// Get existing user
	user, err := b.store.Users().Get(ctx, id)
	if err != nil {
		return err
	}

	// Update fields
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.DeptID > 0 {
		user.DeptID = req.DeptID
	}
	if req.Status >= 0 {
		user.Status = uint8(req.Status)
	}

	return b.store.Users().Update(ctx, user)
}

// Delete soft deletes a user
func (b *userBiz) Delete(ctx context.Context, id uint64) error {
	// Check if user exists
	if _, err := b.store.Users().Get(ctx, id); err != nil {
		return err
	}

	return b.store.Users().Delete(ctx, id)
}

// Get retrieves user by ID
func (b *userBiz) Get(ctx context.Context, id uint64) (*UserResponse, error) {
	user, err := b.store.Users().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return b.toUserResponse(user), nil
}

// List retrieves users with pagination
func (b *userBiz) List(ctx context.Context, req *ListUserRequest) (*ListUserResponse, error) {
	// Build list options
	opts := &store.ListOptions{
		Page:     req.Page,
		PageSize: req.PageSize,
		Filters:  make(map[string]interface{}),
	}

	if req.Username != "" {
		opts.Filters["username"] = req.Username
	}
	if req.Status != nil {
		opts.Filters["status"] = int(*req.Status)
	}

	// Set defaults
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}

	// Query users
	users, total, err := b.store.Users().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Convert to responses
	items := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, b.toUserResponse(user))
	}

	return &ListUserResponse{
		Total: total,
		Items: items,
	}, nil
}

// ChangePassword changes user password
func (b *userBiz) ChangePassword(ctx context.Context, id uint64, req *ChangePasswordRequest) error {
	// Get user
	user, err := b.store.Users().Get(ctx, id)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("old password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	return b.store.Users().Update(ctx, user)
}

// AssignRoles assigns roles to user
func (b *userBiz) AssignRoles(ctx context.Context, userID uint64, roleIDs []uint64) error {
	// Check if user exists
	if _, err := b.store.Users().Get(ctx, userID); err != nil {
		return err
	}

	return b.store.Users().AssignRoles(ctx, userID, roleIDs)
}

// toUserResponse converts model to response
func (b *userBiz) toUserResponse(user *model.User) *UserResponse {
	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		NickName: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		DeptID:   user.DeptID,
		Status:   int8(user.Status),
	}
}

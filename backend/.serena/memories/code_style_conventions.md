# Code Style and Conventions

## General Go Style
- Follow **golang-standards/project-layout**
- Use `gofmt` for formatting
- Use `go vet` for linting
- Interface-driven design (dependency injection)

## Naming Conventions

### Interfaces
- Prefix with `I`: `IBiz`, `IUserBiz`, `IStore`, `IUserStore`
- Define in separate files: `biz/biz.go`, `store/store.go`

### Structs
- PascalCase for exported: `User`, `Role`, `CreateUserRequest`
- camelCase for unexported: `userBiz`, `userStore`

### Functions/Methods
- PascalCase for exported: `Create()`, `GetByUsername()`
- camelCase for unexported: `hashPassword()`, `buildQuery()`

### Variables
- camelCase: `userID`, `roleKey`, `dataScope`
- Constants: UPPER_SNAKE_CASE or PascalCase for exported

### File Names
- snake_case: `user_biz.go`, `role_store.go`, `permission_middleware.go`

## Package Structure
```
package store

// Interface first
type IUserStore interface {
    Create(ctx context.Context, user *model.User) error
    Get(ctx context.Context, id uint64) (*model.User, error)
}

// Implementation
type userStore struct {
    db *gorm.DB
}

// Constructor
func NewUserStore(db *gorm.DB) IUserStore {
    return &userStore{db: db}
}

// Methods
func (s *userStore) Create(ctx context.Context, user *model.User) error {
    return s.db.WithContext(ctx).Create(user).Error
}
```

## Error Handling
- Use `pkg/errors` for wrapped errors
- Return errors, don't panic
- HTTP layer handles status codes

```go
// Bad
if err != nil {
    panic(err)
}

// Good
if err != nil {
    return nil, errors.Wrap(err, "failed to create user")
}
```

## Context Usage
- Always pass `context.Context` as first parameter
- Use for cancellation, timeouts, request-scoped values
- Pass to all database operations

```go
func (s *userStore) Get(ctx context.Context, id uint64) (*model.User, error) {
    var user model.User
    if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

## Struct Tags
- JSON: `json:"username"`
- GORM: `gorm:"column:username;type:varchar(64);not null"`
- Validation: `binding:"required,min=4,max=64"`

```go
type CreateUserRequest struct {
    Username string `json:"username" binding:"required,min=4,max=64"`
    Password string `json:"password" binding:"required,min=6"`
    NickName string `json:"nick_name" binding:"max=64"`
    Email    string `json:"email" binding:"omitempty,email"`
}
```

## Comments
- Godoc style for exported symbols
- Start with symbol name
- End with period

```go
// IUserBiz defines user business logic operations.
// It orchestrates user CRUD, password management, and role assignments.
type IUserBiz interface {
    // Create creates a new user with hashed password.
    Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
}
```

## Testing (Future)
- Files: `*_test.go`
- Function names: `TestUserBiz_Create`
- Use table-driven tests
- Mock interfaces with `go.uber.org/mock`
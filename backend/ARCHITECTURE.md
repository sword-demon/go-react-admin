# Backend Architecture Documentation

## Directory Structure (Following miniblog pattern)

```
backend/
├── cmd/
│   └── server/                 # Application entry point
│       └── main.go             # Server startup
│
├── internal/                   # Internal application code (cannot be imported by external projects)
│   ├── admin/                  # Admin module (core application)
│   │   ├── biz/                # Business logic layer
│   │   │   ├── user/           # User business logic
│   │   │   ├── role/           # Role business logic
│   │   │   ├── dept/           # Department business logic
│   │   │   ├── menu/           # Menu business logic
│   │   │   └── biz.go          # Biz interface definition
│   │   │
│   │   ├── controller/         # HTTP handlers
│   │   │   └── v1/             # API version 1
│   │   │       ├── user.go
│   │   │       ├── role.go
│   │   │       ├── dept.go
│   │   │       ├── menu.go
│   │   │       ├── auth.go
│   │   │       └── controller.go
│   │   │
│   │   ├── store/              # Data access layer (DAO)
│   │   │   ├── user.go         # User CRUD
│   │   │   ├── role.go         # Role CRUD
│   │   │   ├── dept.go         # Department CRUD
│   │   │   ├── menu.go         # Menu CRUD
│   │   │   ├── permission.go   # Permission CRUD
│   │   │   └── store.go        # Store interface
│   │   │
│   │   ├── middleware/         # Application-specific middleware
│   │   │   ├── auth.go         # JWT authentication
│   │   │   ├── permission.go   # Permission validation
│   │   │   └── datascope.go    # Data permission filter
│   │   │
│   │   ├── router.go           # Route registration
│   │   └── admin.go            # Application initialization
│   │
│   └── pkg/                    # Internal shared packages
│       ├── auth/               # Authentication utilities
│       ├── cache/              # Cache wrapper (Redis)
│       │   ├── redis.go
│       │   └── permission.go   # Permission cache (3-layer)
│       ├── errors/             # Custom error types
│       ├── log/                # Logging wrapper
│       ├── model/              # GORM models
│       │   ├── user.go
│       │   ├── role.go
│       │   ├── dept.go
│       │   └── ...
│       ├── middleware/         # Common middleware
│       │   ├── logger.go
│       │   ├── recovery.go
│       │   └── cors.go
│       └── utils/              # Utility functions
│           ├── jwt.go
│           ├── bcrypt.go
│           └── validator.go
│
├── pkg/                        # Public packages (can be imported externally)
│   ├── core/                   # Core response wrapper
│   │   ├── response.go         # Unified response format
│   │   └── pagination.go       # Pagination helper
│   ├── token/                  # JWT token operations
│   │   └── jwt.go
│   └── validator/              # Request validation
│       └── validator.go
│
├── configs/                    # Configuration files
│   ├── config.yaml             # Main config (DO NOT commit!)
│   └── config.example.yaml     # Example config
│
├── go.mod                      # Go module definition
├── go.sum
├── Makefile                    # Build automation
└── README.md
```

## Architecture Layers

### 1. Controller Layer (HTTP Handlers)
- **Location**: `internal/admin/controller/v1/`
- **Responsibility**: Handle HTTP requests, validate input, call business logic
- **Example**:
```go
// internal/admin/controller/v1/user.go
type UserController struct {
    biz biz.IBiz
}

func (c *UserController) Create(ctx *gin.Context) {
    var req CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        core.WriteResponse(ctx, errors.ErrBind, nil)
        return
    }

    user, err := c.biz.Users().Create(ctx, &req)
    if err != nil {
        core.WriteResponse(ctx, err, nil)
        return
    }

    core.WriteResponse(ctx, nil, user)
}
```

### 2. Business Logic Layer (Biz)
- **Location**: `internal/admin/biz/`
- **Responsibility**: Business logic, orchestrate multiple stores, permission checks
- **Example**:
```go
// internal/admin/biz/user/user.go
type UserBiz struct {
    store store.IStore
}

func (b *UserBiz) Create(ctx context.Context, req *CreateUserRequest) (*model.User, error) {
    // Business logic: validate, check duplicates, apply data permissions
    exists, _ := b.store.Users().GetByUsername(ctx, req.Username)
    if exists != nil {
        return nil, errors.ErrUserAlreadyExists
    }

    // Hash password
    hashedPassword, _ := bcrypt.HashPassword(req.Password)

    user := &model.User{
        Username: req.Username,
        Password: hashedPassword,
        // ...
    }

    return b.store.Users().Create(ctx, user)
}
```

### 3. Data Access Layer (Store)
- **Location**: `internal/admin/store/`
- **Responsibility**: Database CRUD operations, GORM queries
- **Example**:
```go
// internal/admin/store/user.go
type UserStore struct {
    db *gorm.DB
}

func (s *UserStore) Create(ctx context.Context, user *model.User) error {
    return s.db.WithContext(ctx).Create(user).Error
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    err := s.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    return &user, err
}
```

## Dependency Flow

```
main.go
  ↓
router.go (register routes)
  ↓
controller (HTTP handling)
  ↓
biz (business logic)
  ↓
store (data access)
  ↓
database (MySQL)
```

## Key Interfaces

### Biz Interface
```go
// internal/admin/biz/biz.go
type IBiz interface {
    Users() IUserBiz
    Roles() IRoleBiz
    Depts() IDeptBiz
    Menus() IMenuBiz
}
```

### Store Interface
```go
// internal/admin/store/store.go
type IStore interface {
    Users() IUserStore
    Roles() IRoleStore
    Depts() IDeptStore
    Menus() IMenuStore
}
```

## Permission Middleware Flow

```
Request → JWT Middleware → Permission Middleware → Handler
            ↓                    ↓
        Parse Token         Check Permission
                            (3-layer cache:
                             Local → Redis → MySQL)
```

## Configuration Management

Using Viper:
```yaml
# configs/config.yaml
server:
  port: 8080
  mode: debug  # debug, release

database:
  host: localhost
  port: 3306
  database: go_react_admin
  username: root
  password: your_password

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your_jwt_secret_key
  expire_hours: 24
```

## Migration from Old Structure

Old → New mapping:
- `controller/` → `internal/admin/controller/v1/`
- `service/` → `internal/admin/biz/`
- `model/` → `internal/pkg/model/`
- `middleware/` → `internal/admin/middleware/` (app-specific)
- `middleware/` → `internal/pkg/middleware/` (common)

## Why This Structure?

1. **Clear Separation**: Controller, Biz, Store layers are clearly separated
2. **Testability**: Each layer can be tested independently with interfaces
3. **Scalability**: Easy to add new modules (just add to biz/, store/, controller/)
4. **Standard**: Follows golang-standards/project-layout
5. **Security**: internal/ prevents external packages from importing sensitive code
6. **Maintainability**: Clear dependency direction (Controller → Biz → Store)

## Development Workflow

1. Define model in `internal/pkg/model/`
2. Implement store interface in `internal/admin/store/`
3. Implement business logic in `internal/admin/biz/`
4. Implement HTTP handler in `internal/admin/controller/v1/`
5. Register route in `internal/admin/router.go`
6. Test each layer independently
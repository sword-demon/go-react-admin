# Backend Directory Migration Guide

## ✅ Migration Completed!

The backend directory structure has been refactored following the **miniblog** pattern (golang-standards/project-layout).

---

## 📊 Before vs After

### Old Structure (Basic)
```
backend/
├── controller/     # HTTP handlers
├── service/        # Business logic
├── model/          # GORM models
├── middleware/     # Middleware
├── config/         # Config files
└── main.go         # Entry point
```

### New Structure (Enterprise-grade)
```
backend/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go
│
├── internal/                # Internal packages (cannot be imported externally)
│   ├── admin/               # Admin application module
│   │   ├── biz/             # Business logic layer
│   │   │   ├── user/
│   │   │   ├── role/
│   │   │   ├── dept/
│   │   │   ├── menu/
│   │   │   └── biz.go       # Interface definitions
│   │   ├── controller/      # HTTP handlers
│   │   │   └── v1/          # API version 1
│   │   ├── store/           # Data access layer
│   │   │   └── store.go     # Interface definitions
│   │   └── middleware/      # App-specific middleware
│   └── pkg/                 # Internal shared packages
│       ├── model/           # GORM models
│       ├── cache/           # Cache wrapper (Redis)
│       ├── auth/            # Auth utilities
│       ├── errors/          # Custom errors
│       ├── log/             # Logging wrapper
│       ├── middleware/      # Common middleware
│       └── utils/           # Utility functions
│
├── pkg/                     # Public packages (can be imported externally)
│   ├── core/                # Core response wrapper
│   ├── token/               # JWT utilities
│   └── validator/           # Request validation
│
├── configs/                 # Configuration files
│   └── config.example.yml
│
├── ARCHITECTURE.md          # Architecture documentation
├── go.mod
└── go.sum
```

---

## 🎯 Key Improvements

### 1. Clear Layer Separation
```
Controller → Biz → Store → Database
(HTTP)     (Logic) (DAO)
```

- **Controller**: Handle HTTP requests, validate input
- **Biz**: Business logic, orchestrate stores, permission checks
- **Store**: Database CRUD operations

### 2. Interface-Driven Design
```go
// Dependency injection via interfaces
type UserController struct {
    biz biz.IBiz  // Not concrete implementation!
}

// Easy to mock for testing
func TestUserController(t *testing.T) {
    mockBiz := &MockBiz{}  // Implement IBiz interface
    controller := &UserController{biz: mockBiz}
    // Test...
}
```

### 3. internal/ Package Protection
- Code in `internal/` **cannot** be imported by external projects
- Ensures API stability (only `pkg/` is public)
- Protects sensitive logic

### 4. Versioned API (v1/)
- Future-proof: Easy to add v2 without breaking v1
- Clean routing: `/api/v1/users`, `/api/v2/users`

---

## 📝 Files Created

1. ✅ `ARCHITECTURE.md` - Complete architecture documentation
2. ✅ `cmd/server/main.go` - Application entry point (updated with TODOs)
3. ✅ `internal/admin/biz/biz.go` - Business logic interfaces
4. ✅ `internal/admin/store/store.go` - Data access interfaces
5. ✅ `configs/config.example.yml` - Example configuration

---

## 🚀 Next Steps

### Week 1-2 Tasks (Current):

1. **Database Connection**
   ```go
   // internal/pkg/db/db.go
   func InitDB(cfg *config.Database) (*gorm.DB, error) {
       dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
           cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
       return gorm.Open(mysql.Open(dsn), &gorm.Config{})
   }
   ```

2. **Implement Models**
   - `internal/pkg/model/user.go`
   - `internal/pkg/model/role.go`
   - `internal/pkg/model/dept.go`
   - `internal/pkg/model/menu.go`

3. **Implement Store Layer**
   - `internal/admin/store/user.go` (implement IUserStore)
   - `internal/admin/store/role.go` (implement IRoleStore)
   - Create `store.New()` constructor

4. **Implement Biz Layer**
   - `internal/admin/biz/user/user.go` (implement IUserBiz)
   - `internal/admin/biz/role/role.go` (implement IRoleBiz)
   - Create `biz.New()` constructor

5. **Implement Controllers**
   - `internal/admin/controller/v1/user.go`
   - `internal/admin/controller/v1/auth.go`

6. **Setup Router**
   - `internal/admin/router.go`
   - Register all routes with middleware

---

## 🔧 Development Commands

```bash
# Run server (from backend/)
go run cmd/server/main.go

# Or with hot reload (install air first)
go install github.com/air-verse/air@latest
air

# Test ping endpoint
curl http://localhost:8080/ping

# Run tests
go test ./...

# Run specific package tests
go test -v ./internal/admin/biz/user

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

---

## 📚 Architecture References

- **Followed**: [marmotedu/miniblog](https://github.com/marmotedu/miniblog)
- **Standard**: [golang-standards/project-layout](https://github.com/golang-standards/project-layout)
- **See**: `ARCHITECTURE.md` for detailed architecture documentation

---

## ⚠️ Migration Notes

### Mapping Old to New:

| Old Location | New Location | Purpose |
|-------------|--------------|---------|
| `controller/` | `internal/admin/controller/v1/` | HTTP handlers |
| `service/` | `internal/admin/biz/` | Business logic |
| `model/` | `internal/pkg/model/` | GORM models |
| `middleware/` (app-specific) | `internal/admin/middleware/` | Auth, permission |
| `middleware/` (common) | `internal/pkg/middleware/` | Logger, recovery |
| `config/` | `configs/` | Configuration files |
| `main.go` | `cmd/server/main.go` | Entry point |

### Breaking Changes:
- Import paths changed (update all imports!)
- Old: `"github.com/sword-demon/go-react-admin/model"`
- New: `"github.com/sword-demon/go-react-admin/internal/pkg/model"`

### Compatibility:
- Old directories still exist (will be removed after migration)
- No database schema changes
- API endpoints remain the same

---

## 🎯 Why This Structure?

1. **Scalability**: Easy to add new modules (just add to biz/, store/, controller/)
2. **Testability**: Interface-driven design makes testing easy
3. **Maintainability**: Clear separation of concerns
4. **Standard**: Follows Go community best practices
5. **Security**: internal/ protects sensitive code
6. **Team-friendly**: Clear structure for team collaboration

---

## 📖 Additional Reading

- Read `ARCHITECTURE.md` for complete layer explanations
- Read `internal/admin/biz/biz.go` for all interface definitions
- Read `internal/admin/store/store.go` for data access interfaces
- Refer to miniblog source code for implementation examples

---

**Status**: ✅ Directory structure refactored, ready for development!
**Next**: Implement database connection and models (Week 1-2)
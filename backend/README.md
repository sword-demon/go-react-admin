# Backend - Go-React-Admin

Enterprise-grade backend architecture following [miniblog](https://github.com/marmotedu/miniblog) pattern.

## 📁 Directory Structure

```
backend/
├── cmd/server/             # Application entry point
├── internal/               # Internal application code
│   ├── admin/              # Admin module (core app)
│   │   ├── biz/            # Business logic layer (interface: IBiz)
│   │   ├── controller/v1/  # HTTP handlers (API v1)
│   │   ├── store/          # Data access layer (interface: IStore)
│   │   └── middleware/     # App-specific middleware
│   └── pkg/                # Internal shared packages
│       ├── model/          # GORM models
│       ├── cache/          # Redis cache (3-layer)
│       └── ...
├── pkg/                    # Public packages
│   ├── core/               # Response wrapper
│   ├── token/              # JWT utilities
│   └── validator/          # Request validation
└── configs/                # Configuration files
```

## 🏗️ Architecture Layers

```
Controller → Biz → Store → Database
(HTTP)     (Logic)(DAO)
```

- **Controller**: Handle HTTP requests, validate input
- **Biz**: Business logic, orchestrate stores, permission checks
- **Store**: Database CRUD operations (GORM)

## 🚀 Quick Start

```bash
# Install dependencies
go mod tidy

# Run server
go run cmd/server/main.go

# Or with hot reload
air

# Test
curl http://localhost:8080/ping
```

## 📚 Documentation

- **ARCHITECTURE.md** - Complete architecture documentation
- **MIGRATION.md** - Migration guide from old structure
- **docs/schema.sql** - Database schema (10 tables)

## 🔑 Key Interfaces

### Biz Interface (Business Logic)
```go
type IBiz interface {
    Users() IUserBiz
    Roles() IRoleBiz
    Depts() IDeptBiz
    Menus() IMenuBiz
}
```

### Store Interface (Data Access)
```go
type IStore interface {
    Users() IUserStore
    Roles() IRoleStore
    Depts() IDeptStore
    Menus() IMenuStore
}
```

## 🎯 Development Workflow

1. Define model in `internal/pkg/model/`
2. Implement store in `internal/admin/store/`
3. Implement biz logic in `internal/admin/biz/`
4. Implement controller in `internal/admin/controller/v1/`
5. Register route in `internal/admin/router.go`

## 📦 Core Dependencies

- **gin-gonic/gin** - Web framework
- **gorm.io/gorm** - ORM
- **redis/go-redis** - Redis client
- **golang-jwt/jwt** - JWT authentication

## 🔐 Permission Flow

```
Request → JWT Middleware → Permission Middleware → Handler
            ↓                    ↓
        Parse Token         Check Permission
                            (3-layer cache:
                             Local → Redis → MySQL)
```

## 📈 Status

- [x] Directory structure (Week 1)
- [x] Interface definitions
- [ ] Database connection (Week 1-2)
- [ ] Implement models (Week 1-2)
- [ ] Implement store layer (Week 2-3)
- [ ] Implement biz layer (Week 3-4)
- [ ] Implement controllers (Week 4-5)
- [ ] Middleware (JWT, Permission) (Week 5)

See main project README for complete roadmap.

## 📖 References

- [marmotedu/miniblog](https://github.com/marmotedu/miniblog) - Architecture reference
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) - Layout standard
- [Gin Documentation](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)

---

**Architecture**: Based on miniblog (enterprise-grade)  
**Status**: ✅ Structure complete, ready for development  
**Next**: Implement database connection and models

# Directory Structure

## Overview
```
backend/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go          # Server startup
│
├── internal/                # Internal packages (cannot be imported externally)
│   ├── admin/               # Core application module
│   │   ├── biz/             # Business logic layer (interface: IBiz)
│   │   │   ├── user/        # User business logic
│   │   │   ├── role/        # Role business logic
│   │   │   ├── dept/        # Department business logic
│   │   │   ├── menu/        # Menu business logic
│   │   │   └── biz.go       # Interface definitions
│   │   ├── controller/      # HTTP handlers
│   │   │   └── v1/          # API version 1
│   │   ├── store/           # Data access layer (interface: IStore)
│   │   │   └── store.go     # Interface definitions
│   │   └── middleware/      # App-specific middleware
│   │
│   └── pkg/                 # Internal shared packages
│       ├── model/           # GORM models
│       ├── cache/           # Redis cache wrapper
│       ├── auth/            # JWT auth utilities
│       ├── errors/          # Custom error types
│       ├── log/             # Logging wrapper
│       ├── middleware/      # Common middleware
│       └── utils/           # Utility functions
│
├── pkg/                     # Public packages (can be imported externally)
│   ├── core/                # Core response wrapper
│   ├── token/               # JWT token utilities
│   └── validator/           # Request validation
│
├── configs/                 # Configuration files
│   └── config.example.yml   # Example config (copy to config.yml)
│
├── ARCHITECTURE.md          # Complete architecture documentation
├── MIGRATION.md             # Migration guide from old structure
├── README.md                # Quick start guide
├── go.mod                   # Go module dependencies
└── go.sum                   # Dependency checksums
```

## Layer Responsibilities

### Controller Layer (`internal/admin/controller/v1/`)
- Handle HTTP requests
- Validate input (gin binding)
- Call Biz layer methods
- Return JSON responses

### Biz Layer (`internal/admin/biz/`)
- Business logic orchestration
- Permission checks
- Orchestrate multiple Store operations
- Data transformation

### Store Layer (`internal/admin/store/`)
- Database CRUD operations
- GORM queries
- Data scope filtering (WHERE clause injection)
- Transaction management

## Key Files
- `cmd/server/main.go` - Server entry point (currently with TODOs)
- `internal/admin/biz/biz.go` - All business logic interfaces
- `internal/admin/store/store.go` - All data access interfaces
- `configs/config.example.yml` - Configuration template
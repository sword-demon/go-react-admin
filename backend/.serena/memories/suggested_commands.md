# Suggested Commands

## Development Workflow

### Initial Setup
```bash
# Install dependencies
go mod tidy

# Copy config file
cp configs/config.example.yml configs/config.yml
# Edit configs/config.yml with your DB/Redis credentials

# Install hot reload tool (optional)
go install github.com/air-verse/air@latest
```

### Running the Server
```bash
# Standard run (from backend/)
go run cmd/server/main.go

# With hot reload (recommended)
air

# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server
```

### Testing
```bash
# Run all tests
go test ./...

# Run specific package tests
go test -v ./internal/admin/biz/user

# Run specific test function
go test -v -run TestUserBiz_Create ./internal/admin/biz/user

# Test with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Format code (auto-fix)
go fmt ./...

# Lint code
go vet ./...

# Check for common mistakes (requires staticcheck)
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

### Database
```bash
# Connect to MySQL
mysql -u root -p

# Import schema (from project root)
mysql -u root -p go_react_admin < docs/schema.sql

# Verify tables
mysql -u root -p -e "USE go_react_admin; SHOW TABLES;"

# Backup database
mysqldump -u root -p go_react_admin > backup_$(date +%Y%m%d).sql
```

### Redis
```bash
# Start Redis server
redis-server

# Connect to Redis CLI
redis-cli

# Check Redis is running
redis-cli ping  # Should return PONG

# Flush all cache (dangerous!)
redis-cli FLUSHALL
```

### Git Workflow
```bash
# Check status
git status

# Add files
git add .

# Commit (use conventional commits)
git commit -m "feat: implement user CRUD operations"
git commit -m "fix: resolve JWT token expiration issue"

# Push
git push origin main
```

### Health Check
```bash
# Test ping endpoint
curl http://localhost:8080/ping

# Test login (when implemented)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Test with JWT (when implemented)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/api/v1/users
```

## System Utilities (macOS Darwin)

### File Operations
```bash
# List files
ls -la

# Find files by name
find . -name "*.go"

# Search in files (ripgrep recommended)
grep -r "IUserBiz" internal/

# Count lines of code
find . -name "*.go" | xargs wc -l
```

### Process Management
```bash
# Find process using port 8080
lsof -i :8080

# Kill process on port 8080
kill -9 $(lsof -t -i:8080)

# Check Go processes
ps aux | grep "go run"
```

### Disk Usage
```bash
# Check disk space
df -h

# Directory size
du -sh backend/

# Clean Go cache
go clean -cache
go clean -modcache
```

## When Task is Completed

1. **Format code**: `go fmt ./...`
2. **Run linter**: `go vet ./...`
3. **Run tests**: `go test ./...`
4. **Test manually**: `curl http://localhost:8080/ping`
5. **Commit changes**: `git add . && git commit -m "feat: ..."`

## Performance Profiling (Advanced)
```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.

# Memory profiling
go test -memprofile=mem.prof -bench=.

# View profile
go tool pprof cpu.prof
```
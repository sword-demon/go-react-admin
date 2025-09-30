# Task Completion Checklist

## Before Submitting Code

### 1. Code Quality Checks ✅
```bash
# Format code (MUST run)
go fmt ./...

# Lint for common issues (MUST run)
go vet ./...

# Optional: Static analysis
staticcheck ./...
```

### 2. Testing ✅
```bash
# Run all tests (if tests exist)
go test ./...

# Check for compilation errors
go build -o /dev/null cmd/server/main.go
```

### 3. Manual Verification ✅
```bash
# Start server
go run cmd/server/main.go

# Test health endpoint
curl http://localhost:8080/ping

# Test new feature endpoint (replace with actual endpoint)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'
```

### 4. Documentation ✅
- [ ] Update comments for new functions
- [ ] Update ARCHITECTURE.md if structure changed
- [ ] Update README.md if commands changed
- [ ] Add TODO comments for incomplete work

### 5. Git Workflow ✅
```bash
# Check what changed
git status
git diff

# Stage changes
git add .

# Commit with conventional commit message
git commit -m "feat(user): implement user CRUD operations"
# Or:
# git commit -m "fix(auth): resolve JWT token expiration"
# git commit -m "refactor(biz): extract common validation logic"
# git commit -m "docs(readme): update API examples"

# Push (if remote exists)
git push origin main
```

## Conventional Commit Types
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code refactoring (no behavior change)
- `docs`: Documentation only
- `style`: Code style (formatting, no logic change)
- `test`: Add/update tests
- `chore`: Build/tooling changes

## Examples
```bash
# Good commit messages
git commit -m "feat(user): add user list API with pagination"
git commit -m "fix(auth): prevent duplicate login sessions"
git commit -m "refactor(store): extract common query builder"

# Bad commit messages
git commit -m "update"
git commit -m "fix bug"
git commit -m "WIP"
```

## Development Phase Checks

### Phase 1 (MVP - Weeks 1-8)
When completing a feature:
1. ✅ CRUD operations work
2. ✅ Error handling implemented
3. ✅ Validation in place
4. ✅ No hardcoded values
5. ✅ TODO comments for future work

### Phase 2 (Extensions)
When completing advanced features:
1. ✅ All Phase 1 checks
2. ✅ Logging added
3. ✅ Cache invalidation works
4. ✅ Performance tested
5. ✅ Security reviewed

## Single Developer Reminders

### Daily Checklist
- [ ] Commit at least once daily
- [ ] Test before committing
- [ ] Don't leave broken code overnight
- [ ] Document complex logic immediately

### Weekly Checklist (Fridays)
- [ ] Review progress vs timeline
- [ ] Clean up TODOs
- [ ] Backup database
- [ ] Update project README

### When Stuck (>2 hours)
1. ❌ Don't keep trying the same approach
2. ✅ Search Stack Overflow / GitHub Issues
3. ✅ Ask ChatGPT/Claude for alternative solutions
4. ✅ Simplify the problem (cut scope if needed)
5. ✅ Take a break (seriously!)

### When Falling Behind Schedule
- **Week 4**: Cut API permission pattern matching
- **Week 6**: Cut data scope filtering
- **Week 7**: Cut import/export, logs

**Rule**: Better to ship 80% working than 50% perfect!
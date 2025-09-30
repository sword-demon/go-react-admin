# Go-React-Admin Backend Project Overview

## Project Purpose
Enterprise-grade RBAC permission management system backend, supporting 500+ API interfaces with:
- Role-Based Access Control (RBAC) + Data Scope filtering
- API permission pattern matching (wildcards to reduce 5000+ records to 50-100)
- Three-layer cache architecture (Local Memory → Redis → MySQL)
- Designed for single developer, 8-week MVP timeline

## Target Scale
- Users: 100-500 (small-to-medium enterprises)
- APIs: 500+ endpoints
- Roles: 5-15
- Single Redis instance (not cluster)

## Tech Stack
- **Language**: Go 1.24.0
- **Web Framework**: Gin v1.11.0
- **ORM**: GORM v1.31.0 + MySQL driver v1.6.0
- **Cache**: Redis v9.14.0 (go-redis)
- **Auth**: JWT (golang-jwt/jwt)
- **Database**: MySQL 8.0+

## Key Dependencies
```
github.com/gin-gonic/gin v1.11.0
gorm.io/gorm v1.31.0
gorm.io/driver/mysql v1.6.0
github.com/redis/go-redis/v9 v9.14.0
github.com/go-playground/validator/v10 v10.27.0
golang.org/x/crypto v0.42.0 (bcrypt)
```

## Architecture Pattern
Based on **marmotedu/miniblog** enterprise pattern:
```
Controller → Biz → Store → Database
(HTTP)     (Logic) (DAO)
```

## Core Features
1. User/Role/Dept/Menu management (CRUD)
2. JWT authentication
3. Permission middleware with 3-layer cache
4. API permission pattern matching
5. Data scope filtering (ALL/DEPT_AND_CHILD/DEPT_ONLY/SELF_ONLY)
6. Audit logs (Phase 2)
7. Data import/export (Phase 2)

## Module: github.com/sword-demon/go-react-admin
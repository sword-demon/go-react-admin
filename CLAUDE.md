# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## ⚠️ 项目真实状态 (2025-10-01)

**架构进度**: 🟡 DDD三层架构设计完成,但**仅完成 30% 实现**

### 后端实现进度
✅ **已完成** (可运行):
- Model层定义 (User/Role/Dept/Menu/Permission/Log)
- Store层接口和实现 (数据库CRUD)
- Biz层接口和部分实现 (业务逻辑)
- 配置管理 (Viper + YAML)
- 数据库连接池 (GORM v1.31.0 + MySQL)
- Redis基础客户端
- Makefile构建脚本

❌ **未实现** (阻塞项目启动):
- Controller层 HTTP handlers (internal/admin/controller/v1/ 目录为空!)
- 路由注册 (router.go 不存在,main.go 无路由)
- JWT认证中间件
- 权限验证中间件
- 三层权限缓存 (只有Redis基础客户端,无缓存逻辑)

**后端当前状态**: 能连数据库,但**没有一个可访问的HTTP API**

### 前端实现进度
✅ **已完成**:
- Vite 7 + React 19 + TypeScript 脚手架
- 基础项目结构

❌ **未实现**:
- React Router (连路由库都没装!)
- Ant Design 5 (package.json 里没有!)
- Axios/状态管理库
- 任何业务页面组件

**前端当前状态**: 只是 `npm create vite` 生成的空项目

### 🎯 当务之急 (Week 1 任务)
1. **实现第一个可用的API**: `POST /api/v1/auth/login` (用户登录)
2. **实现第一个CRUD API**: `GET /api/v1/users` (用户列表)
3. **前端能调通后端API**: 安装 axios,实现登录表单
4. **补全基础依赖**: 安装 React Router + Ant Design

---

## 项目简介

企业级 RBAC 权限管理系统 (单人全栈项目,调整为 **12 周 MVP 计划**)

**核心特性**:
- RBAC 权限模型 (用户-角色-菜单权限)
- API 权限模式匹配 (避免为 500+ 个 API 单独配置)
- 三层权限缓存 (本地内存 → Redis → MySQL)

**技术栈**:
- 后端: Go 1.24 + Gin + GORM v1.31.0 + Redis
- 前端: React 19 + TypeScript + Vite 7
- 数据库: MySQL 8.0+

**⚠️ 重要**: 目前项目处于 **Week 1 状态**,仅完成数据库设计和Model层,**尚无可用的API和前端页面**

---

## 项目结构 (DDD三层架构)

```
backend/
├── cmd/server/main.go           # 应用入口
├── internal/
│   ├── admin/                   # 核心业务模块
│   │   ├── biz/                 # 业务逻辑层 ✅ 已实现
│   │   │   ├── user/            # 用户业务逻辑
│   │   │   ├── role/            # 角色业务逻辑
│   │   │   └── biz.go           # IBiz统一接口
│   │   ├── store/               # 数据访问层 ✅ 已实现
│   │   │   ├── user.go          # 用户CRUD
│   │   │   ├── role.go          # 角色CRUD
│   │   │   └── store.go         # IStore统一接口
│   │   ├── controller/v1/       # HTTP处理层 ❌ 待实现
│   │   └── middleware/          # 业务中间件 ❌ 待实现
│   └── pkg/                     # 内部共享包
│       ├── model/               # GORM模型 ✅ 已实现
│       ├── cache/               # 缓存实现 🟡 部分实现
│       └── config/              # 配置管理 ✅ 已实现
├── Makefile                     # 构建脚本 ✅ 已实现
└── ARCHITECTURE.md              # 详细架构文档 (必读!)

react-admin/
├── src/
│   ├── pages/                   # 页面组件 ❌ 待实现
│   ├── components/              # 通用组件 ❌ 待实现
│   └── App.tsx                  # 应用入口 🟡 仅脚手架
└── package.json                 # 依赖缺失: react-router/antd/axios
```

**架构依赖流向** (严格单向):
```
Controller → Biz → Store → Model → Database
(HTTP处理) → (业务逻辑) → (数据访问) → (GORM模型) → (MySQL)
```

**⚠️ 当前问题**: Controller层为空,导致无法访问任何API

---

## 常用命令

### 后端 (推荐使用 Makefile)
```bash
cd backend

make help           # 查看所有命令
make dev            # 开发模式(热重载,需先 make install-tools)
make run            # 直接运行
make test           # 运行测试
make build          # 编译到 bin/server
make check          # 代码质量检查 (fmt+vet+test)

# 特定测试
go test -v ./internal/admin/biz/user/
go test -v -run TestUserBiz/Create ./internal/admin/biz/user/
```

### 前端
```bash
cd react-admin

npm install         # 安装依赖
npm run dev         # 启动开发服务器 (http://localhost:5173)
npm run build       # 构建生产版本
npm run lint        # ESLint检查
```

### 数据库
```bash
# 连接MySQL
mysql -u root -p

# 导入表结构 (需先创建 docs/schema.sql)
cd backend
make migrate
# 或: mysql -u root -p go_react_admin < ../docs/schema.sql
```

---

## 架构核心原则 (DDD三层架构)

### 依赖规则 (严格单向,禁止跨层)
- ✅ Controller 只能调用 Biz 接口
- ✅ Biz 只能调用 Store 接口
- ✅ Store 只能操作 Model 和 `*gorm.DB`
- ❌ Controller 不能直接访问 Store
- ❌ Controller 不能直接使用 `*gorm.DB`

### 各层职责

**Controller层** (`internal/admin/controller/v1/`):
- 解析 HTTP 请求参数 (ShouldBindJSON)
- 参数验证 (binding tags)
- 调用 Biz 层方法
- 封装 HTTP 响应 (core.WriteResponse)
- ❌ 禁止编写业务逻辑,禁止访问数据库

**Biz层** (`internal/admin/biz/`):
- 核心业务逻辑 (密码加密/验证/数据转换)
- 跨Store协调和事务管理
- 缓存管理 (三层缓存读写)
- 定义 Request/Response 结构体
- ❌ 禁止处理 HTTP 细节,禁止直接使用 `*gorm.DB`

**Store层** (`internal/admin/store/`):
- 数据库 CRUD 操作
- GORM 查询构建
- 数据权限 SQL 拼接
- 返回 `model.*` GORM 模型
- ❌ 禁止编写业务逻辑 (如密码加密)

### 接口设计原则 (SOLID)
```go
// ✅ 正确: 定义领域接口,隐藏实现细节
type IBiz interface {
    Users() user.IUserBiz
    Roles() role.IRoleBiz
}

type IUserBiz interface {
    Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
    GetByID(ctx context.Context, id uint64) (*UserResponse, error)
}

// ❌ 错误: 暴露实现细节
type IUserBiz interface {
    CreateUser(db *gorm.DB, user *model.User) error  // 暴露GORM
}
```

### 错误处理
- Biz层: 返回业务错误 (`fmt.Errorf` 或自定义错误)
- Store层: 返回 `gorm.Error`
- Controller层: 使用 `core.WriteResponse` 统一封装

---

## 开发工作流 (DDD标准流程)

开发新功能时,按此顺序:

### 1. 定义 Model (数据模型)
位置: `internal/pkg/model/xxx.go`
- 定义 GORM 结构体
- 添加索引标签 (gorm:"index")
- 定义关联关系 (has-many/belongs-to)

### 2. 实现 Store 层 (数据访问)
步骤:
1. 在 `internal/admin/store/store.go` 添加接口定义
2. 在 `internal/admin/store/xxx.go` 实现 CRUD 方法
3. 在 `internal/admin/store/store_factory.go` 注册到 IStore

### 3. 实现 Biz 层 (业务逻辑)
步骤:
1. 创建 `internal/admin/biz/xxx/` 目录
2. 定义 IXxxBiz 接口 + Request/Response 结构体
3. 实现业务逻辑 (调用 Store 层)
4. 在 `internal/admin/biz/biz_factory.go` 注册到 IBiz

### 4. 实现 Controller 层 (HTTP 处理)
步骤:
1. 创建 `internal/admin/controller/v1/xxx.go`
2. 定义 XxxController 结构体 (注入 IBiz)
3. 实现 HTTP 处理方法 (调用 Biz 层)

### 5. 注册路由
在 `internal/admin/router.go` 注册路由 (待实现)
或在 `cmd/server/main.go` 临时注册

### 6. 测试
```bash
make test                                 # 全部测试
go test -v ./internal/admin/biz/xxx/      # Biz层测试
go test -v ./internal/admin/store/        # Store层测试
```

---

## 核心参考文档

- **架构文档**: `backend/ARCHITECTURE.md` (DDD分层架构详细说明)
- **需求文档**: `docs/requirement.md` (完整需求和设计)
- **代码示例**:
  - Biz层: `backend/internal/admin/biz/user/user.go`
  - Store层: `backend/internal/admin/store/user.go`
  - Model: `backend/internal/pkg/model/user.go`

---

## 快速开始

### 环境检查
- [ ] MySQL 8.0+ 运行中
- [ ] Redis 运行中 (单机部署即可)
- [ ] Go 1.24+ 已安装
- [ ] Node.js 18+ 已安装

### 启动步骤
```bash
# 1. 启动后端 (仅能连数据库,无API)
cd backend
cp configs/config.example.yml configs/config.yml  # 修改数据库配置
make migrate                                       # 导入表结构
make run                                           # 启动服务

# 2. 启动前端 (仅脚手架)
cd react-admin
npm install
npm run dev                                        # http://localhost:5173
```

**⚠️ 注意**: 后端启动后无可访问的 HTTP API,前端无可用页面

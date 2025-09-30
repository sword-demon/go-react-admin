# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Go-React-Admin 是一个基于 Go + React 的企业级权限管理系统,支持 500+ API 接口规模。

**核心特性:**
- RBAC权限模型 + 数据权限过滤(根据部门/角色)
- API权限模式匹配(通配符支持,避免为每个API单独配置)
- 三层权限缓存架构(本地内存 + Redis + MySQL)
- 单人开发优化(8周MVP计划)

**技术栈:**
- 后端: Go 1.24 + Gin + GORM v2 + Redis
- 前端: React 19 + TypeScript + Vite + Ant Design 5
- 数据库: MySQL 8.0+

---

## 项目结构

```
go-react-admin/
├── backend/              # Go后端
│   ├── controller/       # 控制器层(HTTP处理)
│   ├── service/          # 业务逻辑层
│   ├── model/            # 数据模型(GORM)
│   ├── middleware/       # 中间件(JWT/权限验证)
│   └── go.mod
├── react-admin/          # React前端
│   ├── src/
│   │   ├── pages/        # 页面组件
│   │   ├── components/   # 通用组件
│   │   └── store/        # Redux状态管理
│   └── package.json
└── docs/
    └── requirement.md    # 需求文档(1200+行,必读!)
```

---

## 常用命令

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 运行开发服务器(热重载需要 air)
go install github.com/air-verse/air@latest
air

# 或直接运行
go run main.go

# 运行测试
go test ./...

# 运行特定测试
go test -v -run TestUserService ./service

# 构建
go build -o bin/server main.go

# 格式化代码
go fmt ./...

# 代码检查
go vet ./...
```

### 前端开发

```bash
# 进入前端目录
cd react-admin

# 安装依赖
npm install

# 运行开发服务器(http://localhost:5173)
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# ESLint检查
npm run lint
```

### 数据库

```bash
# MySQL连接
mysql -u root -p

# 导入表结构(需先创建SQL文件)
mysql -u root -p go_react_admin < docs/schema.sql
```

---

## 架构关键点

### 1. 权限验证流程(核心!)

```
请求 → JWT中间件 → 权限中间件 → 业务逻辑 → 数据权限过滤 → 响应
          ↓             ↓                         ↓
      解析Token    三层缓存查权限             SQL自动拼接WHERE
                   (本地→Redis→MySQL)
```

**三层缓存架构**:
- Layer 1: 本地内存(5分钟TTL, 80%命中率, <1ms)
- Layer 2: Redis(30分钟TTL, 95%命中率, <10ms)
- Layer 3: MySQL(持久化, 100%命中率, 10-50ms)

**关键代码位置**:
- 权限中间件: `backend/middleware/permission.go`
- 缓存实现: `backend/service/cache.go`
- 模式匹配: `backend/service/permission.go`

### 2. API权限模式匹配(避免500个API单独配置)

**权限模式示例**:
```
*:*              # 全部API
user:*           # 用户模块所有操作
user:read        # 用户模块只读(GET)
user:write       # 用户模块写入(POST/PUT/DELETE)
/api/admin/*     # 路径通配符
/api/users:GET   # 路径+方法精确匹配
```

**优先级** (高→低):
1. 精确匹配: `/api/users:GET`
2. 路径通配: `/api/users/*`
3. 模块权限: `user:read`, `user:write`
4. 模块通配: `user:*`
5. 全局通配: `*:*`

**数据库设计**:
- 不要用 `sys_role_api` 关联表(会有5000+条记录)!
- 使用 `sys_role_permission` 模式表(只需50-100条)

### 3. 数据权限过滤(根据部门/角色)

**4种数据范围**:
- `ALL`: 全部数据(超管)
- `DEPT_AND_CHILD`: 本部门+子部门(部门经理)
- `DEPT_ONLY`: 仅本部门(部门主管)
- `SELF_ONLY`: 仅本人数据(普通员工)

**实现方式**:
- 在权限中间件中注入 `dataScope` 到Context
- DAO层查询时自动拼接WHERE条件
- 示例: `WHERE dept_id IN (SELECT id FROM sys_dept WHERE FIND_IN_SET(?, ancestors))`

### 4. 核心数据表(7张)

1. **sys_user** - 用户表(索引: username, dept_id, status)
2. **sys_role** - 角色表(索引: role_key)
3. **sys_dept** - 部门表(树形结构,索引: parent_id, ancestors)
4. **sys_menu** - 菜单表(树形结构,索引: parent_id)
5. **sys_role_permission** - 权限模式表(核心!存储 user:*, /api/admin/* 等模式)
6. **sys_user_role** - 用户-角色关联表(复合主键)
7. **sys_role_menu** - 角色-菜单关联表(复合主键)

**注意**: `sys_api_doc` 仅用于API文档展示,不用于权限验证!

---

## 开发约定

### Go后端代码规范

1. **分层架构** (严格遵守):
   ```
   controller → service → model
   (HTTP处理) → (业务逻辑) → (数据访问)
   ```

2. **错误处理**:
   - 使用统一的错误响应格式
   - Service层返回 `error`
   - Controller层处理HTTP状态码

3. **数据库查询**:
   - 所有查询必须考虑数据权限
   - 高频查询必须加索引
   - 使用GORM的Preload避免N+1查询

4. **缓存策略**:
   - 权限数据必须走三层缓存
   - 权限变更时必须清除相关用户缓存
   - 角色变更时批量清除该角色下所有用户缓存

### React前端代码规范

1. **状态管理**(后续添加Redux Toolkit后):
   - 全局状态: Redux(用户信息、权限列表、菜单数据)
   - 本地状态: useState(表单、临时UI状态)

2. **权限控制**:
   - 使用 `usePermission` Hook检查权限
   - 按钮级权限: `{hasPermission('user:create') && <Button />}`
   - 路由级权限: 动态生成路由

3. **API请求**:
   - 统一使用 axios 实例
   - 请求拦截器注入JWT Token
   - 响应拦截器处理401/403

---

## 单人开发注意事项

**这是单人全栈项目,8周MVP计划!**

### 优先级原则
1. **Phase 1 (MVP, 8周)**: 权限核心 + 登录日志 + 基础Excel导入
2. **Phase 2 (后续)**: 审计日志 + 操作日志 + 高级导入导出

### 开发策略
- **前后端切换**: 按模块完整开发(先写完后端API,再写前端UI),减少切换成本
- **代码复用**: 使用 GORM Gen 生成CRUD, Ant Design Pro组件
- **测试策略**: 优先集成测试,单元测试仅覆盖核心逻辑
- **技术债务**: 允许存在,但必须用 `// TODO:` 标记

### 进度监控
- **每周五**: 自查进度,延期2天立即调整范围
- **如果Week 4还没完成用户管理**: 立即砍掉API权限管理
- **如果卡住超过2小时**: 去Stack Overflow求助
- **如果卡住超过1天**: 调整方案或砍功能

### 风险应对
- **Plan B (7周)**: 砍掉数据导入导出、登录日志
- **Plan C (6周)**: 砍掉API权限管理、数据权限过滤

---

## 核心参考文档

- **需求文档**: `docs/requirement.md` (1200+行,包含完整架构设计和代码示例)
- **若依参考**: https://gitee.com/y_project/RuoYi (权限设计参考)
- **Gin文档**: https://gin-gonic.com/
- **GORM文档**: https://gorm.io/
- **Ant Design**: https://ant.design/

---

## 快速开始检查清单

开始开发前确认:

- [ ] MySQL 8.0+ 已安装并运行
- [ ] Redis 已安装并运行(单机部署即可)
- [ ] Go 1.23+ 已安装
- [ ] Node.js 18+ 已安装
- [ ] 已读完 `docs/requirement.md` (特别是架构部分)
- [ ] 已创建数据库表结构(7张核心表)
- [ ] 后端能跑通 Hello World
- [ ] 前端能跑通 Vite 开发服务器

---

**重要**: 这是一个雄心勃勃的单人项目!优先完成MVP,不要追求完美!能上线比完美更重要!
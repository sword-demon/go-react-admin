# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 🚀 项目当前状态(2025-10-01)

**架构重构**: ✅ 已完成DDD三层架构迁移(Controller → Biz → Store)

**后端已实现**:
- ✅ Biz层架构(user/role/dept/menu/permission模块,~26个Go文件)
- ✅ Store层接口(IStore统一接口,包含Users/Roles/Depts/Menus/Permissions)
- ✅ Model层定义(User/Role/Dept/Menu/Permission/Log等)
- ✅ 配置管理(Viper + YAML)
- ✅ 数据库连接(GORM v1.31.0 + MySQL 8.0+)
- ✅ Redis缓存客户端
- ✅ Makefile构建脚本(make dev/run/test/build)

**后端待实现**(TODO):
- ⏳ Controller层(v1目录已创建,但无实现文件)
- ⏳ 路由注册(router.go待实现,main.go有临时路由)
- ⏳ JWT中间件(internal/admin/middleware/)
- ⏳ 权限中间件(三层缓存+模式匹配)
- ⏳ 数据权限过滤中间件

**前端状态**:
- ✅ Vite + React 19 + TypeScript 基础脚手架
- ⏳ Redux状态管理(待实现)
- ⏳ 页面组件(待实现)

**下一步重点**:
1. 实现Controller层(用户/角色/部门CRUD)
2. 实现JWT认证中间件
3. 实现权限验证中间件(三层缓存)
4. 注册路由到main.go或独立router.go

---

## 项目概述

Go-React-Admin 是一个基于 Go + React 的企业级权限管理系统,支持 500+ API 接口规模。

**核心特性:**

- RBAC 权限模型 + 数据权限过滤(根据部门/角色)
- API 权限模式匹配(通配符支持,避免为每个 API 单独配置)
- 三层权限缓存架构(本地内存 + Redis + MySQL)
- 单人开发优化(8 周 MVP 计划)

**技术栈:**

- 后端: Go 1.24 + Gin + GORM v2 + Redis
- 前端: React 19 + TypeScript + Vite + Ant Design 5
- 数据库: MySQL 8.0+

---

## 项目结构(DDD架构 - 2025年重构)

```
go-react-admin/
├── backend/                      # Go后端(遵循golang-standards/project-layout)
│   ├── cmd/
│   │   └── server/main.go        # 应用入口
│   ├── internal/                 # 内部代码(不可外部import)
│   │   ├── admin/                # 核心业务模块
│   │   │   ├── biz/              # 业务逻辑层(替代service)
│   │   │   │   ├── user/         # 用户业务逻辑
│   │   │   │   ├── role/         # 角色业务逻辑
│   │   │   │   ├── dept/         # 部门业务逻辑
│   │   │   │   ├── menu/         # 菜单业务逻辑
│   │   │   │   ├── permission/   # 权限业务逻辑
│   │   │   │   └── biz.go        # IBiz统一接口
│   │   │   ├── store/            # 数据访问层(DAO)
│   │   │   │   ├── user.go       # 用户CRUD
│   │   │   │   ├── role.go       # 角色CRUD
│   │   │   │   ├── dept.go       # 部门CRUD
│   │   │   │   ├── menu.go       # 菜单CRUD
│   │   │   │   ├── permission.go # 权限CRUD
│   │   │   │   └── store.go      # IStore统一接口
│   │   │   ├── controller/v1/    # HTTP处理层
│   │   │   └── middleware/       # 业务中间件
│   │   └── pkg/                  # 内部共享包
│   │       ├── model/            # GORM模型
│   │       ├── cache/            # 三层缓存实现
│   │       ├── config/           # 配置管理
│   │       └── auth/             # JWT工具
│   ├── pkg/                      # 可外部import的公共包
│   │   ├── core/                 # 核心响应封装
│   │   ├── token/                # Token工具
│   │   └── validator/            # 验证器
│   ├── configs/                  # 配置文件
│   │   ├── config.yml            # 实际配置(DO NOT commit!)
│   │   └── config.example.yml    # 配置模板
│   ├── Makefile                  # 构建脚本(make dev/run/test/build)
│   ├── ARCHITECTURE.md           # 架构文档(必读!)
│   └── go.mod
├── react-admin/                  # React前端
│   ├── src/
│   │   ├── pages/                # 页面组件
│   │   ├── components/           # 通用组件
│   │   └── store/                # Redux状态管理(待实现)
│   └── package.json
└── docs/
    └── requirement.md            # 需求文档(1200+行,必读!)
```

**关键变化(2025重构)**:
- ❌ 废弃: `backend/service/` → ✅ 改用: `backend/internal/admin/biz/`
- ❌ 废弃: `backend/model/` → ✅ 改用: `backend/internal/pkg/model/`
- ❌ 废弃: `backend/controller/` → ✅ 改用: `backend/internal/admin/controller/v1/`
- ✅ 新增: `backend/internal/admin/store/` (专职数据访问)

---

## 常用命令

### 后端开发(推荐使用Makefile)

```bash
# 进入后端目录
cd backend

# 查看所有可用命令
make help

# 开发模式(热重载, 需先安装air)
make dev

# 直接运行(无热重载)
make run

# 安装开发工具(air + staticcheck)
make install-tools

# 安装依赖
make deps
# 或: go mod tidy

# 运行测试
make test
# 或: go test ./...

# 运行特定包的测试
go test -v -run TestUserBiz ./internal/admin/biz/user/

# 运行特定函数的测试
go test -v -run TestUserBiz/Create ./internal/admin/biz/user/

# 代码质量检查(fmt+vet+test)
make check

# 构建二进制文件
make build
# 输出: ./bin/server

# 格式化代码
make fmt
# 或: go fmt ./...

# 代码检查
make vet
# 或: go vet ./...

# 静态分析
make lint
# 或: staticcheck ./...

# 清理编译文件
make clean

# 查看项目信息(依赖数/代码行数)
make info
```

**核心命令速记**:
- `make dev` → 启动开发服务器(推荐)
- `make test` → 运行测试
- `make check` → 代码质量检查
- `make build` → 编译生产版本

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

## 架构关键点(DDD三层架构)

### 1. 依赖流向(严格单向)

```
cmd/server/main.go
       ↓
[初始化] → Store → Biz → Controller → Router
       ↓      ↓      ↓        ↓          ↓
   Config  GORM  Business  HTTP     Routes
              ↓      Logic  Handlers
           MySQL
```

**依赖注入示例**(main.go:67-75):
```go
// 5. 初始化store层(数据访问)
dataStore := store.NewStore(database)

// 6. 初始化biz层(业务逻辑)
bizLayer := biz.NewBiz(dataStore, redisClient)

// 7. 初始化controller层(HTTP处理,传入bizLayer)
// TODO: userController := controller.NewUserController(bizLayer)
```

### 2. 三层架构职责划分

#### Controller层(HTTP处理)
- **位置**: `internal/admin/controller/v1/`
- **职责**:
  - 解析HTTP请求参数
  - 参数验证(binding)
  - 调用Biz层方法
  - 封装HTTP响应
- **禁止**:
  - ❌ 直接访问Store层
  - ❌ 编写业务逻辑
  - ❌ 直接操作数据库
- **示例**:
```go
// controller/v1/user.go
type UserController struct {
    biz biz.IBiz  // 仅依赖IBiz接口
}

func (c *UserController) Create(ctx *gin.Context) {
    var req user.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        core.WriteResponse(ctx, errors.ErrBind, nil)
        return
    }

    // 调用Biz层,不关心内部实现
    result, err := c.biz.Users().Create(ctx, &req)
    core.WriteResponse(ctx, err, result)
}
```

#### Biz层(业务逻辑)
- **位置**: `internal/admin/biz/`
- **职责**:
  - 核心业务逻辑
  - 数据验证(业务规则)
  - 多Store协调(事务)
  - 权限检查(数据范围)
  - 缓存管理
- **可访问**: Store层接口(IStore)
- **禁止**:
  - ❌ 处理HTTP细节
  - ❌ 直接使用*gorm.DB
- **关键特性**:
  - ✅ 每个模块独立目录(user/, role/, dept/)
  - ✅ 定义自己的Request/Response结构体
  - ✅ 实现IBiz接口的子接口(IUserBiz等)
- **示例**(biz/user/user.go:85-118):
```go
func (b *userBiz) Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
    // 1. 业务规则验证
    _, err := b.store.Users().GetByUsername(ctx, req.Username)
    if err == nil {
        return nil, fmt.Errorf("username already exists")
    }

    // 2. 密码加密(业务逻辑)
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    // 3. 构建领域模型
    user := &model.User{
        Username: req.Username,
        Password: string(hashedPassword),
        // ...
    }

    // 4. 调用Store层持久化
    if err := b.store.Users().Create(ctx, user); err != nil {
        return nil, err
    }

    return b.toUserResponse(user), nil
}
```

#### Store层(数据访问)
- **位置**: `internal/admin/store/`
- **职责**:
  - 数据库CRUD操作
  - GORM查询构建
  - 数据权限SQL拼接
  - 关联查询(Preload)
- **直接使用**: `*gorm.DB`
- **禁止**:
  - ❌ 编写业务逻辑
  - ❌ 密码加密等非数据操作
- **关键特性**:
  - ✅ 实现IStore接口的子接口
  - ✅ 所有方法接收context.Context
  - ✅ 返回model.*模型(GORM结构体)
- **示例**(store/user.go片段):
```go
func (s *userStore) Create(ctx context.Context, user *model.User) error {
    return s.db.WithContext(ctx).Create(user).Error
}

func (s *userStore) GetByUsername(ctx context.Context, username string) (*model.User, error) {
    var user model.User
    err := s.db.WithContext(ctx).
        Where("username = ?", username).
        First(&user).Error
    return &user, err
}
```

### 3. 权限验证流程(核心!)

```
请求 → JWT中间件 → 权限中间件 → 业务逻辑 → 数据权限过滤 → 响应
          ↓             ↓                         ↓
      解析Token    三层缓存查权限             SQL自动拼接WHERE
                   (本地→Redis→MySQL)
```

**三层缓存架构**:

- Layer 1: 本地内存(5 分钟 TTL, 80%命中率, <1ms)
- Layer 2: Redis(30 分钟 TTL, 95%命中率, <10ms)
- Layer 3: MySQL(持久化, 100%命中率, 10-50ms)

**关键代码位置**(已更新):
- 权限中间件: `internal/admin/middleware/permission.go`
- 缓存实现: `internal/pkg/cache/` (三层缓存)
- 权限Biz: `internal/admin/biz/permission/permission.go`
- 权限Store: `internal/admin/store/permission.go`

### 4. API 权限模式匹配(避免 500 个 API 单独配置)

**权限模式示例**:

```
*:*              # 全部API
user:*           # 用户模块所有操作
user:read        # 用户模块只读(GET)
user:write       # 用户模块写入(POST/PUT/DELETE)
/api/admin/*     # 路径通配符
/api/users:GET   # 路径+方法精确匹配
```

**优先级** (高 → 低):

1. 精确匹配: `/api/users:GET`
2. 路径通配: `/api/users/*`
3. 模块权限: `user:read`, `user:write`
4. 模块通配: `user:*`
5. 全局通配: `*:*`

**数据库设计**:

- 不要用 `sys_role_api` 关联表(会有 5000+条记录)!
- 使用 `sys_role_permission` 模式表(只需 50-100 条)

### 5. 数据权限过滤(根据部门/角色)

**4 种数据范围**:

- `ALL`: 全部数据(超管)
- `DEPT_AND_CHILD`: 本部门+子部门(部门经理)
- `DEPT_ONLY`: 仅本部门(部门主管)
- `SELF_ONLY`: 仅本人数据(普通员工)

**实现方式**(DDD架构):
- 在权限中间件中注入 `dataScope` 到 Context
- Store层查询时从Context读取并自动拼接 WHERE 条件
- 示例: `WHERE dept_id IN (SELECT id FROM sys_dept WHERE FIND_IN_SET(?, ancestors))`

### 6. 核心数据表(7 张)

1. **sys_user** - 用户表(索引: username, dept_id, status)
2. **sys_role** - 角色表(索引: role_key)
3. **sys_dept** - 部门表(树形结构,索引: parent_id, ancestors)
4. **sys_menu** - 菜单表(树形结构,索引: parent_id)
5. **sys_role_permission** - 权限模式表(核心!存储 user:_, /api/admin/_ 等模式)
6. **sys_user_role** - 用户-角色关联表(复合主键)
7. **sys_role_menu** - 角色-菜单关联表(复合主键)

**注意**: `sys_api_doc` 仅用于 API 文档展示,不用于权限验证!

---

## 开发约定(DDD架构规范)

### Go 后端代码规范

1. **分层架构** (严格遵守,不可跨层):

   ```
   Controller → Biz → Store → Model
   (HTTP处理) → (业务逻辑) → (数据访问) → (GORM模型)
   ```

   **依赖规则**:
   - ✅ Controller 只能调用 Biz 接口
   - ✅ Biz 只能调用 Store 接口
   - ✅ Store 只能操作 Model 和 *gorm.DB
   - ❌ Controller 不能直接访问 Store
   - ❌ Controller 不能直接使用 *gorm.DB

2. **接口设计原则** (遵循SOLID):

   ```go
   // ✅ 正确: 定义领域接口
   // internal/admin/biz/biz.go
   type IBiz interface {
       Users() user.IUserBiz
       Roles() role.IRoleBiz
   }

   // internal/admin/biz/user/user.go
   type IUserBiz interface {
       Create(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
       Update(ctx context.Context, id uint64, req *UpdateUserRequest) error
       // ...
   }

   // ❌ 错误: 不要在接口中暴露实现细节
   type IUserBiz interface {
       CreateUser(db *gorm.DB, user *model.User) error  // ❌ 暴露GORM
   }
   ```

3. **错误处理**:

   - Biz层: 返回业务错误(使用 `fmt.Errorf` 或自定义错误类型)
   - Store层: 返回数据库错误(直接返回 `gorm.Error`)
   - Controller层: 使用统一响应格式封装错误
   ```go
   // ✅ Biz层示例
   if err == nil {
       return nil, fmt.Errorf("username already exists")
   }

   // ✅ Controller层示例
   if err != nil {
       core.WriteResponse(ctx, err, nil)  // 统一错误响应
       return
   }
   ```

4. **数据库查询**:

   - Store层所有方法必须接收 `context.Context`
   - 使用 `db.WithContext(ctx)` 传递上下文
   - 高频查询必须加索引
   - 使用 GORM 的 Preload 避免 N+1 查询
   ```go
   // ✅ 正确示例
   func (s *userStore) Get(ctx context.Context, id uint64) (*model.User, error) {
       var user model.User
       err := s.db.WithContext(ctx).First(&user, id).Error
       return &user, err
   }
   ```

5. **缓存策略**:
   - 权限数据必须走三层缓存(本地内存 → Redis → MySQL)
   - 权限变更时必须清除相关用户缓存
   - 角色变更时批量清除该角色下所有用户缓存
   - 缓存逻辑在Biz层处理,不要在Store层缓存

6. **Request/Response 设计**:
   - 每个Biz模块定义自己的Request/Response结构体
   - Request用于参数验证(binding tags)
   - Response用于隐藏敏感字段(如password)
   ```go
   // ✅ biz/user/user.go
   type CreateUserRequest struct {
       Username string `json:"username" binding:"required"`
       Password string `json:"password" binding:"required"`
   }

   type UserResponse struct {
       ID       uint64 `json:"id"`
       Username string `json:"username"`
       // 不包含Password字段
   }
   ```

### React 前端代码规范

1. **状态管理**(后续添加 Redux Toolkit 后):

   - 全局状态: Redux(用户信息、权限列表、菜单数据)
   - 本地状态: useState(表单、临时 UI 状态)

2. **权限控制**:

   - 使用 `usePermission` Hook 检查权限
   - 按钮级权限: `{hasPermission('user:create') && <Button />}`
   - 路由级权限: 动态生成路由

3. **API 请求**:
   - 统一使用 axios 实例
   - 请求拦截器注入 JWT Token
   - 响应拦截器处理 401/403

---

## 单人开发注意事项

**这是单人全栈项目,8 周 MVP 计划!**

### 优先级原则

1. **Phase 1 (MVP, 8 周)**: 权限核心 + 登录日志 + 基础 Excel 导入
2. **Phase 2 (后续)**: 审计日志 + 操作日志 + 高级导入导出

### 开发策略

- **前后端切换**: 按模块完整开发(先写完后端 API,再写前端 UI),减少切换成本
- **代码复用**: 使用 GORM Gen 生成 CRUD, Ant Design Pro 组件
- **测试策略**: 优先集成测试,单元测试仅覆盖核心逻辑
- **技术债务**: 允许存在,但必须用 `// TODO:` 标记

### 进度监控

- **每周五**: 自查进度,延期 2 天立即调整范围
- **如果 Week 4 还没完成用户管理**: 立即砍掉 API 权限管理
- **如果卡住超过 2 小时**: 去 Stack Overflow 求助
- **如果卡住超过 1 天**: 调整方案或砍功能

### 风险应对

- **Plan B (7 周)**: 砍掉数据导入导出、登录日志
- **Plan C (6 周)**: 砍掉 API 权限管理、数据权限过滤

---

## 核心参考文档

- **架构文档**: `backend/ARCHITECTURE.md` (DDD分层架构说明,必读!)
- **需求文档**: `docs/requirement.md` (1200+行,包含完整架构设计和代码示例)
- **已实现代码**:
  - Biz层示例: `backend/internal/admin/biz/user/user.go`
  - Store层示例: `backend/internal/admin/store/user.go`
  - Model定义: `backend/internal/pkg/model/user.go`
- **若依参考**: https://gitee.com/y_project/RuoYi (权限设计参考)
- **Gin 文档**: https://gin-gonic.com/
- **GORM 文档**: https://gorm.io/ (v1.31.0+)
- **Ant Design**: https://ant.design/

---

## 开发工作流(DDD标准流程)

开发新功能时,必须按以下顺序:

### 1. 定义Model(数据模型)
```bash
# 文件位置: internal/pkg/model/xxx.go
# 定义GORM结构体,添加索引标签
```

### 2. 实现Store层(数据访问)
```bash
# 步骤:
# 1) 在 internal/admin/store/store.go 添加接口定义
# 2) 在 internal/admin/store/xxx.go 实现CRUD方法
# 3) 在 internal/admin/store/store_factory.go 注册到IStore
```

### 3. 实现Biz层(业务逻辑)
```bash
# 步骤:
# 1) 创建 internal/admin/biz/xxx/ 目录
# 2) 定义 IXxxBiz 接口 + Request/Response结构体
# 3) 实现业务逻辑(调用Store层)
# 4) 在 internal/admin/biz/biz_factory.go 注册到IBiz
```

### 4. 实现Controller层(HTTP处理)
```bash
# 步骤:
# 1) 创建 internal/admin/controller/v1/xxx.go
# 2) 定义 XxxController 结构体(注入IBiz)
# 3) 实现HTTP处理方法(调用Biz层)
```

### 5. 注册路由
```bash
# 在 internal/admin/router.go 注册路由(待实现)
# 或在 cmd/server/main.go 临时注册
```

### 6. 测试
```bash
make test                      # 运行所有测试
go test -v ./internal/admin/biz/xxx/  # 测试Biz层
go test -v ./internal/admin/store/    # 测试Store层
```

---

## 快速开始检查清单

开始开发前确认:

- [ ] MySQL 8.0+ 已安装并运行
- [ ] Redis 已安装并运行(单机部署即可)
- [ ] Go 1.24+ 已安装
- [ ] Node.js 18+ 已安装
- [ ] 已读完 `backend/ARCHITECTURE.md` (理解DDD分层架构)
- [ ] 已读完 `docs/requirement.md` (特别是权限设计部分)
- [ ] 已创建数据库表结构(7 张核心表,运行 `make migrate`)
- [ ] 后端能跑通 (`make run` 访问 http://localhost:8080/ping)
- [ ] 前端能跑通 Vite 开发服务器 (`cd react-admin && npm run dev`)

---

**重要**: 这是一个雄心勃勃的单人项目!优先完成 MVP,不要追求完美!能上线比完美更重要!

---

## DDD架构迁移说明(2025年重构)

**如果你看到旧代码引用以下路径,请更新**:

| 旧路径(已废弃)            | 新路径(DDD架构)                        |
|---------------------------|----------------------------------------|
| `backend/service/`        | `backend/internal/admin/biz/`          |
| `backend/model/`          | `backend/internal/pkg/model/`          |
| `backend/controller/`     | `backend/internal/admin/controller/v1/`|
| `backend/middleware/`     | `backend/internal/admin/middleware/`   |
| `go run main.go`          | `cd backend && make run`               |

**架构升级收益**:
- ✅ 清晰的依赖方向(Controller → Biz → Store → Model)
- ✅ 更好的可测试性(接口隔离)
- ✅ 符合golang-standards/project-layout标准
- ✅ internal/包保护(防止外部错误引用)
- ✅ 更易扩展(添加新模块只需复制user/模板)

## 开发规则

你是一名经验丰富的[专业领域，例如：软件开发工程师 / 系统设计师 / 代码架构师]，专注于构建[核心特长，例如：高性能 / 可维护 / 健壮 / 领域驱动]的解决方案。

你的任务是：**审查、理解并迭代式地改进/推进一个[项目类型，例如：现有代码库 / 软件项目 / 技术流程]。**

在整个工作流程中，你必须内化并严格遵循以下核心编程原则，确保你的每次输出和建议都体现这些理念：

- **简单至上 (KISS):** 追求代码和设计的极致简洁与直观，避免不必要的复杂性。
- **精益求精 (YAGNI):** 仅实现当前明确所需的功能，抵制过度设计和不必要的未来特性预留。
- **坚实基础 (SOLID):**
  - **S (单一职责):** 各组件、类、函数只承担一项明确职责。
  - **O (开放/封闭):** 功能扩展无需修改现有代码。
  - **L (里氏替换):** 子类型可无缝替换其基类型。
  - **I (接口隔离):** 接口应专一，避免“胖接口”。
  - **D (依赖倒置):** 依赖抽象而非具体实现。
- **杜绝重复 (DRY):** 识别并消除代码或逻辑中的重复模式，提升复用性。

**请严格遵循以下工作流程和输出要求：**

1.  **深入理解与初步分析（理解阶段）：**

    - 详细审阅提供的[资料/代码/项目描述]，全面掌握其当前架构、核心组件、业务逻辑及痛点。
    - 在理解的基础上，初步识别项目中潜在的**KISS, YAGNI, DRY, SOLID**原则应用点或违背现象。

2.  **明确目标与迭代规划（规划阶段）：**

    - 基于用户需求和对现有项目的理解，清晰定义本次迭代的具体任务范围和可衡量的预期成果。
    - 在规划解决方案时，优先考虑如何通过应用上述原则，实现更简洁、高效和可扩展的改进，而非盲目增加功能。

3.  **分步实施与具体改进（执行阶段）：**

    - 详细说明你的改进方案，并将其拆解为逻辑清晰、可操作的步骤。
    - 针对每个步骤，具体阐述你将如何操作，以及这些操作如何体现**KISS, YAGNI, DRY, SOLID**原则。例如：
      - “将此模块拆分为更小的服务，以遵循 SRP 和 OCP。”
      - “为避免 DRY，将重复的 XXX 逻辑抽象为通用函数。”
      - “简化了 Y 功能的用户流，体现 KISS 原则。”
      - “移除了 Z 冗余设计，遵循 YAGNI 原则。”
    - 重点关注[项目类型，例如：代码质量优化 / 架构重构 / 功能增强 / 用户体验提升 / 性能调优 / 可维护性改善 / Bug 修复]的具体实现细节。

4.  **总结、反思与展望（汇报阶段）：**
    - 提供一个清晰、结构化且包含**实际代码/设计变动建议（如果适用）**的总结报告。
    - 报告中必须包含：
      - **本次迭代已完成的核心任务**及其具体成果。
      - **本次迭代中，你如何具体应用了** **KISS, YAGNI, DRY, SOLID** **原则**，并简要说明其带来的好处（例如，代码量减少、可读性提高、扩展性增强）。
      - **遇到的挑战**以及如何克服。
      - **下一步的明确计划和建议。**

---

# MCP 服务调用规则

## 核心策略

- **审慎单选**：优先离线工具，确需外呼时每轮最多 1 个 MCP 服务
- **序贯调用**：多服务需求时必须串行，明确说明每步理由和产出预期
- **最小范围**：精确限定查询参数，避免过度抓取和噪声
- **可追溯性**：答复末尾统一附加"工具调用简报"

## 服务选择优先级

### 1. Serena（本地代码分析优先）

**工具能力**：find_symbol, find_referencing_symbols, get_symbols_overview, search_for_pattern, read_file, replace_symbol_body, create_text_file, execute_shell_command
**触发场景**：代码检索、架构分析、跨文件引用、项目理解
**调用策略**：

- 先用 get_symbols_overview 快速了解文件结构
- find_symbol 精确定位（支持 name_path 模式匹配）
- search_for_pattern 用于复杂正则搜索
- 限制 relative_path 到相关目录，避免全项目扫描

### 2. Context7（官方文档查询）

**流程**：resolve-library-id → get-library-docs
**触发场景**：框架 API、配置文档、版本差异、迁移指南
**限制参数**：tokens≤5000, topic 指定聚焦范围

### 3. Sequential Thinking（复杂规划）

**触发场景**：多步骤任务分解、架构设计、问题诊断流程
**输出要求**：6-10 步可执行计划，不暴露推理过程
**参数控制**：total_thoughts≤10, 每步一句话描述

### 4. DuckDuckGo（外部信息）

**触发场景**：最新信息、官方公告、breaking changes
**查询优化**：≤12 关键词 + 限定词（site:, after:, filetype:）
**结果控制**：≤35 条，优先官方域名，过滤内容农场

### 5. Playwright（浏览器自动化）

**触发场景**：网页截图、表单测试、SPA 交互验证
**安全限制**：仅开发测试用途

## 错误处理和降级

### 失败策略

- **429 限流**：退避 20s，降低参数范围
- **5xx/超时**：单次重试，退避 2s
- **无结果**：缩小范围或请求澄清

### 降级链路

1. Context7 → DuckDuckGo(site:官方域名)
2. DuckDuckGo → 请求用户提供线索
3. Serena → 使用 Claude Code 本地工具
4. 最终降级 → 保守离线答案 + 标注不确定性

## 实际调用约束

### 禁用场景

- 网络受限且未明确授权
- 查询包含敏感代码/密钥
- 本地工具可充分完成任务

### 并发控制

- **严格串行**：禁止同轮并发调用多个 MCP 服务
- **意图分解**：多服务需求时拆分为多轮对话
- **明确预期**：每次调用前说明预期产出和后续步骤

## 工具调用简报格式

【MCP 调用简报】
服务: <serena|context7|sequential-thinking|ddg-search|playwright>
触发: <具体原因>
参数: <关键参数摘要>
结果: <命中数/主要来源>
状态: <成功|重试|降级>

## 典型调用模式

### 代码分析模式

1. serena.get_symbols_overview → 了解文件结构
2. serena.find_symbol → 定位具体实现
3. serena.find_referencing_symbols → 分析调用关系

### 文档查询模式

1. context7.resolve-library-id → 确定库标识
2. context7.get-library-docs → 获取相关文档段落

### 规划执行模式

1. sequential-thinking → 生成执行计划
2. serena 工具链 → 逐步实施代码修改
3. 验证测试 → 确保修改正确性

### 编码输出/语言偏好###

## Communication & Language

- Default language: Simplified Chinese for issues, PRs, and assistant replies, unless a thread explicitly requests English.
- Keep code identifiers, CLI commands, logs, and error messages in their original language; add concise Chinese explanations when helpful.
- To switch languages, state it clearly in the conversation or PR description.

## File Encoding

When modifying or adding any code files, the following coding requirements must be adhered to:

- Encoding should be unified to UTF-8 (without BOM). It is strictly prohibited to use other local encodings such as GBK/ANSI, and it is strictly prohibited to submit content containing unreadable characters.
- When modifying or adding files, be sure to save them in UTF-8 format; if you find any files that are not in UTF-8 format before submitting, please convert them to UTF-8 before submitting.

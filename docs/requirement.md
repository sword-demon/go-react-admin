# 🎯 后台权限管理系统需求文档

## 📋 项目概述

### 项目名称
通用后台权限管理系统 (Go-React-Admin)

### 技术栈
- **后端**: Go 1.23+ (Gin框架) ✅ 已确认
- **前端**: React 18.x + Ant Design 5.x + Redux Toolkit ✅ 已确认
- **数据库**: MySQL 8.0+ / PostgreSQL 14+
- **ORM**: GORM v2 ✅ 已确认
- **缓存**: Redis 7.0+ (单机部署) ✅ 已确认
- **认证**: JWT Token
- **权限模型**: RBAC + 数据权限过滤 (Data Scope)

### 核心目标
构建一个企业级的后台权限管理系统,支持**500+API接口**规模,实现:
1. ✅ 基于角色的API接口访问控制 (RBAC + 通配符模式)
2. ✅ 基于部门/角色的数据权限过滤
3. ✅ 灵活的菜单和按钮权限管理
4. ✅ 完整的用户、角色、部门、权限管理
5. ✅ 三层权限缓存架构 (本地缓存 + Redis + MySQL)

---

## 🏗️ 系统架构

### 系统规模定位
- **API接口数量**: 500+ (大型系统)
- **预估用户数**: 100-500 (中小型企业) ⚡ 已调整
- **角色数量**: 5-15 (精简配置)
- **并发要求**: 500+ QPS (单机足够)
- **开发团队**: 1人 (全栈独立开发) ⚠️ 风险提示

### 权限模型设计

```
用户 (User)
  ├── 所属部门 (Department)
  ├── 绑定角色 (Roles) - 多对多
  └── 数据权限范围 (Data Scope)

角色 (Role)
  ├── 菜单权限 (Menu Permissions)
  ├── API权限模式 (API Permission Patterns) ⚡ 核心优化
  │   ├── 模块级通配符 (user:*, product:*)
  │   ├── 操作级权限 (user:read, user:write)
  │   └── 路径模式匹配 (/api/admin/*, /api/users:GET)
  └── 数据权限范围 (Data Scope Level)
       ├── 全部数据 (ALL)
       ├── 本部门及子部门 (DEPT_AND_CHILD)
       ├── 仅本部门 (DEPT_ONLY)
       └── 仅本人 (SELF_ONLY)
```

### 数据权限过滤逻辑

| 角色类型 | 数据权限范围 | 示例 |
|---------|------------|------|
| 超级管理员 | 全部数据 | 可查看所有部门的数据 |
| 部门经理 | 本部门+子部门 | 可查看本部门及下属部门数据 |
| 部门主管 | 仅本部门 | 只能查看本部门数据 |
| 普通员工 | 仅本人 | 只能查看/操作自己创建的数据 |

---

## 📦 功能模块

### 1. 用户管理模块

#### 1.1 用户基础信息
- 用户账号、姓名、邮箱、手机号
- 所属部门、直属上级
- 用户状态(启用/禁用)
- 创建时间、最后登录时间

#### 1.2 用户操作
- ✅ 创建用户 (分配部门、角色)
- ✅ 编辑用户信息
- ✅ 重置用户密码
- ✅ 启用/禁用用户
- ✅ 删除用户
- ✅ 批量导入用户 (Excel)

#### 1.3 权限控制
- **API权限**:
  - `POST /api/users` - 创建用户 (需要 `user:create` 权限)
  - `PUT /api/users/:id` - 编辑用户 (需要 `user:update` 权限)
  - `DELETE /api/users/:id` - 删除用户 (需要 `user:delete` 权限)
  - `GET /api/users` - 查询用户列表 (需要 `user:list` 权限)
- **数据权限**: 根据角色数据范围过滤用户列表

---

### 2. 角色管理模块

#### 2.1 角色基础信息
- 角色名称、角色编码 (roleKey)
- 角色描述
- 数据权限范围 (dataScope)
- 显示顺序、状态

#### 2.2 角色操作
- ✅ 创建角色
- ✅ 编辑角色
- ✅ 删除角色 (检查是否有用户绑定)
- ✅ 分配权限 (菜单权限 + API权限)
- ✅ 设置数据权限范围

#### 2.3 权限分配
**菜单权限树:**
```
系统管理
  ├── 用户管理 (menu:user)
  ├── 角色管理 (menu:role)
  ├── 部门管理 (menu:dept)
  └── 菜单管理 (menu:menu)
```

**API权限点:**
```
用户管理:
  - user:list (查询)
  - user:create (新增)
  - user:update (修改)
  - user:delete (删除)
  - user:export (导出)
```

---

### 3. 部门管理模块

#### 3.1 部门结构
- 树形结构 (支持无限层级)
- 部门名称、部门编码
- 负责人、联系电话
- 显示顺序、状态

#### 3.2 部门操作
- ✅ 创建部门 (选择父部门)
- ✅ 编辑部门信息
- ✅ 删除部门 (检查子部门和用户)
- ✅ 部门树形展示

#### 3.3 数据权限关联
- 用户所属部门决定数据权限范围
- 支持跨部门查询 (根据角色配置)

---

### 4. 菜单管理模块

#### 4.1 菜单类型
- **目录 (Directory)**: 一级菜单,无路由
- **菜单 (Menu)**: 带路由的页面
- **按钮 (Button)**: 页面内操作按钮

#### 4.2 菜单属性
- 菜单名称、路由地址
- 组件路径 (前端组件)
- 权限标识 (permission)
- 图标、排序
- 是否显示、是否缓存

#### 4.3 菜单操作
- ✅ 创建菜单 (树形结构)
- ✅ 编辑菜单
- ✅ 删除菜单
- ✅ 菜单树形展示

---

### 5. API权限管理模块 ⚡ (500+接口核心优化)

#### 5.1 API接口自动扫描与分组

**问题**: 手动管理500+个API不现实!

**解决方案**: 路由注册时自动扫描 + 模块分组

```go
// 使用中间件装饰器自动注册API
type APIMetadata struct {
    Module      string  // 模块名 (user/product/order)
    Permission  string  // 权限标识 (user:*)
    Description string
}

// 路由注册示例
func RegisterRoutes(r *gin.Engine) {
    // 用户模块 - 自动归类
    userAPI := r.Group("/api/users",
        WithPermission("user:*", "用户管理模块"))
    {
        userAPI.GET("", GetUsers)           // user:list
        userAPI.POST("", CreateUser)        // user:create
        userAPI.PUT("/:id", UpdateUser)     // user:update
        userAPI.DELETE("/:id", DeleteUser)  // user:delete
        userAPI.GET("/:id", GetUser)        // user:detail
    }

    // 产品模块
    productAPI := r.Group("/api/products",
        WithPermission("product:*", "产品管理模块"))
    {
        productAPI.GET("", GetProducts)
        productAPI.POST("", CreateProduct)
    }

    // 启动时自动扫描所有路由,生成API文档
    apiScanner := NewAPIScanner(r)
    apiScanner.ScanAndSyncToDB() // 同步到数据库
}
```

#### 5.2 API权限模式匹配 (核心!)

**不要给每个API单独配权限!用模式匹配!**

| 权限模式 | 匹配规则 | 示例 | 适用场景 |
|---------|---------|------|---------|
| `*:*` | 全部API | 超级管理员 | 超管 |
| `user:*` | 用户模块所有操作 | `/api/users/*` 所有方法 | 用户管理员 |
| `user:read` | 用户模块只读 | `/api/users (GET)` | 只读角色 |
| `user:write` | 用户模块写入 | `/api/users (POST/PUT/DELETE)` | 编辑角色 |
| `/api/admin/*` | 路径通配符 | 管理员路由 | 特殊路径控制 |
| `/api/users:GET` | 路径+方法 | 只能GET用户接口 | 精确控制 |

**权限匹配优先级:**
```
1. 精确匹配: /api/users:GET (最高优先级)
2. 路径通配: /api/users/*
3. 模块权限: user:read, user:write
4. 模块通配: user:*
5. 全局通配: *:* (最低优先级)
```

#### 5.3 API权限验证中间件 (带缓存)

```go
type PermissionMiddleware struct {
    cache *PermissionCache  // 三层缓存
}

func (m *PermissionMiddleware) CheckPermission() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 从JWT解析用户信息
        claims := c.MustGet("claims").(*JWTClaims)

        // 2. 从缓存获取用户权限模式 (快速!)
        patterns, err := m.cache.GetUserPermissions(claims.UserID)
        if err != nil {
            c.JSON(403, gin.H{"message": "权限获取失败"})
            c.Abort()
            return
        }

        // 3. 匹配API路径和方法
        apiPath := c.Request.URL.Path
        apiMethod := c.Request.Method

        if !m.matchPermission(patterns, apiPath, apiMethod) {
            c.JSON(403, gin.H{"message": "无权限访问此接口"})
            c.Abort()
            return
        }

        // 4. 注入数据权限范围
        c.Set("userID", claims.UserID)
        c.Set("dataScope", claims.DataScope)
        c.Set("deptID", claims.DeptID)

        c.Next()
    }
}

// 权限匹配算法 (支持通配符和模式)
func (m *PermissionMiddleware) matchPermission(
    patterns []string,
    apiPath string,
    method string,
) bool {
    for _, pattern := range patterns {
        // 全局通配
        if pattern == "*:*" {
            return true
        }

        // 路径+方法精确匹配: /api/users:GET
        if strings.Contains(pattern, ":") {
            parts := strings.Split(pattern, ":")
            if len(parts) == 2 {
                pathPattern, methodPattern := parts[0], parts[1]
                if matchPath(pathPattern, apiPath) && method == methodPattern {
                    return true
                }
            }
        }

        // 路径通配: /api/users/*
        if strings.HasSuffix(pattern, "/*") {
            prefix := strings.TrimSuffix(pattern, "/*")
            if strings.HasPrefix(apiPath, prefix) {
                return true
            }
        }

        // 模块权限: user:read, user:write, user:*
        if module, action := parseModulePermission(pattern); module != "" {
            if matchModulePermission(module, action, apiPath, method) {
                return true
            }
        }
    }
    return false
}

// 模块权限解析
func parseModulePermission(pattern string) (module, action string) {
    // user:* → module=user, action=*
    // user:read → module=user, action=read
    // user:write → module=user, action=write
    parts := strings.Split(pattern, ":")
    if len(parts) == 2 {
        return parts[0], parts[1]
    }
    return "", ""
}

// 模块权限匹配
func matchModulePermission(module, action, apiPath, method string) bool {
    // 判断API路径是否属于该模块
    // 例: /api/users/123 属于 user 模块
    if !strings.Contains(apiPath, "/"+module) {
        return false
    }

    // 模块通配符
    if action == "*" {
        return true
    }

    // 读写权限判断
    if action == "read" && method == "GET" {
        return true
    }
    if action == "write" && (method == "POST" || method == "PUT" || method == "DELETE") {
        return true
    }

    return false
}
```

#### 5.4 API权限管理UI (分层树形结构)

**前端实现 - 不要500个checkbox!**

```tsx
// API权限树结构
interface PermissionNode {
  key: string;           // user:* 或 /api/users/*
  title: string;         // "用户管理"
  type: 'module' | 'action' | 'path';
  children?: PermissionNode[];
}

// 权限树示例
const permissionTree: PermissionNode[] = [
  {
    key: 'system',
    title: '📁 系统管理',
    type: 'module',
    children: [
      {
        key: 'user:*',
        title: '👥 用户管理 (全部权限)',
        type: 'module',
        children: [
          { key: 'user:read', title: '🔍 只读', type: 'action' },
          { key: 'user:write', title: '✏️ 编辑', type: 'action' },
          { key: 'user:list', title: 'GET /api/users', type: 'path' },
          { key: 'user:create', title: 'POST /api/users', type: 'path' },
        ]
      },
      {
        key: 'role:*',
        title: '🔐 角色管理 (全部权限)',
        type: 'module',
      }
    ]
  },
  {
    key: 'business',
    title: '📁 业务管理',
    type: 'module',
    children: [
      { key: 'product:*', title: '📦 产品管理', type: 'module' },
      { key: 'order:*', title: '📋 订单管理', type: 'module' },
    ]
  }
];

// React组件
function PermissionTreeSelect() {
  const [checkedKeys, setCheckedKeys] = useState<string[]>([]);

  const handleCheck = (checked: string[]) => {
    // 智能展开: 勾选 user:* 自动勾选子权限
    setCheckedKeys(expandCheckedKeys(checked));
  };

  return (
    <Tree
      checkable
      checkedKeys={checkedKeys}
      onCheck={handleCheck}
      treeData={permissionTree}
      defaultExpandAll
    />
  );
}
```

---

### 6. 数据权限过滤模块

#### 6.1 数据权限实现

**后端SQL自动拼接:**
```go
// 根据用户数据权限范围,自动拼接WHERE条件
func BuildDataScopeSQL(userID int64, dataScope string, deptID int64) string {
    switch dataScope {
    case "ALL":
        return "" // 不限制
    case "DEPT_AND_CHILD":
        return "dept_id IN (SELECT id FROM dept WHERE find_in_set(dept_id, ancestors))"
    case "DEPT_ONLY":
        return "dept_id = " + deptID
    case "SELF_ONLY":
        return "create_by = " + userID
    }
}
```

#### 6.2 数据权限范围定义

| 数据范围代码 | 说明 | 应用场景 |
|------------|------|---------|
| `ALL` | 全部数据 | 超级管理员 |
| `DEPT_AND_CHILD` | 本部门及子部门 | 部门经理 |
| `DEPT_ONLY` | 仅本部门 | 部门主管 |
| `SELF_ONLY` | 仅本人 | 普通员工 |
| `CUSTOM` | 自定义部门 | 特殊角色(可指定多个部门) |

---

## 🔐 认证与授权流程

### 登录流程
```
1. 用户输入账号密码
   ↓
2. 后端验证密码 (bcrypt)
   ↓
3. 生成JWT Token (包含 userID, 角色, 数据权限范围)
   ↓
4. 加载用户权限模式 → 写入Redis缓存
   ↓
5. 返回Token + 用户信息 + 权限列表
   ↓
6. 前端存储Token (localStorage)
   ↓
7. 前端根据权限列表生成路由和菜单
```

### API调用流程 (带三层缓存优化)
```
前端发起请求 (携带JWT Token)
   ↓
后端中间件验证Token
   ↓
解析用户角色和权限
   ↓
三层缓存查询权限 ⚡
   ├─ Layer 1: 本地内存 (5分钟TTL) → 命中率 80%
   ├─ Layer 2: Redis (30分钟TTL) → 命中率 95%
   └─ Layer 3: MySQL → 写回缓存
   ↓
模式匹配API权限
   ├── 无权限 → 返回403
   └── 有权限 → 注入数据权限范围
       ↓
   业务逻辑层
       ↓
   DAO层自动拼接数据权限SQL
       ↓
   返回过滤后的数据
```

### 三层权限缓存架构 ⚡ (核心优化)

**问题**: 500+接口每次请求都查数据库?性能崩溃!

**解决方案**: 本地缓存 + Redis + MySQL 三层架构

```go
type PermissionCache struct {
    local *sync.Map       // 本地进程缓存 (5分钟TTL)
    redis *redis.Client   // Redis分布式缓存 (30分钟TTL)
    db    *gorm.DB       // MySQL持久化
}

// 获取用户权限 (三层查询)
func (c *PermissionCache) GetUserPermissions(userID int64) ([]string, error) {
    cacheKey := fmt.Sprintf("user:permissions:%d", userID)

    // Layer 1: 本地内存缓存 (最快,命中率80%)
    if val, ok := c.local.Load(cacheKey); ok {
        if cached, ok := val.(*CachedPermission); ok {
            if time.Now().Before(cached.ExpireAt) {
                return cached.Permissions, nil
            }
        }
    }

    // Layer 2: Redis缓存 (快,命中率95%)
    permsJSON, err := c.redis.Get(ctx, cacheKey).Result()
    if err == nil {
        var perms []string
        json.Unmarshal([]byte(permsJSON), &perms)

        // 回写本地缓存
        c.local.Store(cacheKey, &CachedPermission{
            Permissions: perms,
            ExpireAt:    time.Now().Add(5 * time.Minute),
        })
        return perms, nil
    }

    // Layer 3: MySQL数据库 (慢,命中率100%)
    perms := c.loadFromDatabase(userID)

    // 回写Redis和本地缓存
    permsJSON, _ = json.Marshal(perms)
    c.redis.Set(ctx, cacheKey, permsJSON, 30*time.Minute)
    c.local.Store(cacheKey, &CachedPermission{
        Permissions: perms,
        ExpireAt:    time.Now().Add(5 * time.Minute),
    })

    return perms, nil
}

// 从数据库加载权限
func (c *PermissionCache) loadFromDatabase(userID int64) []string {
    var permissions []string

    // 查询用户 → 角色 → 权限模式
    c.db.Raw(`
        SELECT DISTINCT rp.permission_pattern
        FROM sys_user_role ur
        JOIN sys_role_permission rp ON ur.role_id = rp.role_id
        WHERE ur.user_id = ? AND rp.status = 1
    `, userID).Scan(&permissions)

    return permissions
}

// 清除用户权限缓存 (权限变更时调用)
func (c *PermissionCache) ClearUserCache(userID int64) {
    cacheKey := fmt.Sprintf("user:permissions:%d", userID)

    // 清除本地缓存
    c.local.Delete(cacheKey)

    // 清除Redis缓存
    c.redis.Del(ctx, cacheKey)
}

// 清除角色权限缓存 (角色变更时调用)
func (c *PermissionCache) ClearRoleCache(roleID int64) {
    // 查询该角色下所有用户
    var userIDs []int64
    c.db.Raw(`
        SELECT user_id FROM sys_user_role WHERE role_id = ?
    `, roleID).Scan(&userIDs)

    // 批量清除缓存
    for _, userID := range userIDs {
        c.ClearUserCache(userID)
    }
}
```

### 缓存性能分析

| 缓存层 | TTL | 命中率 | 响应时间 | 容量 |
|-------|-----|--------|---------|------|
| 本地内存 | 5分钟 | 80% | < 1ms | 进程内存限制 |
| Redis | 30分钟 | 95% | < 10ms | 几乎无限 |
| MySQL | 永久 | 100% | 10-50ms | 持久化 |

**性能提升:**
- 无缓存: 每次请求查MySQL (20ms) → 2000 QPS = 40s总延迟
- 三层缓存: 80%命中本地(1ms) + 15%命中Redis(10ms) + 5%查MySQL(20ms)
- **延迟降低 95%!**

---

## 📊 数据库设计 (优化版)

### 核心表结构

#### 用户表 (sys_user)
```sql
id            bigint       PK
username      varchar(50)  唯一,用户名,索引
password      varchar(100) bcrypt加密
nick_name     varchar(50)  昵称
email         varchar(100)
phone         varchar(20)
dept_id       bigint       所属部门ID,索引
status        tinyint      状态(0正常 1停用)
create_time   datetime
update_time   datetime

INDEX idx_username (username)
INDEX idx_dept_id (dept_id)
INDEX idx_status (status)
```

#### 角色表 (sys_role)
```sql
id            bigint       PK
role_name     varchar(50)  角色名称
role_key      varchar(50)  唯一,角色权限字符串,索引
data_scope    varchar(20)  数据权限范围(ALL/DEPT_AND_CHILD/DEPT_ONLY/SELF_ONLY)
status        tinyint
sort          int
create_time   datetime

INDEX idx_role_key (role_key)
```

#### 部门表 (sys_dept)
```sql
id            bigint       PK
parent_id     bigint       父部门ID,索引
ancestors     varchar(500) 祖级列表(逗号分隔,用于查询子部门),索引
dept_name     varchar(50)
sort          int
leader        varchar(50)  负责人
phone         varchar(20)
status        tinyint

INDEX idx_parent_id (parent_id)
INDEX idx_ancestors (ancestors)
```

#### 菜单表 (sys_menu)
```sql
id            bigint       PK
menu_name     varchar(50)
parent_id     bigint       索引
menu_type     char(1)      类型(D目录 M菜单 B按钮)
path          varchar(200) 路由地址
component     varchar(200) 组件路径
perms         varchar(100) 权限标识
icon          varchar(100)
sort          int
visible       tinyint      是否显示
status        tinyint

INDEX idx_parent_id (parent_id)
```

#### ⚡ 权限模式表 (sys_role_permission) - 核心优化!
```sql
-- 不再是 sys_role_api 关联表!改为模式表!
id                  bigint        PK
role_id             bigint        角色ID,索引
permission_pattern  varchar(100)  权限模式 (user:*, /api/users/*, user:read)
permission_type     varchar(20)   模式类型 (module/path/action)
description         varchar(200)  描述
status              tinyint       状态
create_time         datetime

INDEX idx_role_id (role_id)
INDEX idx_pattern_type (permission_type)

-- 示例数据
INSERT INTO sys_role_permission VALUES
(1, 1, '*:*', 'global', '超级管理员全部权限', 1),
(2, 2, 'user:*', 'module', '用户模块全部权限', 1),
(3, 2, 'product:read', 'action', '产品模块只读', 1),
(4, 3, '/api/admin/*', 'path', '管理员路径', 1);
```

#### API文档表 (sys_api_doc) - 仅用于文档展示
```sql
-- 不用于权限验证!只用于API文档管理和展示
id            bigint       PK
api_path      varchar(200) API路径
api_method    varchar(10)  请求方法(GET/POST/PUT/DELETE)
api_module    varchar(50)  所属模块 (user/product/order)
description   varchar(200)
create_time   datetime

INDEX idx_module (api_module)

-- 由路由扫描器自动同步,不手动维护
```

#### 用户-角色关联表 (sys_user_role)
```sql
user_id       bigint       FK,索引
role_id       bigint       FK,索引

PRIMARY KEY (user_id, role_id)
INDEX idx_user_id (user_id)
INDEX idx_role_id (role_id)
```

#### 角色-菜单关联表 (sys_role_menu)
```sql
role_id       bigint       FK,索引
menu_id       bigint       FK,索引

PRIMARY KEY (role_id, menu_id)
INDEX idx_role_id (role_id)
INDEX idx_menu_id (menu_id)
```

### 🔥 关键优化对比

| 旧设计 (错误) | 新设计 (正确) | 优势 |
|-------------|-------------|------|
| sys_role_api (5000+条记录) | sys_role_permission (50-100条) | 减少 98% 数据量 |
| 每个API一条记录 | 权限模式匹配 | 管理成本降低 100倍 |
| 查询500个API权限 | 查询10个模式 | 查询速度提升 50倍 |
| 新增API需手动配置 | 自动模式匹配 | 零维护成本 |

---

## 🎨 前端实现要点

### 1. 权限指令 (React Hooks)

```tsx
// usePermission Hook
export const usePermission = () => {
  const { permissions } = useAuth();

  const hasPermission = (permission: string) => {
    return permissions.includes(permission);
  };

  return { hasPermission };
};

// 使用示例
function UserManage() {
  const { hasPermission } = usePermission();

  return (
    <>
      {hasPermission('user:create') && (
        <Button onClick={handleCreate}>新增用户</Button>
      )}
      {hasPermission('user:delete') && (
        <Button onClick={handleDelete}>删除</Button>
      )}
    </>
  );
}
```

### 2. 动态路由生成

```tsx
// 根据后端返回的菜单权限生成路由
function generateRoutes(menuList: Menu[]): RouteObject[] {
  return menuList
    .filter(menu => menu.menuType !== 'B') // 排除按钮
    .map(menu => ({
      path: menu.path,
      element: lazy(() => import(`@/pages${menu.component}`)),
      children: menu.children ? generateRoutes(menu.children) : []
    }));
}
```

### 3. API请求拦截

```tsx
// axios请求拦截器
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response?.status === 403) {
      message.error('无权限访问');
    }
    return Promise.reject(error);
  }
);
```

---

## 🔧 后端实现要点

### 1. JWT Token生成

```go
type JWTClaims struct {
    UserID    int64    `json:"user_id"`
    Username  string   `json:"username"`
    Roles     []string `json:"roles"`
    DataScope string   `json:"data_scope"`
    DeptID    int64    `json:"dept_id"`
    jwt.RegisteredClaims
}

func GenerateToken(user *User) (string, error) {
    claims := JWTClaims{
        UserID:    user.ID,
        Username:  user.Username,
        Roles:     user.GetRoleKeys(),
        DataScope: user.GetDataScope(),
        DeptID:    user.DeptID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

### 2. 权限中间件

```go
func PermissionMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 验证Token
        token := c.GetHeader("Authorization")
        claims, err := ParseToken(token)
        if err != nil {
            c.JSON(401, gin.H{"message": "未授权"})
            c.Abort()
            return
        }

        // 2. 检查API权限
        apiPath := c.Request.URL.Path
        apiMethod := c.Request.Method

        hasPermission := checkAPIPermission(claims.Roles, apiPath, apiMethod)
        if !hasPermission {
            c.JSON(403, gin.H{"message": "无权限访问"})
            c.Abort()
            return
        }

        // 3. 注入数据权限范围到Context
        c.Set("userID", claims.UserID)
        c.Set("dataScope", claims.DataScope)
        c.Set("deptID", claims.DeptID)

        c.Next()
    }
}
```

### 3. 数据权限过滤

```go
type DataScopeFilter struct {
    UserID    int64
    DataScope string
    DeptID    int64
}

func (f *DataScopeFilter) ApplyToQuery(db *gorm.DB, tableName string) *gorm.DB {
    switch f.DataScope {
    case "ALL":
        return db // 不过滤
    case "DEPT_AND_CHILD":
        // 查询本部门及子部门
        return db.Where(fmt.Sprintf(
            "%s.dept_id IN (SELECT id FROM sys_dept WHERE FIND_IN_SET(?, ancestors) OR id = ?)",
            tableName,
        ), f.DeptID, f.DeptID)
    case "DEPT_ONLY":
        return db.Where(fmt.Sprintf("%s.dept_id = ?", tableName), f.DeptID)
    case "SELF_ONLY":
        return db.Where(fmt.Sprintf("%s.create_by = ?", tableName), f.UserID)
    default:
        return db
    }
}

// 使用示例
func (s *UserService) GetUserList(ctx *gin.Context) ([]User, error) {
    // 从Context获取数据权限
    filter := &DataScopeFilter{
        UserID:    ctx.GetInt64("userID"),
        DataScope: ctx.GetString("dataScope"),
        DeptID:    ctx.GetInt64("deptID"),
    }

    db := s.db.Table("sys_user")
    db = filter.ApplyToQuery(db, "sys_user")

    var users []User
    err := db.Find(&users).Error
    return users, err
}
```

---

## 🚀 实施计划 (单人开发 8周版)

### ⚠️ 单人开发注意事项
- **前后端切换**: 建议按模块完整开发 (后端API→前端UI), 减少切换成本
- **代码复用**: 使用代码生成工具 (如GORM Gen, Ant Design Pro模板)
- **测试策略**: 优先集成测试,单元测试仅覆盖核心逻辑
- **进度监控**: 每周五自查进度,如延期2天立即调整范围
- **技术债务**: 允许存在,但必须注释TODO标记

### Week 1-2: 基础框架搭建
**后端 (3天)**
- ✅ Go项目初始化 (Gin + GORM + Redis)
- ✅ 数据库设计和表创建 (7张核心表)
- ✅ JWT认证实现 (Token生成/验证)
- ✅ 基础错误处理和日志框架

**前端 (3天)**
- ✅ React项目初始化 (Vite + React 18)
- ✅ Ant Design Pro 模板集成
- ✅ Redux Toolkit 配置
- ✅ Axios封装 + 路由配置

**集成 (1天)**
- ✅ 登录功能联调
- ✅ Token刷新机制

**风险**: 如果卡在环境配置,最多1天必须求助社区!

---

### Week 3-4: 权限核心功能 (最关键!)
**后端 (5天)**
- ✅ 用户、角色、部门CRUD (3天)
- ✅ 菜单管理 (树形结构) (1天)
- ✅ 权限中间件 + 三层缓存 (1天)

**前端 (2天)**
- ✅ 用户管理页面 (表格+表单)
- ✅ 角色管理页面 (权限树选择)

**集成 (1天)**
- ✅ 权限验证联调
- ✅ 动态菜单渲染

**风险**: 权限树逻辑复杂,预留半天调试时间!

---

### Week 5: API权限管理 ⚡ (核心优化)
**后端 (3天)**
- ✅ API自动扫描注册 (1天)
- ✅ 权限模式匹配实现 (1天)
- ✅ sys_role_permission 表和逻辑 (1天)

**前端 (2天)**
- ✅ API权限树UI (借鉴菜单树)
- ✅ 权限模式配置界面

**风险**: 模式匹配算法需要仔细测试,预留1天测试时间!

---

### Week 6: 数据权限过滤
**后端 (3天)**
- ✅ DataScopeFilter 实现 (1天)
- ✅ SQL自动拼接逻辑 (1天)
- ✅ 4种数据范围测试 (1天)

**前端 (2天)**
- ✅ 部门树组件
- ✅ 数据权限选择UI

---

### Week 7: 日志和导入导出 (MVP功能)
**后端 (3天)**
- ✅ 登录日志 (1天)
- ✅ 用户Excel导入 (excelize库) (2天)

**前端 (2天)**
- ✅ 登录日志页面 (1天)
- ✅ 文件上传组件 (1天)

**说明**: 只做基础Excel导入,不做复杂验证!

---

### Week 8: 测试、优化、部署
**测试 (3天)**
- ✅ 集成测试 (主流程)
- ✅ 性能测试 (缓存命中率)
- ✅ Bug修复

**优化 (2天)**
- ✅ 代码重构 (去除重复代码)
- ✅ SQL性能优化 (加索引)

**部署 (2天)**
- ✅ Docker镜像构建
- ✅ 部署文档编写
- ✅ 生产环境部署

---

## 📊 工作量分解 (40天 × 8小时 = 320小时)

| 模块 | 后端 | 前端 | 测试 | 总计 |
|------|------|------|------|------|
| 基础框架 | 24h | 24h | 8h | 56h |
| 权限核心 | 40h | 16h | 8h | 64h |
| API权限 | 24h | 16h | 8h | 48h |
| 数据权限 | 24h | 16h | 8h | 48h |
| 日志导入 | 24h | 16h | - | 40h |
| 测试部署 | - | - | 64h | 64h |
| **总计** | 136h | 88h | 96h | **320h** |

**缓冲时间**: 预留20%应对突发问题 = 64小时

---

## 🔥 单人开发生存指南

### 每日工作流程
```
09:00-12:00  专注编码 (后端或前端,不要切换!)
12:00-13:00  午休
13:00-15:00  专注编码
15:00-15:30  代码提交 + 测试
15:30-17:00  另一端开发 (如果上午写后端,下午写前端)
17:00-18:00  联调 + Bug修复
```

### 每周检查点
- **周三晚**: 检查本周进度,如延期立即调整
- **周五晚**: 代码提交 + 写周报
- **周日晚**: 规划下周任务

### 求助策略
- **卡住超过2小时**: 去Stack Overflow/GitHub Issues
- **卡住超过半天**: 付费咨询 (值得!)
- **卡住超过1天**: 调整方案或砍功能

### 健康提醒 ⚠️
- **每天最多10小时**: 超时会降低效率
- **每周休息1天**: 避免burnout
- **如果生病**: 立即休息,否则进度会崩

---

## ⚡ 快速开发技巧

### 代码生成工具
```bash
# GORM Gen - 生成CRUD代码
go install gorm.io/gen/tools/gentool@latest

# Ant Design Pro - 脚手架
npm create vite@latest frontend -- --template react-ts
```

### 复用开源组件
- **前端表格**: ProTable (Ant Design Pro)
- **权限树**: Ant Design Tree + 递归组件
- **Excel导入**: excelize (Go) + antd Upload

### 性能优化懒人法
- **后端**: 只给高频查询加索引
- **前端**: 只给大表格加虚拟滚动
- **缓存**: Redis用默认配置就行

---

## 🎯 如果8周做不完怎么办?

### Plan B: 7周MVP (砍功能版)
- ❌ 砍掉: 数据导入导出
- ❌ 砍掉: 登录日志
- ❌ 简化: 部门管理 (只支持2层)

### Plan C: 6周核心版 (最小可用)
- ❌ 砍掉: API权限管理 (手动配置)
- ❌ 砍掉: 数据权限过滤 (全部用全局权限)
- ✅ 保留: 用户、角色、菜单管理

**我的建议**: 按8周计划走,如果Week 4发现进度慢,立即启动Plan B!

---

## 📈 非功能需求 (500+接口优化版)

### 性能要求
- 登录响应时间 < 500ms
- **API权限验证 < 5ms** (三层缓存优化后)
- 菜单加载 < 200ms
- **支持2000+ QPS** (单实例)
- **支持5000+ 并发用户** (水平扩展)

### 缓存策略
- **本地缓存**: 5分钟TTL, 80%命中率
- **Redis缓存**: 30分钟TTL, 95%命中率
- **权限变更**: 实时清除相关用户缓存
- **缓存预热**: 系统启动时加载热点数据

### 安全要求
- 密码bcrypt加密 (cost=10)
- JWT Token有效期24小时
- Token刷新机制 (RefreshToken 7天)
- SQL注入防护 (参数化查询)
- XSS防护 (内容转义)
- **API限流**: 单用户 100 req/min
- **暴力破解防护**: 登录失败5次锁定15分钟

### 可扩展性
- **无状态设计**: 支持水平扩展
- **Redis集群**: 分布式缓存
- **数据库读写分离**: 主从架构
- **支持多租户扩展**: 租户ID隔离

### 监控告警
- API响应时间监控
- 缓存命中率监控
- 权限验证失败告警
- 异常登录告警

---

## 🎯 核心价值

1. **开箱即用**: 提供完整的权限管理解决方案
2. **灵活扩展**: RBAC + 数据权限双重保障
3. **高性能**: Go后端高并发处理能力
4. **现代化**: React 18 + Ant Design 5最新技术栈
5. **安全可靠**: JWT认证 + 多层权限校验

---

## ⚠️ 已知风险和限制 (500+接口场景)

### 风险
1. **数据权限过滤性能**: 复杂SQL可能影响性能 (需加索引优化 + 查询优化)
2. **权限缓存一致性**: 权限变更后需要及时刷新缓存 (已通过清除机制解决)
3. **本地缓存内存占用**: 高并发时本地缓存可能占用较多内存 (建议限制大小)
4. **权限模式冲突**: 多个模式可能产生冲突 (需要明确优先级规则)

### 限制
1. **不支持动态规则引擎**: 如果需要复杂的属性判断(如时间、地点),需要升级到ABAC
2. **不支持时间维度权限**: 比如"只能在工作时间访问" (可扩展)
3. **不支持细粒度字段权限**: 只能控制到接口级,无法控制字段级 (可通过DTO过滤实现)
4. **权限模式学习成本**: 开发者需要理解通配符和模式匹配规则

### 针对500+接口的特殊注意事项
1. **必须使用Redis缓存**: 没有缓存性能会崩溃
2. **必须定期清理API文档表**: 避免废弃API堆积
3. **必须做好权限模式规划**: 避免权限碎片化
4. **必须监控缓存命中率**: 低于70%需要优化

---

## 📝 后续扩展方向 (Phase 2)

### 核心功能 (Phase 1 - MVP, 8周)
- ✅ 用户、角色、部门、菜单管理
- ✅ API权限管理 (模式匹配)
- ✅ 数据权限过滤
- ✅ JWT认证
- ✅ 基础日志 (登录日志)
- ✅ 基础导入导出 (用户Excel导入)

### Phase 2 扩展 (后续迭代, 4-6周)
1. **审计日志**: 记录所有敏感操作 (权限变更、删除操作)
2. **操作日志**: 记录所有CRUD操作 (带数据快照)
3. **在线用户管理**: 实时查看在线用户,强制下线
4. **字典管理**: 系统配置项管理
5. **高级导入导出**: 批量导入、模板生成、导出任务队列
6. **通知公告**: 系统消息推送
7. **文件管理**: OSS对象存储集成
8. **代码生成器**: 根据数据库表生成CRUD代码

### ⚠️ 单人开发优先级说明

**Phase 1 (MVP) 必须实现:**
- 权限管理核心功能 (无法妥协)
- 登录日志 (安全必需)
- 用户导入 (业务必需)

**Phase 2 暂缓原因:**
- 审计日志: 可以先用操作日志代替
- 操作日志: 性能影响大,需要异步方案
- 高级导入导出: 业务复杂度高,时间成本大

**建议策略**: 先上线MVP,收集用户反馈后再决定Phase 2优先级!

---

## 📚 参考资料

- [若依后台管理系统](https://gitee.com/y_project/RuoYi)
- [Gin Web Framework](https://gin-gonic.com/)
- [Ant Design Pro](https://pro.ant.design/)
- [Casbin权限框架](https://casbin.org/)

---

**文档版本**: v3.0 (单人开发优化版)
**创建日期**: 2025-09-30
**最后更新**: 2025-09-30
**作者**: Claude
**系统规模**: 500+ API接口, 100-500用户
**开发团队**: 1人全栈 (8周MVP)
**技术栈**: Gin + GORM + Redis + React 18 + Redux Toolkit
**状态**: ✅ 单人方案优化完成,可开始开发
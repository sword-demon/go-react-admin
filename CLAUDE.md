# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

- Layer 1: 本地内存(5 分钟 TTL, 80%命中率, <1ms)
- Layer 2: Redis(30 分钟 TTL, 95%命中率, <10ms)
- Layer 3: MySQL(持久化, 100%命中率, 10-50ms)

**关键代码位置**:

- 权限中间件: `backend/middleware/permission.go`
- 缓存实现: `backend/service/cache.go`
- 模式匹配: `backend/service/permission.go`

### 2. API 权限模式匹配(避免 500 个 API 单独配置)

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

### 3. 数据权限过滤(根据部门/角色)

**4 种数据范围**:

- `ALL`: 全部数据(超管)
- `DEPT_AND_CHILD`: 本部门+子部门(部门经理)
- `DEPT_ONLY`: 仅本部门(部门主管)
- `SELF_ONLY`: 仅本人数据(普通员工)

**实现方式**:

- 在权限中间件中注入 `dataScope` 到 Context
- DAO 层查询时自动拼接 WHERE 条件
- 示例: `WHERE dept_id IN (SELECT id FROM sys_dept WHERE FIND_IN_SET(?, ancestors))`

### 4. 核心数据表(7 张)

1. **sys_user** - 用户表(索引: username, dept_id, status)
2. **sys_role** - 角色表(索引: role_key)
3. **sys_dept** - 部门表(树形结构,索引: parent_id, ancestors)
4. **sys_menu** - 菜单表(树形结构,索引: parent_id)
5. **sys_role_permission** - 权限模式表(核心!存储 user:_, /api/admin/_ 等模式)
6. **sys_user_role** - 用户-角色关联表(复合主键)
7. **sys_role_menu** - 角色-菜单关联表(复合主键)

**注意**: `sys_api_doc` 仅用于 API 文档展示,不用于权限验证!

---

## 开发约定

### Go 后端代码规范

1. **分层架构** (严格遵守):

   ```
   controller → service → model
   (HTTP处理) → (业务逻辑) → (数据访问)
   ```

2. **错误处理**:

   - 使用统一的错误响应格式
   - Service 层返回 `error`
   - Controller 层处理 HTTP 状态码

3. **数据库查询**:

   - 所有查询必须考虑数据权限
   - 高频查询必须加索引
   - 使用 GORM 的 Preload 避免 N+1 查询

4. **缓存策略**:
   - 权限数据必须走三层缓存
   - 权限变更时必须清除相关用户缓存
   - 角色变更时批量清除该角色下所有用户缓存

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

- **需求文档**: `docs/requirement.md` (1200+行,包含完整架构设计和代码示例)
- **若依参考**: https://gitee.com/y_project/RuoYi (权限设计参考)
- **Gin 文档**: https://gin-gonic.com/
- **GORM 文档**: https://gorm.io/
- **Ant Design**: https://ant.design/

---

## 快速开始检查清单

开始开发前确认:

- [ ] MySQL 8.0+ 已安装并运行
- [ ] Redis 已安装并运行(单机部署即可)
- [ ] Go 1.23+ 已安装
- [ ] Node.js 18+ 已安装
- [ ] 已读完 `docs/requirement.md` (特别是架构部分)
- [ ] 已创建数据库表结构(7 张核心表)
- [ ] 后端能跑通 Hello World
- [ ] 前端能跑通 Vite 开发服务器

---

**重要**: 这是一个雄心勃勃的单人项目!优先完成 MVP,不要追求完美!能上线比完美更重要!

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

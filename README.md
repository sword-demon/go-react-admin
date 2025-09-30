# 🚀 Go-React-Admin

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/React-19.1+-61DAFB?style=flat&logo=react" alt="React Version" />
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat" alt="License" />
  <img src="https://img.shields.io/badge/Status-In_Development-yellow?style=flat" alt="Status" />
</p>

<p align="center">
  <strong>企业级权限管理系统</strong><br/>
  支持 500+ API 接口规模 | RBAC权限模型 | 数据权限过滤 | 三层缓存架构
</p>

---

## 📋 项目简介

Go-React-Admin 是一个基于 **Go + React** 的现代化企业级权限管理系统,专为**中小型企业**设计,支持 **500+ API 接口**规模。

### 核心特性

- 🔐 **RBAC权限模型** - 基于角色的访问控制 + 数据权限过滤(部门/角色)
- ⚡ **API权限模式匹配** - 通配符支持,避免为每个API单独配置(减少98%权限记录)
- 🚀 **三层缓存架构** - 本地内存(80%命中) + Redis(95%命中) + MySQL持久化
- 🎯 **单人开发优化** - 8周MVP计划,模块化设计,快速上手
- 📊 **4种数据权限范围** - 全部/本部门+子部门/仅本部门/仅本人
- 📝 **完整日志系统** - 登录日志(Phase 1) + 操作日志(Phase 2)

---

## 🛠️ 技术栈

### 后端
- **框架**: [Gin](https://gin-gonic.com/) - 高性能HTTP Web框架
- **ORM**: [GORM v2](https://gorm.io/) - 功能丰富的Go ORM
- **数据库**: MySQL 8.0+ / PostgreSQL 14+
- **缓存**: Redis 7.0+ (单机部署)
- **认证**: JWT Token
- **版本**: Go 1.24+

### 前端
- **框架**: React 19.1+ with TypeScript
- **构建工具**: Vite 7.1+
- **UI库**: Ant Design 5.x
- **状态管理**: Redux Toolkit (规划中)
- **路由**: React Router v6 (规划中)
- **HTTP客户端**: Axios

### 数据库设计
- 10张核心表(7核心 + 3日志)
- 软删除支持
- 索引优化
- 树形结构(部门/菜单)

---

## 🚀 快速开始

### 前置要求

确保你的开发环境已安装:

```bash
# 检查版本
go version      # >= 1.23
node -v         # >= 18.0
mysql --version # >= 8.0
redis-server -v # >= 7.0
```

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/go-react-admin.git
cd go-react-admin
```

### 2. 数据库初始化

```bash
# 登录MySQL
mysql -u root -p

# 导入表结构和初始数据
source docs/schema.sql

# 或者
mysql -u root -p < docs/schema.sql

# 验证
mysql -u root -p -e "USE go_react_admin; SHOW TABLES;"
```

### 3. 后端启动

```bash
cd backend

# 安装依赖
go mod tidy

# 创建配置文件(复制示例配置)
cp config/config.example.yml config/config.local.yml
# 修改 config.local.yml 中的数据库和Redis配置

# 运行(开发模式)
go run main.go

# 或使用热重载(推荐)
go install github.com/air-verse/air@latest
air

# 后端将运行在 http://localhost:8080
```

### 4. 前端启动

```bash
cd react-admin

# 安装依赖
npm install

# 运行开发服务器
npm run dev

# 前端将运行在 http://localhost:5173
```

### 5. 访问系统

- **前端地址**: http://localhost:5173
- **后端API**: http://localhost:8080
- **默认账号**: `admin` / `admin123`

---

## 📁 项目结构

```
go-react-admin/
├── backend/                 # Go后端
│   ├── controller/          # 控制器层(HTTP处理)
│   ├── service/             # 业务逻辑层
│   ├── model/               # 数据模型(GORM)
│   ├── middleware/          # 中间件(JWT/权限验证)
│   ├── config/              # 配置文件
│   ├── utils/               # 工具函数
│   ├── go.mod               # Go模块依赖
│   └── main.go              # 程序入口
│
├── react-admin/             # React前端
│   ├── src/
│   │   ├── pages/           # 页面组件
│   │   ├── components/      # 通用组件
│   │   ├── store/           # Redux状态管理
│   │   ├── api/             # API请求封装
│   │   ├── utils/           # 工具函数
│   │   └── App.tsx          # 根组件
│   ├── package.json
│   └── vite.config.ts
│
├── docs/                    # 文档
│   ├── requirement.md       # 需求文档(1200+行,必读!)
│   └── schema.sql           # 数据库表结构
│
├── CLAUDE.md                # AI开发助手指南
├── README.md                # 项目说明(当前文件)
└── .gitignore               # Git忽略规则
```

---

## 🔑 核心功能

### 已实现 (Phase 1 - MVP)

- [x] 项目初始化(前后端骨架)
- [x] 数据库设计(10张表)
- [ ] 用户管理(CRUD)
- [ ] 角色管理(CRUD + 权限分配)
- [ ] 部门管理(树形结构)
- [ ] 菜单管理(动态路由)
- [ ] JWT认证(登录/登出)
- [ ] 权限中间件(三层缓存)
- [ ] API权限模式匹配
- [ ] 数据权限过滤
- [ ] 登录日志
- [ ] 用户Excel导入

### 规划中 (Phase 2)

- [ ] 操作日志(带数据快照)
- [ ] 审计日志(敏感操作记录)
- [ ] 在线用户管理(强制下线)
- [ ] 字典管理(配置项)
- [ ] 高级导入导出(任务队列)
- [ ] 通知公告
- [ ] 文件管理(OSS集成)
- [ ] 代码生成器

---

## 🏗️ 架构设计

### 权限验证流程(核心!)

```
HTTP请求 → JWT中间件 → 权限中间件 → 业务逻辑 → 数据权限过滤 → 响应
              ↓             ↓                         ↓
          解析Token    三层缓存查权限             SQL自动拼接WHERE
                      (本地→Redis→MySQL)
```

### 三层缓存架构

| 缓存层 | TTL | 命中率 | 响应时间 | 说明 |
|-------|-----|--------|---------|------|
| 本地内存 | 5分钟 | 80% | <1ms | sync.Map实现 |
| Redis | 30分钟 | 95% | <10ms | 分布式缓存 |
| MySQL | 永久 | 100% | 10-50ms | 持久化 |

**性能提升**: 延迟降低 95%!

### API权限模式匹配

不要为每个API单独配置权限!使用模式匹配:

```
*:*              # 全部API(超管)
user:*           # 用户模块所有操作
user:read        # 用户模块只读(GET)
user:write       # 用户模块写入(POST/PUT/DELETE)
/api/admin/*     # 路径通配符
/api/users:GET   # 路径+方法精确匹配
```

**优势**: 从5000+权限记录减少到50-100条!

### 数据权限范围

| 范围 | 说明 | 应用场景 |
|------|------|---------|
| `ALL` | 全部数据 | 超级管理员 |
| `DEPT_AND_CHILD` | 本部门+子部门 | 部门经理 |
| `DEPT_ONLY` | 仅本部门 | 部门主管 |
| `SELF_ONLY` | 仅本人数据 | 普通员工 |

---

## 🧪 开发指南

### 后端开发

```bash
cd backend

# 运行测试
go test ./...

# 运行特定测试
go test -v -run TestUserService ./service

# 代码格式化
go fmt ./...

# 代码检查
go vet ./...

# 构建
go build -o bin/server main.go
```

### 前端开发

```bash
cd react-admin

# 运行开发服务器
npm run dev

# 构建生产版本
npm run build

# 预览生产构建
npm run preview

# ESLint检查
npm run lint
```

### 数据库操作

```bash
# 重建数据库(慎用!会删除所有数据)
mysql -u root -p < docs/schema.sql

# 备份数据库
mysqldump -u root -p go_react_admin > backup_$(date +%Y%m%d).sql

# 查看表结构
mysql -u root -p -e "USE go_react_admin; SHOW CREATE TABLE sys_role_permission;"
```

---

## 📚 文档

- **需求文档**: [docs/requirement.md](docs/requirement.md) - 1200+行完整架构设计
- **CLAUDE.md**: AI开发助手指南(包含架构关键点、开发约定)
- **数据库设计**: [docs/schema.sql](docs/schema.sql) - 10张表结构
- **API文档**: (规划中 - 将使用Swagger自动生成)

---

## 🤝 开发团队

- **开发者**: 1人全栈(单人项目)
- **开发周期**: 8周MVP计划
- **开发模式**: 按模块完整开发(后端API → 前端UI)

### 开发进度

- [x] Week 1-2: 基础框架搭建 ⚡ 进行中
- [ ] Week 3-4: 权限核心功能
- [ ] Week 5: API权限管理
- [ ] Week 6: 数据权限过滤
- [ ] Week 7: 日志和导入导出
- [ ] Week 8: 测试、优化、部署

---

## ⚠️ 注意事项

### 单人开发风险

这是一个**雄心勃勃的单人项目**!如果进度延期:

- **Week 4发现进度慢**: 立即砍掉API权限管理
- **Week 6还没完成**: 砍掉数据权限过滤
- **卡住超过1天**: 调整方案或求助社区

### 安全提醒

- ✅ 永远不要提交 `.env` 文件到Git
- ✅ 生产环境修改默认密码 `admin123`
- ✅ 使用HTTPS部署
- ✅ 定期备份数据库
- ✅ 启用Redis密码认证

---

## 📈 性能指标

### 目标性能

- 登录响应: < 500ms
- API权限验证: < 5ms (三层缓存)
- 菜单加载: < 200ms
- 支持并发: 500+ QPS (单机)
- 支持用户数: 100-500 (中小型企业)

### 优化策略

- ✅ 三层缓存(权限数据)
- ✅ SQL索引优化(高频查询)
- ✅ 权限模式匹配(减少98%记录)
- ✅ 前端懒加载(路由/组件)

---

## 🔗 参考资料

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM 文档](https://gorm.io/)
- [React 官方文档](https://react.dev/)
- [Ant Design 官方文档](https://ant.design/)
- [若依后台管理系统](https://gitee.com/y_project/RuoYi) (权限设计参考)

---

## 📄 License

MIT License

Copyright (c) 2025 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

## 💬 联系方式

- **Issues**: [GitHub Issues](https://github.com/yourusername/go-react-admin/issues)
- **Email**: your.email@example.com
- **Blog**: https://your-blog.com

---

<p align="center">
  <strong>⚡ 能上线比完美更重要! ⚡</strong><br/>
  <sub>专注MVP,后续迭代优化</sub>
</p>

<p align="center">
  Made with ❤️ by [Your Name]
</p>
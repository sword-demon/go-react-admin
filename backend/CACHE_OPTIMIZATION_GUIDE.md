# 缓存优化功能使用指南

> **版本**: v1.1.0
> **新增功能**: 缓存预热 + 监控API + 性能压测
> **日期**: 2025-10-01

---

## 🔥 功能一: 缓存预热 (Cache Warmup)

### 什么是缓存预热?

系统启动时预先加载热点数据到缓存,避免冷启动时大量请求直接打到数据库。

### 使用方法

#### 1. 在 main.go 中启用缓存预热

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/sword-demon/go-react-admin/internal/pkg/cache"
    "github.com/sword-demon/go-react-admin/internal/admin/biz/permission"
)

func main() {
    // ... 初始化代码 ...

    // 1. 创建PermissionBiz (实现了PermissionLoader接口)
    permBiz := permission.NewPermissionBiz(dataStore, localCache, redisClient)

    // 2. 配置预热参数
    warmupConfig := &cache.WarmupConfig{
        SuperAdminUserIDs: []uint64{1},           // 超级管理员用户ID
        CommonRoleIDs:     []uint64{1, 2, 3},     // 常用角色ID (管理员/经理/普通用户)
        Concurrency:       5,                     // 并发数
        Timeout:           30 * time.Second,      // 超时时间
        EnableLogging:     true,                  // 启用日志
    }

    // 3. 创建预热器
    warmer := cache.NewPermissionWarmer(warmupConfig, permBiz)

    // 4. 异步执行预热 (不阻塞启动)
    go func() {
        log.Println("🔥 Starting cache warmup...")
        stats, err := warmer.Warm(context.Background())
        if err != nil {
            log.Printf("⚠️  Warmup error: %v", err)
        } else {
            log.Printf("✅ Warmup completed: %d/%d success, took %v",
                stats.SuccessCount, stats.TotalItems, stats.Duration)
        }
    }()

    // ... 启动服务器 ...
}
```

#### 2. 预热效果

**启动日志示例**:
```
🔥 Starting cache warmup: 1 users, 3 roles
  → Warming user permissions: userID=1
  ✓ Success: user:1
  → Warming role permissions: roleID=1
  ✓ Success: role:1
  → Warming role permissions: roleID=2
  ✓ Success: role:2
  → Warming role permissions: roleID=3
  ✓ Success: role:3
🎉 Cache warmup completed: 4/4 success, 0 failed, took 125ms
```

**性能对比**:
```
┌─────────────────┬──────────┬───────────┐
│ 请求类型         │ 冷启动   │ 预热后     │
├─────────────────┼──────────┼───────────┤
│ 超管权限检查     │ 30ms     │ 0.5ms     │
│ 普通角色权限     │ 30ms     │ 0.5ms     │
│ 首次登录用户     │ 30ms     │ 30ms      │
└─────────────────┴──────────┴───────────┘
```

---

## 📊 功能二: 缓存监控API

### API端点

#### 1. 查看缓存统计

**请求**:
```bash
GET /api/v1/cache/stats
```

**响应**:
```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "hits": 8532,
    "misses": 1245,
    "hit_rate": 87.28,
    "size": 756,
    "max_size": 10000,
    "health": "excellent"
  }
}
```

**健康状态评级**:
- `excellent`: 命中率 >= 80%
- `good`: 命中率 >= 60%
- `fair`: 命中率 >= 40%
- `poor`: 命中率 >= 20%
- `critical`: 命中率 < 20%

#### 2. 清除指定缓存

**清除用户缓存**:
```bash
DELETE /api/v1/cache?user_id=123
```

**清除角色缓存**:
```bash
DELETE /api/v1/cache?role_id=5
```

**按前缀批量清除**:
```bash
DELETE /api/v1/cache?prefix=user:123:
```

**响应**:
```json
{
  "code": 0,
  "msg": "cache cleared",
  "data": {
    "cleared_count": 5
  }
}
```

#### 3. 触发手动预热

```bash
POST /api/v1/cache/warmup
```

**响应**:
```json
{
  "code": 0,
  "msg": "warmup triggered",
  "data": {
    "status": "queued"
  }
}
```

### 集成到main.go

```go
package main

import (
    "github.com/gin-gonic/gin"
    v1 "github.com/sword-demon/go-react-admin/internal/admin/controller/v1"
)

func main() {
    // ... 初始化代码 ...

    // 创建CacheController
    cacheCtrl := v1.NewCacheController(permBiz, localCache)

    // 注册路由
    r := gin.Default()
    api := r.Group("/api/v1")
    {
        cache := api.Group("/cache")
        {
            cache.GET("/stats", cacheCtrl.GetStats)           // 查看统计
            cache.DELETE("", cacheCtrl.ClearCache)            // 清除缓存
            cache.POST("/warmup", cacheCtrl.WarmupCache)      // 手动预热
        }
    }

    r.Run(":8080")
}
```

---

## 🏎️ 功能三: 性能压测

### 运行Benchmark

#### 1. 本地缓存读性能

```bash
cd backend
go test -bench=BenchmarkLocalCacheGet ./internal/pkg/cache/
```

**预期输出**:
```
BenchmarkLocalCacheGet-8    50000000    25.3 ns/op    0 B/op    0 allocs/op
```
解读: 每次读操作 25纳秒, 零内存分配

#### 2. 本地缓存写性能

```bash
go test -bench=BenchmarkLocalCacheSet ./internal/pkg/cache/
```

**预期输出**:
```
BenchmarkLocalCacheSet-8    5000000    312 ns/op    224 B/op    3 allocs/op
```
解读: 每次写操作 312纳秒, 224字节内存分配

#### 3. 并发读写测试

```bash
go test -bench=BenchmarkLocalCacheConcurrent ./internal/pkg/cache/
```

**预期输出**:
```
BenchmarkLocalCacheConcurrent-8    2000000    785 ns/op    64 B/op    1 allocs/op
```
解读: 8读2写并发场景, 平均 785纳秒/操作

#### 4. LRU淘汰性能

```bash
go test -bench=BenchmarkLRUEviction ./internal/pkg/cache/
```

**预期输出**:
```
BenchmarkLRUEviction-8    10000000    142 ns/op    96 B/op    2 allocs/op
```
解读: LRU淘汰算法开销 142纳秒/次

#### 5. 预热并发性能对比

```bash
go test -bench=BenchmarkWarmupConcurrency ./internal/pkg/cache/
```

**预期输出**:
```
BenchmarkWarmupConcurrency/Concurrency-1-8     20     65432109 ns/op
BenchmarkWarmupConcurrency/Concurrency-5-8     50     23456789 ns/op
BenchmarkWarmupConcurrency/Concurrency-10-8    80     15234567 ns/op
BenchmarkWarmupConcurrency/Concurrency-20-8    100    12345678 ns/op
```
解读: 并发度越高, 预热越快 (但有边际递减)

#### 6. 完整性能报告

```bash
go test -bench=. -benchmem -benchtime=3s ./internal/pkg/cache/ > bench_results.txt
```

### 运行功能测试

#### 1. 缓存命中率模拟

```bash
go test -v -run=TestCacheHitRateSimulation ./internal/pkg/cache/
```

**预期输出**:
```
=== RUN   TestCacheHitRateSimulation
    cache_benchmark_test.go:85: Cache Hit Rate: 87.45%
    cache_benchmark_test.go:86: Hits: 8745, Misses: 1255
    cache_benchmark_test.go:87: Cache Size: 1000/1000
--- PASS: TestCacheHitRateSimulation (0.12s)
```

#### 2. 前缀批量删除测试

```bash
go test -v -run=TestCachePrefixDelete ./internal/pkg/cache/
```

**预期输出**:
```
=== RUN   TestCachePrefixDelete
    cache_benchmark_test.go:135: Initial cache size: 5000
    cache_benchmark_test.go:140: Deleted 5 entries in 125µs
--- PASS: TestCachePrefixDelete (0.01s)
```

---

## 📈 性能优化建议

### 1. 根据业务调整缓存容量

**小型系统** (< 1000 活跃用户):
```go
localCache := cache.NewLocalCache(1000, 5*time.Minute)
```

**中型系统** (1000-10000 活跃用户):
```go
localCache := cache.NewLocalCache(5000, 5*time.Minute)
```

**大型系统** (> 10000 活跃用户):
```go
localCache := cache.NewLocalCache(10000, 5*time.Minute)
```

### 2. 根据并发量调整预热并发度

**测试环境**:
```go
Concurrency: 5   // 温和预热,避免数据库压力
```

**生产环境** (充足资源):
```go
Concurrency: 20  // 快速预热,缩短启动时间
```

**生产环境** (资源受限):
```go
Concurrency: 3   // 保守预热,避免启动时资源争抢
```

### 3. 监控告警阈值

```yaml
# Prometheus告警规则示例
- alert: CacheHitRateLow
  expr: cache_hit_rate < 60
  for: 5m
  annotations:
    summary: "缓存命中率过低 ({{ $value }}%)"
    description: "建议检查缓存容量或TTL设置"

- alert: CacheSizeNearLimit
  expr: cache_size / cache_max_size > 0.9
  for: 5m
  annotations:
    summary: "缓存容量接近上限 ({{ $value }}%)"
    description: "建议增加maxSize或检查是否有内存泄漏"
```

---

## 🔍 故障排查

### 问题: 预热失败

**症状**:
```
⚠️  Warmup error: context deadline exceeded
```

**解决**:
1. 增加Timeout: `Timeout: 60 * time.Second`
2. 减少并发度: `Concurrency: 3`
3. 减少预热数据量

### 问题: 命中率低于预期

**症状**:
```json
{ "hit_rate": 35.6, "health": "poor" }
```

**诊断**:
```bash
# 1. 检查缓存容量是否足够
curl http://localhost:8080/api/v1/cache/stats

# 2. 检查是否有大量不同用户访问 (导致缓存频繁淘汰)
# 3. 考虑增加maxSize或延长TTL
```

### 问题: 内存占用过高

**症状**:
系统内存使用持续增长

**解决**:
1. 检查maxSize配置: 是否设置了上限?
2. 启动清理Worker:
```go
stopCh := localCache.StartCleanupWorker(1 * time.Minute)
defer close(stopCh)
```

---

## 📚 参考资料

- [CACHE_DESIGN_ZH.md](./CACHE_DESIGN_ZH.md) - 完整架构设计文档
- [internal/pkg/cache/](./internal/pkg/cache/) - 源码实现
- [Go Benchmark Tutorial](https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go)

---

**作者**: Claude Code
**版本**: v1.1.0
**最后更新**: 2025-10-01

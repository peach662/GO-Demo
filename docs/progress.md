# 学习进度

## 当前状态

当前阶段：企业级后端第 1 步，正在把内存版 Todo API 迁移到 MySQL。

当前项目：`D:\demo\awesomeProject`

## 已完成

- [x] Go module、Gin 和健康检查接口。
- [x] Todo 查询、按 ID 查询、创建和更新状态接口。
- [x] 请求参数校验和统一错误响应。
- [x] Handler -> Service 分层。
- [x] Service 单元测试和 HTTP 路由测试。
- [x] `sync.RWMutex` 保护内存 Todo 数据。
- [x] `List()` 返回切片副本，避免泄露内部数据。
- [x] goroutine、`sync.WaitGroup`、channel 基础练习。
- [x] 并发创建 Todo，验证数量与 ID 唯一性。
- [x] 使用 Docker Compose 启动 MySQL 8.4。
- [x] 创建 `migrations/001_create_todos.sql` 并在 MySQL 中创建 `todos` 表。
- [x] 完成 `internal/database/mysql.go`：连接池创建、`Ping` 验证和失败关闭。
- [x] 在 `main.go` 启动阶段调用 `database.OpenMySQL()`，确认 Go 服务可以连接 MySQL。
- [x] 通过 `GET /health` 验证 Gin 服务可访问。
- [x] 定义 `todo.Repository` 接口，明确查询、创建和更新的数据访问契约。
- [x] 创建 `MySQLRepository` 和构造函数，注入 `*sql.DB` 连接池。
- [x] 实现 `MySQLRepository.List`：使用 `QueryContext` 查询并扫描 Todo 列表。

## 最近验证

用户已确认以下命令通过：

```powershell
go test -race -count=1 ./...
```

结果：根包和 `internal/todo` 包均通过，未报告数据竞争。

用户已确认以下命令通过：

```powershell
go test ./...
```

结果：根包、`internal/database` 和 `internal/todo` 包均可编译；其中 `internal/database` 暂无测试文件。

已通过只读 `gofmt -d internal/database/mysql.go` 检查，未发现格式差异。

用户已执行 `go run .`，Gin 成功注册全部路由并监听 `:9090`。由于 MySQL `Ping` 在 Gin 启动前执行，启动过程中未出现 `panic` 也证明 Go 到 MySQL 的连接成功。

用户已在另一终端执行 `Invoke-RestMethod http://localhost:9090/health`，返回 `code=0`、`data=pong`、`message=ok`。

用户已执行 `go test ./...`，根包、`internal/database` 和 `internal/todo` 包均通过。

用户已执行 `go fmt ./...` 和 `go test ./...`，`MySQLRepository.List` 已通过编译检查；目前还没有真实数据库集成测试。

## 已知依赖说明

Gin `v1.12.0` 依赖 `github.com/goccy/go-yaml v1.19.2`，但当前下载到的该版本缺少 Gin 运行时导入的包，导致 `go mod tidy` 和完整 race 测试失败。

`go.mod` 已通过 `replace` 临时使用兼容的 `github.com/goccy/go-yaml v1.18.0`，并已验证：

```powershell
go test -race -count=1 ./...
```

通过。不要删除这条 `replace`，除非升级 Gin 或确认上游 YAML 发布问题已修复，并重新运行完整测试。

Docker 环境也已确认可用：Docker Desktop 4.55.0，Docker Engine 29.1.3。

用户已确认 MySQL 容器状态：

```text
awesome-project-mysql   mysql:8.4   Up (healthy)
0.0.0.0:13306 -> 3306/tcp
```

用户已使用 `DESCRIBE todos;` 确认表存在。

当前学习重点已根据目标 Go 后端简历调整：优先训练状态机、MySQL 事务与幂等、Redis 缓存、RabbitMQ、TraceID/pprof 和 Docker；AI 能力在传统后端基本盘稳定后再加入。

## 当前架构

```text
Gin Handler
    |
Todo Service
    |
内存 []Todo + sync.RWMutex
```

## 唯一下一步

为 `MySQLRepository.List` 编写真实 MySQL 集成测试，验证查询链路。

要求：

- 在 `internal/todo/mysql_repository_test.go` 中新增测试。
- 使用 `database.OpenMySQL()` 连接当前 Docker MySQL。
- 测试开始前插入一条带唯一标题的 Todo 测试数据。
- 调用 `NewMySQLRepository(db).List(ctx)`。
- 断言返回列表中包含刚插入的标题。
- 测试结束后删除这条测试数据，避免污染数据库。
- 数据库不可用时测试应明确失败，不要静默跳过。

写完运行 `go fmt ./...` 和 `go test ./...`，再贴出文件内容和结果。

## 跨设备与跨 Agent 续接

切换电脑或更换 AI Agent 时，不需要依赖对话记忆。新 Agent 应按这个顺序读取：

```text
1. docs/README.md
2. docs/progress.md
3. docs/roadmap.md
4. 与“唯一下一步”直接相关的 Go 文件和测试
5. 必要时运行当前验证命令
```

可直接发送给新 Agent 的提示词：

```text
我是 Go 后端学习项目的维护者。请先阅读 docs/README.md、docs/progress.md 和 docs/roadmap.md，
再读取与 progress.md 中“唯一下一步”相关的代码和测试。
先告诉我当前已完成什么、当前唯一下一步是什么；然后像导师一样只给我一个小任务。
我自己写代码，你负责讲解、审查、排错和在我给出验证结果后更新 docs/progress.md。
不要跳过基础，也不要直接替我大段完成项目代码。
```

项目路径：`D:\demo\awesomeProject`。

Git 同步状态：

- [x] 本地 Git 仓库已初始化，当前分支为 `main`。
- [x] 已创建首个本地提交：`24aabe4 chore: 初始化 Go 学习项目`。
- [x] 已配置 HTTPS 远程：`https://github.com/peach662/GO-Demo.git`。
- [x] 已推送 `main` 到 `origin/main`，本地分支已建立跟踪关系。

远程推送完成后，两台电脑只需要克隆同一个私有仓库，就能同步代码和本目录中的学习文档。


## 后续里程碑

```text
MySQL
-> Repository
-> 事务与索引
-> JWT
-> Redis
-> Docker、日志与健康检查
-> RabbitMQ、事务外盒与幂等
-> DeepSeek API、Tool Calling、RAG
-> NexusAgent 源码精读
```

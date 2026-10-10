# 学习进度

## 当前状态

当前阶段：企业级后端第 1 步，正在把内存版 Todo API 迁移到 MySQL。

当前项目：本仓库。本机路径 `D:\demo\awesomeProject`。MySQL 在服务器 `124.221.130.183:33603` 的库 `awesome_project`。

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
- [x] 为 `MySQLRepository.List` 编写真实 MySQL 集成测试，覆盖插入测试数据、查询、断言和清理。
- [x] 实现 `MySQLRepository.GetByID`：区分查询成功、无数据和数据库错误。
- [x] 为 `GetByID` 编写真实 MySQL 测试，覆盖查到数据和查不到数据。
- [x] 为 `GetByID` 编写数据库错误测试，验证连接关闭时返回 `err`。
- [x] 实现 `MySQLRepository.Create`：执行 INSERT 并获取自增 ID。
- [x] 为 `Create` 编写真实 MySQL 集成测试，验证默认状态、持久化和测试数据清理。
- [x] 实现 `MySQLRepository.UpdateStatus`：更新状态并查询返回最新 Todo。
- [x] 为 `UpdateStatus` 编写真实 MySQL 集成测试，覆盖更新成功和不存在 ID。
- [x] 将 `Service` 迁移为依赖 `Repository` 接口，Service 测试使用 `fakeRepository`。
- [x] 将 HTTP 测试迁移为使用测试 `fakeRepository`，适配新的 Service context/error 签名。
- [x] 通过真实 HTTP 验证 MySQL-backed Todo：健康检查、列表、创建和创建后查询。
- [x] 通过真实 HTTP 验证 PATCH 状态更新，并用 GET 确认 `done=true` 已持久化。
- [x] 重启 Go 服务后再次查询 ID=39，数据仍存在，确认 Todo 已由 MySQL 持久化。
- [x] 使用 `if err := r.Run(":9090"); err != nil` 显式处理 Gin 启动错误。
- [x] 增加 `config.Load`，从 `.env` 或系统环境变量读取 `MYSQL_DSN`。
- [x] 修改 `database.OpenMySQL(dsn)`，移除数据库 DSN 硬编码。
- [x] 统一 MySQL 集成测试配置加载，并处理测试工作目录与项目根目录不同的问题。
- [x] 定义 Todo 状态类型、合法流转规则和同状态幂等规则，并用表格测试验证。
- [x] 将 Todo 模型、请求 DTO、Repository、Service、Handler 和测试从 `Done bool` 迁移为 `Status Status`。
- [x] 在 Service 更新状态前调用 `CanTransition`，拒绝非法状态流转。
- [x] 修正 Repository 同状态更新的幂等行为：通过更新后查询区分“记录不存在”和“状态未变化”。
- [x] 新增 `migrations/002_add_status_to_todos.sql`，将旧 `done` 数据迁移到 `status`。
- [x] 启动 Docker Desktop 和 MySQL 容器，执行 `002_add_status_to_todos.sql`。
- [x] 验证 `status` 字段、历史数据回填和 `chk_todos_status` CHECK 约束。
- [x] 迁移后执行 `go test -count=1 ./...`，根包、配置、数据库和 Todo 集成测试全部通过。
- [x] 确认 `todo_status_logs` 不加外键：`todo_id` 是否真实存在由应用层查询和事务保证。
- [x] 在本地 Docker MySQL 执行 `migrations/003_create_todo_status_logs.sql`。
- [x] 验证 `todo_status_logs` 已创建，包含 `idx_todo_status_logs_todo_id`、两条状态 CHECK，且没有 FOREIGN KEY。
- [x] `MySQLRepository.UpdateStatus` 改为事务：`SELECT ... FOR UPDATE` 读取当前行，状态变化时更新 `todos.status` 并插入 `todo_status_logs`，失败由 `Rollback` 撤销。
- [x] 同状态直接返回，不写审计；不存在的 Todo 返回未找到，不写审计。
- [x] `TestMySQLRepositoryUpdateStatus` 断言 `PENDING -> PROCESSING` 后审计恰好 1 行，同状态再更新仍是这 1 行，清理时先删审计再删 Todo。
- [x] `EXPLAIN` 按 `todo_id` 查询 `todo_status_logs` 时，实际使用索引 `idx_todo_status_logs_todo_id`，`type` 为 `ref`。
- [x] 测试里的 `COUNT(*)`、`MIN(from_status)`、`MIN(to_status)` 聚合查询同样走这个索引，`Extra` 为空，因为状态列不在索引里，需要回表。
- [x] `SHOW INDEX` 确认 `idx_todo_status_logs_todo_id` 只有 `todo_id` 一列，并且允许同一个 `todo_id` 有多行审计。
- [x] 新增 `migrations/004_create_users.sql`，并在服务器 `awesome_project` 库创建 `users` 表。`username` 有唯一索引 `uk_users_username`，密码字段是 `password_hash`。
- [x] 新增 `internal/user/model.go`。`User` 的 JSON 只包含 `id` 和 `username`，`PasswordHash` 使用 `json:"-"`，不会出现在响应里。
- [x] `TestUser` 用 `json.Marshal` 验证响应包含 `alice`，且不包含 `secret-hash` 和 `password_hash`。`go test ./internal/user/` 通过。
- [x] 定义 `user.Repository` 接口：`Create` 用于注册，`GetByUsername` 用 `bool` 表示用户是否存在。
- [x] 定义 `user.MySQLRepository` 和 `NewMySQLRepository`，保存 `*sql.DB` 连接池。
- [x] 实现 `GetByUsername`：查到返回用户，`sql.ErrNoRows` 返回未找到，其他错误原样返回。
- [x] `TestMySQLRepositoryGetByUsername` 用真实 MySQL 覆盖查到用户和查不到用户。`go test ./internal/user/` 通过。
- [x] 实现 `Create`：插入 `username` 和 `password_hash`，用 `LastInsertId` 填入返回的 `User.ID`。
- [x] `TestMySQLRepositoryCreate` 调用 `repo.Create` 创建用户，断言 `ID > 0`、用户名和密码哈希一致，再用 `GetByUsername` 确认库中记录与返回的 `ID` 相同。测试结束后按用户名删除。
- [x] 同一测试里第二次用相同用户名 `Create`，确认返回 MySQL 错误号 `1062`（唯一索引 `uk_users_username`）。
- [x] 定义 `user.ErrUsernameTaken`。`Create` 在插入遇到 MySQL `1062` 时返回它，其他数据库错误原样返回。
- [x] 项目级安装 TuTor（`.agents/skills/tutor` 与 `.cursor/skills/tutor`）和 cc-skills-golang（`.cursor/skills/golang-*`）。`npx skills add --all` 会往很多 Agent 目录复制，已删掉 `.claude` 和 `agent` 里的重复副本。
- [x] 定义 `user.ErrInvalidUsername`。
- [x] 新增 `user.Service` 和 `Register`：空用户名直接失败；bcrypt 哈希明文密码后再调用 `repo.Create`。
- [x] `TestServiceRegister` 使用假仓库：成功注册存的是哈希、空用户名返回 `ErrInvalidUsername`、重复用户名返回 `ErrUsernameTaken`。`go test ./internal/user/` 通过。
- [x] `Register` 拒绝空密码，返回 `ErrInvalidPassword`。测试覆盖空密码。`go test ./internal/user/` 通过。
- [x] 新增 `user.RegisterRequest`（`username`、`password`，`binding:"required"`）。
- [x] `main` 创建 User Repository/Service，`newRouter` 接收两个 Service。
- [x] 增加 `POST /users/register`：绑定 JSON；`ErrUsernameTaken` / `ErrInvalidPassword` / `ErrInvalidUsername` 返回 400；其它错误 500；成功返回用户（不含密码哈希）。
- [x] `main_test.go` 中所有 `newRouter` 调用补上 `user.NewService(nil)`。`go test ./...` 通过。
- [x] 新增 `fakeUserRepository`、`newTestUserService` 和 `TestRegisterRoute`：`POST /users/register` 断言 200、响应含用户名、不含明文密码和 `password_hash`。`go test ./...` 通过。
- [x] `fakeUserRepository.Create` 查重后返回 `ErrUsernameTaken`。新增 `TestRegisterDuplicateRoute`：预置同名用户再注册，断言 HTTP 400。`go test . -run TestRegisterDuplicateRoute` 通过。
- [x] 新增 `ErrInvalidCredentials` 和 `Service.Login`：不存在用户与密码错误都返回同一错误；用 `bcrypt.CompareHashAndPassword` 校验。
- [x] `TestServiceLogin` 覆盖成功登录、错误密码、不存在用户。`go test ./...` 通过。
- [x] 新增 `user.LoginRequest` 与 `POST /users/login`：`ErrInvalidCredentials` → 401；空用户名/空密码 → 400；成功返回用户。`go test ./...` 通过。
- [x] 新增 `TestLoginRoute`：先注册再登录断言 200 与敏感字段不泄露；错密码断言 401。`go test . -run TestLoginRoute` 通过。
- [x] `config.Load` 读取 `JWT_SECRET`（空则报错）；`.env.example` 增加占位。
- [x] `main` 使用 `auth.NewJWT`；`newRouter` 注入 `*auth.JWT`；登录成功响应 `data` 含 `user` + `token`。
- [x] HTTP 测试补上 JWT；`TestLoginRoute` 断言成功响应含 `token`。`go test ./...` 通过。
- [x] `ParseToken`：验签、HMAC 检查、`user_id`（float64→int）、空密钥拒绝。
- [x] `AuthMiddleware`：Bearer 校验；失败统一 `code/message/data` 401；成功写入 `ContextUserIDKey`。
- [x] `GET /users/me` 挂中间件，返回当前 `user_id`。登录/注册保持公开。`go test ./...` 通过。
- [x] 新增 `TestMeRoute`：注册→登录解析 `data.token`→带 Bearer 访问 `/users/me` 断言 200/`user_id`；无 Header 断言 401。`go test . -run TestMeRoute` 通过。
- [x] migration `005_add_user_id_to_todos.sql`：`todos.user_id BIGINT NULL` + 索引；已在服务器 MySQL 执行。
- [x] `todo.Todo` 增加 `UserID *int`（JSON `user_id`）；`Create` 链路（repo/service/fake/测试）接收 `userID int`。
- [x] `scanTodo` 用 `sql.NullInt64` 扫 `user_id`，避免历史 NULL 行 Scan 失败。
- [x] `POST /todos` 挂 `AuthMiddleware`；从 Context 取 `userID` 并 `.(int)`；`Create` 使用真实 userID（不再写死 `0`）。
- [x] HTTP 测试：`TestCreateTodoRoute` / `TestCreateTodoValidationError` 带 Bearer（`GenerateToken`）；新增 `TestCreateTodoUnauthorized`（无 Header → 401，列表不变）。`go test . -run 'TestCreateTodo'` 通过。
- [x] `List` / `GetByID` / `UpdateStatus` 全链路带 `userID`：MySQL `WHERE user_id = ?`；Service 透传；HTTP 的 GET/PATCH 挂 `AuthMiddleware`。
- [x] HTTP 测试：fake 按 `UserID` 过滤；GET/PATCH 带 Bearer；列表不含他人 Todo；他人 id 当 404。`go test ./internal/todo/ -count=1` 与 `go test . -count=1` 通过。
- [x] migration `006_make_todos_user_id_not_null.sql`：`UPDATE` 把 `user_id IS NULL` 回填为用户 `1`，再 `MODIFY ... NOT NULL`；已在服务器 `mysql-33603` / `awesome_project` 执行。验证：`null_user_id=0`，列 `Null=NO`。
- [x] 服务器 Docker 新建 `redis-36379`（`redis:7-alpine`，`36379->6379`，volume `redis-36379-data`，AOF）；容器内 `PONG`。
- [x] `REDIS_ADDR` 写入 `.env.example` / `.env`；`config.Load` 读取 `RedisAddr`。
- [x] `database.OpenRedis`：`go-redis/v9` + `Ping`；`main` 启动时连接并 `defer Close`。`go run .` 监听 `:9090` 无 panic；`go test ./...` 通过。
- [x] `GetByID` 读缓存：key `todo:{userID}:{id}`，TTL 5m；未命中查库并回填；`redis==nil` 仅查库。`NewService(repo, redis)`；测试传 `nil`。`go test ./...` 通过。
- [x] `UpdateStatus` 写库成功后 `Del` 同 key（`redis==nil` 跳过）；`go test ./...` 通过。
- [x] `List` 缓存：key `todos:{userID}`，TTL 2m；Create / UpdateStatus 成功后失效列表 key（Update 同时删详情 key）。`go test ./...` 通过。
- [x] `GET /health`：注入 redis；`Ping` 失败返回 503；测试 `newRouter(..., nil)`。`go test .` 通过。
- [x] 新增 `docs/middleware.md`：说明 `mysql-33603` / `redis-36379`、端口与 `.env`；`docs/README.md` 增加入口链接。
- [x] 服务器 Docker 启动 `rabbitmq-35672`（`rabbitmq:3-management-alpine`）：AMQP `35672→5672`，管理台 `35673→15672`，用户 `app` / `app123`，volume `rabbitmq-35672-data`；`rabbitmq-diagnostics ping` 成功。`docs/middleware.md` / `.env.example` 已记 `RABBITMQ_URL`。
- [x] `config.Load` 读取 `RabbitMQURL`（`RABBITMQ_URL`）；`database.OpenRabbitMQ` 用 `amqp091-go` Dial；`main` 启动连接并 `defer conn.Close()`。`go run .` 监听 `:9090`，无 RabbitMQ panic。
- [x] `database.DeclareQueue`：Channel + `QueueDeclare`（durable）；`main` 声明 `todo.events`；管理台可见该队列（Features `D`，running）。
- [x] `database.PublishMessage`：默认交换机 `""` + 路由键队列名；`main` 启动时发 `hello`。`go run .` 监听 `:9090`，无 panic。管理台 Get messages：payload `hello`，`content_type: text/plain`，Exchange `(AMQP default)`，Routing Key `todo.events`。
- [x] `database.ConsumeMessage`：`Consume` + `<-msgs` + 手动 `Ack`；`main` 在 Publish 后取一条并打印。`go run .` 启动日志出现 `hello` 后监听 `:9090`。
- [x] 从 `main` 启动路径移除学习用 `PublishMessage` / `ConsumeMessage` / `fmt.Println`；保留 `OpenRabbitMQ`、`defer Close`、`DeclareQueue`。`Publish`/`Consume` 函数仍留在 `database` 包供下一步业务接入。
- [x] `todo.Service` 注入 `*amqp.Connection`（`nil` 跳过）；`Create` 成功后可选清 Redis 列表缓存，再 `json.Marshal` + `database.PublishMessage` 发到 `todo.events`。`main` 传 `conn`；测试传 `nil`。HTTP 创建 Todo 后管理台/API 可见 JSON（如 `id`/`title`/`status`/`user_id`）。
- [x] `database.StartConsumer`：goroutine 中 `range` + `Ack`；`main` 声明队列后启动。创建 Todo（如 id=82）后队列 Ready/Unacked 均为 0；进程应打印 `Received message: ...`。
- [x] 消费循环：JSON 成功则打印并 `Ack`；失败则 `Nack(false, false)` + `continue`。验证：合法 Todo（id=83）被消费；非法 `not-json` 入队后 Ready/Unacked 仍为 0（丢弃不重试）。
- [x] 进程内 `map[int]struct{}` 幂等：同一 `id` 第二次只打 Duplicate 并 Ack。验证：连发两次 `{"id":99,"title":"dup"}`，第一次 `Todo: {ID:99 Title:dup}`，第二次 `Duplicate message`；队列 Ready/Unacked 为 0。
- [x] Redis 幂等：`StartConsumer` 注入 `*redis.Client`；key `todo:event:{id}`，`Exists` 判重，`Set` TTL 24h；nil 时回退内存 `seen`。验证：id=1001 连发两次，先 Todo 后 Duplicate；Redis `GET todo:event:1001` = `1`，TTL≈24h。
- [x] `internal/trace`：中间件读/生成 `X-Request-ID`，写入响应头、`c.Set` 与 request context；`newRouter` 中 `r.Use(trace.Middleware())`。`curl.exe -i` 验证自带头与自动生成。
- [x] `trace.IDFromContext` + `/health` 使用 `slog.Info(..., "trace_id", ...)`。验证：请求头 `X-Request-ID: my-test-id` 时日志出现 `INFO health ok trace_id=my-test-id`。
- [x] `POST /todos`：入口 `slog.Info`；绑定/类型/Create 失败路径在 `c.JSON` 前 `slog.Error`（含 `trace_id`、`error`）；成功再 `slog.Info`。`go build .` 通过。
- [x] `PATCH /todos/:id`：同样接 TraceID + slog；非法流转/404 用 Info，500 用 Error；成功只打 `todo_id`。`go build .` 通过。
- [x] 多阶段 `Dockerfile`（builder `golang:1.25-alpine` + runtime `alpine`）；CN 构建加 `GOPROXY=https://goproxy.cn,direct`。已推送 `249c6b9`。
- [ ] 服务器 `docker build`：首次因 `proxy.golang.org` 超时失败；加 GOPROXY 后构建超时；随后 SSH banner 超时（主机可 ping，疑似构建占满资源）。部署验证待服务器恢复后继续。

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

用户已执行 `go fmt ./...` 和 `go test ./...`，`MySQLRepository.List` 的真实 MySQL 集成测试通过；测试包含测试数据清理。

用户已执行无缓存定向测试：

```powershell
go test -count=1 ./internal/todo -run 'TestMySQLRepository(GetByID|List)'
```

结果通过，覆盖 List 和 GetByID 的成功、未找到、数据库错误场景。

用户已执行 `go fmt ./...` 和 `go test ./...`，`MySQLRepository.Create` 已通过编译检查；目前还没有 Create 的真实数据库集成测试。

用户已执行 `go test ./...`，Create 集成测试通过；测试验证返回 ID、标题、默认 `Done=false`，并通过 `GetByID` 验证持久化。

用户已执行 `go fmt ./...` 和 `go test ./...`，`UpdateStatus` 已通过编译检查；目前还没有 UpdateStatus 的真实数据库集成测试。

用户已执行 `go test ./...`，UpdateStatus 集成测试通过；测试验证更新结果、数据库持久化和不存在 ID。

用户已执行 `go test ./internal/todo`，Service 的 Repository 接口迁移测试通过。

用户已执行无缓存完整测试：

```powershell
go test -count=1 ./...
```

根包、`internal/database` 和 `internal/todo` 均通过。

配置迁移后已执行无缓存完整测试：

```powershell
go test -count=1 ./...
```

根包、`internal/config`、`internal/database` 和 `internal/todo` 均通过；本地 `.env` 未提交。

用户已手动验证：

```text
GET /health -> code=0, data=pong
GET /todos -> 初始空列表
POST /todos -> 创建 MySQL Todo，返回 id=39
GET /todos -> 能查到 id=39 的 MySQL Todo
PATCH /todos/39 -> done=true
GET /todos/39 -> done=true
重启服务后 GET /todos/39 -> 仍为 done=true
```

## 已知依赖说明

Gin `v1.12.0` 依赖 `github.com/goccy/go-yaml v1.19.2`，但当前下载到的该版本缺少 Gin 运行时导入的包，导致 `go mod tidy` 和完整 race 测试失败。

`go.mod` 已通过 `replace` 临时使用兼容的 `github.com/goccy/go-yaml v1.18.0`，并已验证：

```powershell
go test -race -count=1 ./...
```

通过。不要删除这条 `replace`，除非升级 Gin 或确认上游 YAML 发布问题已修复，并重新运行完整测试。

Docker 环境也已确认可用：Docker Desktop 4.55.0，Docker Engine 29.1.3。

事务版 `UpdateStatus` 已写入代码，但测试还没有检查审计表。已执行：

```powershell
go test -count=1 ./internal/todo -run TestMySQLRepositoryUpdateStatus
```

结果通过。这个测试当时只验证了状态更新和不存在的 ID。

2026-09-25 已把连接改到服务器 MySQL，并执行：

```powershell
go test -count=1 ./...
```

根包和 `internal/todo` 通过。`TestMySQLRepositoryUpdateStatus` 覆盖了审计行数和同状态不新增审计。

用户在服务器 MySQL 执行：

```sql
EXPLAIN SELECT id, from_status, to_status FROM todo_status_logs WHERE todo_id = 1;
```

`key` 为 `idx_todo_status_logs_todo_id`，`type` 为 `ref`，`key_len` 为 8，`rows` 为 1。

同一张表上的聚合查询：

```sql
EXPLAIN SELECT COUNT(*), MIN(from_status), MIN(to_status) FROM todo_status_logs WHERE todo_id = 1;
```

结果相同：`key` 仍是 `idx_todo_status_logs_todo_id`，`type` 为 `ref`，`Extra` 为 `NULL`。`WHERE todo_id = 1` 先用索引定位，再回表计算 `COUNT` 和 `MIN`。

`SHOW INDEX FROM todo_status_logs` 有两行：

- `PRIMARY` 的 `Column_name` 是 `id`，`Non_unique` 为 0，主键不能重复。
- `idx_todo_status_logs_todo_id` 的 `Column_name` 只有 `todo_id`，`Non_unique` 为 1，同一个任务可以有多条审计。`Cardinality` 为 0，是因为当前表里没有数据。

用户已确认：

```powershell
go test ./internal/user/
```

结果：`ok awesomeProject/internal/user`。`TestMySQLRepositoryCreate` 覆盖了成功创建和重复用户名 `1062`。

用户已确认 `Create` 将 `1062` 转成 `ErrUsernameTaken` 后：

```powershell
go test ./internal/user/
```

结果再次通过。

用户已确认 `Register` 假仓库测试：

```powershell
go test ./internal/user/
```

结果：`ok awesomeProject/internal/user`（约 2.8s）。覆盖成功注册、空用户名和重复用户名。

空密码校验后再次确认：

```powershell
go test ./internal/user/
```

结果：`ok awesomeProject/internal/user`。

用户已确认注册 HTTP 接入后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。

用户已确认注册 HTTP 测试后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。`TestRegisterRoute` 覆盖成功注册与敏感字段不泄露。

用户已确认重复注册 HTTP 测试后：

```powershell
go test . -run TestRegisterDuplicateRoute
```

结果：`ok awesomeProject`。覆盖 `ErrUsernameTaken` 映射为 400。

用户已确认 `Login` 假仓库测试后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。`TestServiceLogin` 覆盖成功、错密码、不存在用户。

用户已确认登录 HTTP 接入后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。

用户已确认登录 HTTP 测试后：

```powershell
go test . -run TestLoginRoute
```

结果：`ok awesomeProject`。覆盖成功登录与错密码 401。

用户已确认 JWT 接入登录后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。登录成功响应含 `token`。

用户已确认 JWT 中间件与 `GET /users/me` 后：

```powershell
go test ./...
```

结果：根包、`internal/todo`、`internal/user` 均通过。

用户已确认 `TestMeRoute` 后：

```powershell
go test . -run TestMeRoute
```

结果：`ok awesomeProject`。覆盖带 token 的 200 与无 Authorization 的 401。

本机 `.env` 曾指向 `127.0.0.1:13306`，测试报 `users` 表不存在。已把本地 `MYSQL_DSN` 改到服务器 `124.221.130.183:33603`。`.env` 仍不提交。

已在容器 `awesome-project-mysql` 的 `awesome_project` 库执行 `003_create_todo_status_logs.sql`。`SHOW CREATE TABLE todo_status_logs` 确认：

- 字段：`id`、`todo_id`、`from_status`、`to_status`、`created_at`
- 索引：`PRIMARY (id)`、`idx_todo_status_logs_todo_id (todo_id)`
- CHECK：`chk_todo_status_logs_from_status`、`chk_todo_status_logs_to_status`
- 无 FOREIGN KEY

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
Todo Service / User Service
    |
Todo Repository / User Repository
    |
MySQL（服务器 124.221.130.183:33603）

另：internal/auth（签发 + 验签 + 中间件；`GET /users/me` 已验证接线）
```

## 唯一下一步

Dockerfile 已入库。下一步：服务器恢复 SSH 后完成 `docker build` / `docker run`（`--network host` + 本机中间件地址）并用 `/health` 验证；或改为本机交叉编译二进制再上传运行，避开服务器内 `go mod download`。

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

项目路径：`E:\projects\GO-Demo\GO-Demo`。换电脑时先 `git pull`，SSH 私钥放在 `~/.ssh/id_ed25519`，并按 `.cursor/skills/server-middleware/SKILL.md` 配置主机别名 `middleware`。

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

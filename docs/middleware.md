# 服务器中间件（本项目）

本机一般不跑 Docker Desktop。MySQL 和 Redis 跑在腾讯云服务器上，开发机通过公网端口连接。

服务器地址：`124.221.130.183`（SSH 用户多为 `ubuntu`，详见仓库内 Cursor skill `server-middleware`）。

## 已有容器

| 容器名 | 镜像（约） | 宿主机端口 | 用途 |
|--------|------------|------------|------|
| `mysql-33603` | `mysql:8.0` | `33603` → 容器 `3306` | 业务库 `awesome_project` |
| `redis-36379` | `redis:7-alpine` | `36379` → 容器 `6379` | Todo 列表/详情缓存；`/health` Ping |
| `rabbitmq-35672` | `rabbitmq:3-management-alpine` | `35672` → `5672`（AMQP）；`35673` → `15672`（管理台） | 消息队列（学习连通中） |

在服务器上查看：

```bash
sudo docker ps --filter name=mysql-33603 --filter name=redis-36379 --filter name=rabbitmq-35672
```

RabbitMQ 默认学习账号（容器环境变量）：用户 `app`，密码与 MySQL 学习账号同风格见 `.env.example`（`RABBITMQ_*`）。管理界面：浏览器打开 `http://124.221.130.183:35673`（若安全组未放行 35673，需在云控制台开端口或仅用 AMQP）。
## 本机 `.env` 如何指向它们

仓库跟踪 `.env.example`；本机可复制为 `.env`（已 gitignore）。`config.Load` 先读 `.env`，再回退 `.env.example`。

和本项目相关的项：

| 变量 | 含义 | 示例形态 |
|------|------|----------|
| `MYSQL_DSN` | Go `database/sql` 连接串 | `用户:密码@tcp(124.221.130.183:33603)/awesome_project?...` |
| `REDIS_ADDR` | go-redis 的 `Addr` | `124.221.130.183:36379` |
| `RABBITMQ_URL` | AMQP 连接串（`amqp` 驱动） | `amqp://app:app123@124.221.130.183:35672/` |
| `JWT_SECRET` | 签发/校验 JWT | 本地自拟，勿提交真实生产密钥 |

完整占位见仓库根目录 `.env.example`。换电脑后：`git pull` → 需要时复制 `.env.example` 为 `.env` → `go run .` 或 `go test ./...`。

## 不要做的事

- **不要**对 `mysql-33603` / `redis-36379` / `rabbitmq-35672` 执行 `docker rm -v`、随意换 volume 名或重建容器，除非你明确要清空学习数据。
- **不要**把本机 Docker Desktop 再起一套同端口库，和文档、测试假设冲突。
- **不要**把真实生产密码写进 Git；学习占位只留在 `.env.example` / 本地 `.env`。

## 健康检查

应用起来后：

```text
GET http://localhost:9090/health
```

会 Ping Redis：失败时 HTTP `503`。MySQL 在进程启动时已 `Ping`，连不上会直接启动失败。

## 和 Cursor skill 的关系

更完整的 SSH / 新开中间件步骤见 `.cursor/skills/server-middleware/SKILL.md`。本文只固定 **本仓库当前在用的三个容器**，方便换电脑续学。

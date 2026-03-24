# alioth-hrc — Claude 工作指引

面向 AI 辅助开发与代码审查的**项目级约定**。风格与结构对齐 Go 官方与社区共识：[`Effective Go`](https://go.dev/doc/effective_go)、[`Code Review Comments`](https://go.dev/wiki/CodeReviewComments)、[`Organizing a Go module`](https://go.dev/doc/modules/layout)、[`golang-standards/project-layout`](https://github.com/golang-standards/project-layout)（非官方标准，但广泛参考）。

---

## 项目是什么

- **定位**：HTTP API 服务（当前为健康检查与 MySQL/Redis 连通性演示）。
- **模块名**：`alioth-hrc`（见 `go.mod`）。
- **运行时**：`cmd/server` 启动；`internal/app` 组装配置、数据库、缓存、路由与 `http.Server`；`main` 负责信号与优雅关闭。

---

## 常用命令

在项目根目录执行：

```bash
go mod tidy
go vet ./...
go build -o bin/server ./cmd/server
go run ./cmd/server
```

当前仓库**尚无** `_test.go` 与 Makefile；新增测试后使用 `go test ./...`。

---

## 目录与职责（本仓库实际结构）

| 路径 | 职责 |
|------|------|
| `cmd/server` | 可执行入口：`main`、信号处理、调用 `internal/app`。保持瘦入口。 |
| `internal/app` | 应用生命周期：加载配置、初始化 MySQL/GORM、Redis、注册路由、启动/关闭 HTTP。 |
| `internal/config` | 配置加载与校验（Viper + `.env`，环境变量可覆盖）。 |
| `internal/db` | MySQL：GORM 与底层 `*sql.DB` 连接池。 |
| `internal/cache` | Redis 客户端与连通性检查。 |
| `internal/router` | Gin 引擎与路由注册。 |
| `internal/handler` | HTTP 处理器（依赖注入 `*sql.DB`、`redis.Client`）。 |
| `internal/model` | 数据模型（如 GORM struct）。 |

**约定**：可对外复用的库若未来出现，再考虑 `pkg/`；仅本应用使用的代码放在 `internal/`。未使用顶层 `src/`（易与 GOPATH 时代习惯混淆）。

---

## 配置与环境

- 复制 `.env.example` 为 `.env`，按需填写 `MYSQL_DSN`、`REDIS_*` 等。
- 必填项逻辑见 `internal/config` 的 `Validate()`。
- **Gin 模式**：`internal/router` 根据 `APP_ENV` 设置——`prod`/`production` 为 `ReleaseMode`，`test` 为 `TestMode`，其余为 `DebugMode`。

---

## 编码约定（本仓库）

- **错误**：向上返回时尽量 `fmt.Errorf("...: %w", err)`；HTTP 层避免把内部错误细节暴露给不可信客户端（健康检查接口可按现有风格适度暴露，业务 API 需收敛）。
- **Context**：HTTP 与 DB/Redis 调用使用 `c.Request.Context()` 或带超时的子 context（与现有 `health`/`ready` 一致）。
- **并发**：`http.Server` 已配置 `ReadHeaderTimeout`；扩展时继续为 Server 与客户端设置合理超时。
- **依赖**：通过构造函数注入（如 `NewHealthHandler`），避免未必要的全局单例。
- **格式与静态分析**：提交前 `gofmt` / `go vet`；团队可统一引入 `staticcheck` 等（与社区推荐一致）。

---

## 与助手协作时注意

- **优先小步修改**：匹配现有包名、错误包装与日志（`log/slog`）风格。
- **不要**为“可能复用”提前抽 `pkg` 或复杂分层；随业务增长再拆。
- **数据库**：启动时对 `model.User` 执行 GORM `AutoMigrate`；生产若改用独立迁移工具（如 goose/atlas），可改为在 `app` 中按环境开关或移除自动迁移。

---

## 架构简评（供维护参考）

当前结构（`cmd` + `internal` + 配置/DB/缓存/路由/处理器分层）对中小型 API **合理、常见**，边界清晰。已落实：Gin 按 `APP_ENV` 设模式、`go mod tidy` 校正直接依赖、GORM 启动迁移、关闭时释放 MySQL/Redis 连接。待改进点：补 `_test.go`、业务增长时可引入 service/repository 分层（仍保持 YAGNI）。

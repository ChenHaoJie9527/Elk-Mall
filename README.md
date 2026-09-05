# Elk-Mall

Go + Echo 商城后端，按版本迭代基础设施。当前已发布 `v1.0.0` / `v2.0.0` / `v3.0.0` / `v4.0.0` / `v5.0.0`。

## 已完成功能

### v1.0.0 — 服务可启动、可探活

- 项目基础目录与模块拆分（`cmd` / `internal`）
- 用 Viper 加载根目录 `config.yaml`，支持 `SERVER_PORT` 环境变量覆盖，默认端口 `8080`
- 统一响应结构 `Response{code, msg, data}`，成功走 `response.Success`
- 注册路由，提供 `GET /ping` 健康检查，返回 `{"code":0,"msg":"success","data":"pong"}`
- 入口读配置、创建 Echo、挂载路由并启动 HTTP 服务

### v2.0.0 — 业务错误码、统一错误响应、Recover 防崩

- 基于 `error` 标准库自定义业务错误码 `Errno`（`OK` / `ParmError` / `NotFound` / `ServerError`）
- `WithMsg` 只改文案、保持原 Code；`Err` 将 `Errno` 映射为统一 JSON
- 全局 `HTTPErrorHandler`：
  1. 响应已提交则不再覆盖
  2. 业务错误走 `Errno`，HTTP 200 + JSON
  3. 其余按 HTTP 状态码返回
- 挂载 `Recover` 中间件：panic 不会打崩进程，对外返回 500，不泄漏 panic 内容
- 覆盖测试：`Errno` 映射、错误处理分流、Recover 防崩

### v3.0.0 — 请求 ID 与结构化日志

- 挂载 Request ID 中间件，用 UUID v7 生成请求 ID
- 用 `slog` JSON 输出请求日志，记录 `request_id` / `uri` / `method` / `status`
- 成功请求打 `REQUEST`，出错打 `REQUEST_ERROR`（含 error 信息）

### v4.0.0 — 数据接入：MySQL / Redis 连接池与探活

- `docker-compose.yml` 提供本地依赖：MySQL 8 映射 `3308:3306`，Redis 8 映射 `6380:6379`（避开宿主机已占用的 3306 / 6379）
- `config.yaml` 拆字段描述连接（host / port / user / password 等），adaptor 再拼 DSN；带下划线的池参数用 `mapstructure` tag 才能解进结构体
- `internal/adaptor` 用 `database/sql` + 官方 MySQL 驱动、go-redis 建连接池；启动时 Ping，失败则进程退出、不对外听端口
- `GET /ping` 注入连接后探测两个依赖，成功返回 `{"mysql":"ok","redis":"ok"}`；不通则走全局错误处理变成 500
- 本版不建表、不上 GORM、不写 repository；adaptor 只负责连上和探活

### v5.0.0 — 数据模型与仓储：UserDO / UserRepo

- `internal/model/do` 增加 `UserDO`：对照用户表写列和约束（主键、用户名唯一索引、密码/昵称长度、`CreatedAt` / `UpdatedAt`、软删除 `DeletedAt`）
- DO 只给 GORM 落库用，不加 `json` tag；对外 DTO 本版仍未加
- `internal/repository` 增加 `UserRepo`，持有 `*gorm.DB`，封装用户存取：`CreateUser` / `GetByID` / `GetList`（`Count` 总条数 + `Offset` / `Limit` 分页）
- 引入 `gorm.io/gorm`；插入、按主键查、列表都走 GORM 链式调用，错误从 `.Error` 取出
- 本版不挂用户 HTTP 接口、不改 `main` 接入 GORM、不 `AutoMigrate`；adaptor 仍只负责连上 `database/sql`

本地依赖：

```bash
docker compose up -d
```

应用仍在宿主机运行，连 `127.0.0.1:3308` 与 `127.0.0.1:6380`。Navicat 用同一套账号（`elk` / `elk`，端口 3308）**双击打开连接**。`docker compose down` 只拆容器、数据还在；`down -v` 会连数据卷一起删掉。

## 目录结构

```
elk-mall/
├── cmd/server/main.go          # 入口：读配置、连 MySQL/Redis、建 Echo、挂中间件、注册路由、启动
├── internal/
│   ├── adaptor/                # 数据接入：连接池 + Ping（v4）
│   ├── config/                 # 配置加载
│   ├── router/                 # 路由层
│   ├── controller/             # 控制层（当前仅 health）
│   ├── service/                # 服务层（后续模块）
│   ├── repository/             # 仓储：UserRepo 用户 CRUD（v5）
│   ├── model/do/               # DO：UserDO 表映射（v5）；DTO 后续再加
│   ├── middleware/             # 自定义中间件（后续）
│   └── common/
│       ├── Errno/              # 业务错误码（v2）
│       └── response/           # 统一响应 + 全局错误处理
├── docker-compose.yml          # 本地 MySQL / Redis
├── config.yaml
└── go.mod
```

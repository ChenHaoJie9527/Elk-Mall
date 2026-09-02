# Elk-Mall

Go + Echo 商城后端，按版本迭代基础设施。当前已发布 `v1.0.0` / `v2.0.0` / `v3.0.0`。

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

## 目录结构

```
elk-mall/
├── cmd/server/main.go          # 入口：读配置、建 Echo、挂中间件、注册路由、启动
├── internal/
│   ├── config/                 # 配置加载
│   ├── router/                 # 路由层
│   ├── controller/             # 控制层（当前仅 health）
│   ├── service/                # 服务层（后续模块）
│   ├── repository/             # 数据层（后续）
│   ├── model/                  # DTO / DO（后续）
│   ├── middleware/             # 自定义中间件（后续）
│   └── common/
│       ├── Errno/              # 业务错误码（v2）
│       └── response/           # 统一响应 + 全局错误处理
├── config.yaml
└── go.mod
```

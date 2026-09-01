elk-mall/
├── cmd/server/main.go          # 入口：读配置、建 Echo、注册路由、启动
├── internal/                   # 业务代码（不对外暴露）
│   ├── config/                 # 配置加载
│   ├── router/                 # 路由层
│   ├── controller/             # 控制层（先放 health，后续分 admin/user）
│   ├── service/                # 服务层（v1 空，后续放模块）
│   ├── repository/             # 数据层（v1 空，v4+ 用）
│   ├── model/                  # DTO / DO（v1 空）
│   ├── middleware/             # 中间件（v2/v3 开始加）
│   └── common/
│       └── response/           # 统一响应结构（v1 就要用）
├── config.yaml
└── go.mod

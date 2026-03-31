# Kite Admin

基于 Go + Vue3 的后台管理系统

## 技术栈

### 后端
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **认证**: JWT (golang-jwt/jwt)
- **密码加密**: bcrypt

### 前端
- **框架**: Vue 3
- **UI 库**: Naive UI
- **构建工具**: Vite
- **状态管理**: Pinia
- **路由**: Vue Router

## 项目结构

```
kite-admin/
├── backend/              # Go 后端
│   ├── config/          # 配置
│   ├── controllers/     # 控制器
│   ├── middleware/      # 中间件
│   ├── models/          # 数据模型
│   ├── routes/          # 路由
│   ├── utils/           # 工具函数
│   └── main.go          # 入口文件
├── frontend/            # Vue3 前端
│   ├── src/
│   │   ├── api/        # API 接口
│   │   ├── components/ # 组件
│   │   ├── layouts/    # 布局
│   │   ├── router/     # 路由
│   │   ├── store/      # 状态管理
│   │   ├── views/      # 页面
│   │   └── main.js     # 入口文件
│   └── package.json
└── README.md
```

## 快速开始

### 后端启动

1. 安装依赖
```bash
cd backend
go mod tidy
```

2. 配置数据库

创建 MySQL 数据库：
```sql
CREATE DATABASE kite_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改 `backend/.env` 配置文件：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=kite_admin
```

3. 运行服务
```bash
go run main.go
```

后端服务将在 `http://localhost:8090` 启动

### 前端启动

1. 安装依赖
```bash
cd frontend
pnpm install
```

2. 启动开发服务器
```bash
pnpm run dev
```

前端服务将在 `http://localhost:5173` 启动

## 默认账号

- 用户名: `admin`
- 密码: `123456`

## 功能特性

- ✅ 用户管理
- ✅ 角色管理
- ✅ 权限管理（菜单权限）
- ✅ JWT 认证
- ✅ 角色切换
- ✅ 个人资料
- ✅ 响应式布局
- ✅ 主题切换（亮色/暗色）

## API 接口

### 认证
- `POST /auth/login` - 登录
- `GET /auth/captcha` - 获取验证码
- `POST /auth/logout` - 退出登录
- `POST /auth/current-role/switch/:roleCode` - 切换角色

### 用户管理
- `GET /user/detail` - 获取当前用户详情
- `GET /user` - 获取用户列表（分页）
- `POST /user` - 新增用户
- `PATCH /user/:id` - 修改用户
- `DELETE /user/:id` - 删除用户
- `PATCH /user/password/reset/:id` - 重置密码

### 角色管理
- `GET /role/page` - 获取角色列表（分页）
- `GET /role` - 获取所有角色
- `POST /role` - 新增角色
- `PATCH /role/:id` - 修改角色
- `DELETE /role/:id` - 删除角色

### 权限管理
- `GET /role/permissions/tree` - 获取角色权限树
- `GET /permission/menu/tree` - 获取菜单权限树
- `GET /permission/tree` - 获取所有权限树
- `POST /permission` - 新增权限
- `PATCH /permission/:id` - 修改权限
- `DELETE /permission/:id` - 删除权限

## 开发说明

- 后端使用环境变量配置，修改 `.env` 文件后需重启服务
- 前端使用代理访问后端 API，配置在 `vite.config.js`
- 数据库表会在首次启动时自动创建
- 默认会创建管理员账号和基础权限数据

## License

MIT

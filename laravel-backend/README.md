# Laravel Admin API

基于 Laravel 11 实现的后台管理系统 API

## 功能模块

- 用户认证（JWT）
- 用户管理
- 角色管理
- 权限管理
- 系统日志

## 系统要求

- PHP >= 8.2
- Composer
- MySQL >= 5.7

## 安装步骤

1. 安装依赖
```bash
composer install
```

2. 配置环境变量
```bash
cp .env.example .env
php artisan key:generate
php artisan jwt:secret
```

3. 配置数据库并运行迁移
```bash
php artisan migrate
```

4. 初始化数据
```bash
php artisan db:seed
```

5. 启动服务
```bash
php artisan serve --port=8090
```

## API 文档

详见 `接口文档.md`

## 技术栈

- Laravel 11
- JWT Authentication (tymon/jwt-auth)
- MySQL

## 默认账号

- 用户名：admin
- 密码：123456

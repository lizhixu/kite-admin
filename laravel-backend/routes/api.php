<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use App\Http\Controllers\AuthController;
use App\Http\Controllers\UserController;
use App\Http\Controllers\RoleController;
use App\Http\Controllers\PermissionController;
use App\Http\Controllers\SysLogController;

// 认证路由
Route::prefix('auth')->group(function () {
    Route::post('/login', [AuthController::class, 'login']);
    Route::get('/captcha', [AuthController::class, 'captcha']);
    
    Route::middleware('auth:api')->group(function () {
        Route::post('/logout', [AuthController::class, 'logout']);
        Route::post('/current-role/switch/{roleCode}', [AuthController::class, 'switchRole']);
    });
});

// 需要认证的路由
Route::middleware('auth:api')->group(function () {
    
    // 用户路由
    Route::prefix('user')->group(function () {
        Route::get('/detail', [UserController::class, 'detail']);
        Route::get('/', [UserController::class, 'index']);
        Route::post('/', [UserController::class, 'store']);
        Route::patch('/{id}', [UserController::class, 'update']);
        Route::patch('/profile/{id}', [UserController::class, 'updateProfile']); // 新增：个人资料更新
        Route::delete('/{id}', [UserController::class, 'destroy']);
        Route::patch('/password/reset/{id}', [UserController::class, 'resetPassword']);
    });

    // 认证扩展路由（需要认证）
    Route::post('/auth/password', [AuthController::class, 'changePassword']); // 新增：修改密码

    // 角色路由
    Route::prefix('role')->group(function () {
        Route::get('/permissions/tree', [RoleController::class, 'permissionsTree']);
        Route::get('/page', [RoleController::class, 'page']);
        Route::get('/', [RoleController::class, 'index']);
        Route::post('/', [RoleController::class, 'store']);
        Route::patch('/{id}', [RoleController::class, 'update']);
        Route::delete('/{id}', [RoleController::class, 'destroy']);
        Route::patch('/users/add/{id}', [RoleController::class, 'addUsers']);
        Route::patch('/users/remove/{id}', [RoleController::class, 'removeUsers']);
    });

    // 权限路由
    Route::prefix('permission')->group(function () {
        Route::get('/tree', [PermissionController::class, 'tree']);
        Route::get('/menu/tree', [PermissionController::class, 'menuTree']);
        Route::get('/button/{parentId}', [PermissionController::class, 'buttonPermissions']);
        Route::post('/', [PermissionController::class, 'store']);
        Route::patch('/{id}', [PermissionController::class, 'update']);
        Route::delete('/{id}', [PermissionController::class, 'destroy']);
    });

    // 日志路由
    Route::get('/syslog/list', [SysLogController::class, 'index']);
});

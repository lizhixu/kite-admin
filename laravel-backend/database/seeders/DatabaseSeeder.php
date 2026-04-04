<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\User;
use App\Models\Profile;
use App\Models\Role;
use App\Models\Permission;
use Illuminate\Support\Facades\Hash;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        // 创建权限
        $permissions = [
            // 系统管理
            ['id' => 2, 'name' => '系统管理', 'code' => 'SysMgt', 'type' => 'MENU', 'icon' => 'i-fe:grid', 'show' => true, 'enable' => true, 'order' => 1],

            // 日志管理
            ['id' => 6, 'name' => '日志管理', 'code' => 'LogMgt', 'type' => 'MENU', 'parent_id' => 2, 'path' => '/log/list', 'icon' => 'i-fe:file-text', 'component' => '/src/views/log/index.vue', 'show' => true, 'enable' => true, 'order' => 4],

            // 资源管理
            ['id' => 1, 'name' => '资源管理', 'code' => 'Resource_Mgt', 'type' => 'MENU', 'parent_id' => 2, 'path' => '/pms/resource', 'icon' => 'i-fe:list', 'component' => '/src/views/pms/resource/index.vue', 'show' => true, 'enable' => true, 'order' => 1],
            ['id' => 11, 'name' => '新增资源', 'code' => 'AddResource', 'type' => 'BUTTON', 'parent_id' => 1, 'show' => true, 'enable' => true, 'order' => 1],
            ['id' => 12, 'name' => '编辑资源', 'code' => 'EditResource', 'type' => 'BUTTON', 'parent_id' => 1, 'show' => true, 'enable' => true, 'order' => 2],
            ['id' => 13, 'name' => '删除资源', 'code' => 'DeleteResource', 'type' => 'BUTTON', 'parent_id' => 1, 'show' => true, 'enable' => true, 'order' => 3],

            // 角色管理
            ['id' => 3, 'name' => '角色管理', 'code' => 'RoleMgt', 'type' => 'MENU', 'parent_id' => 2, 'path' => '/pms/role', 'icon' => 'i-fe:user-check', 'component' => '/src/views/pms/role/index.vue', 'show' => true, 'enable' => true, 'order' => 2],
            ['id' => 5, 'name' => '分配用户', 'code' => 'RoleUser', 'type' => 'MENU', 'parent_id' => 3, 'path' => '/pms/role/user/:roleId', 'icon' => 'i-fe:user-plus', 'component' => '/src/views/pms/role/role-user.vue', 'show' => false, 'enable' => true, 'order' => 1],
            ['id' => 14, 'name' => '新增角色', 'code' => 'AddRole', 'type' => 'BUTTON', 'parent_id' => 3, 'show' => true, 'enable' => true, 'order' => 1],
            ['id' => 15, 'name' => '编辑角色', 'code' => 'EditRole', 'type' => 'BUTTON', 'parent_id' => 3, 'show' => true, 'enable' => true, 'order' => 2],
            ['id' => 16, 'name' => '删除角色', 'code' => 'DeleteRole', 'type' => 'BUTTON', 'parent_id' => 3, 'show' => true, 'enable' => true, 'order' => 3],
            ['id' => 17, 'name' => '分配权限', 'code' => 'AssignPermission', 'type' => 'BUTTON', 'parent_id' => 3, 'show' => true, 'enable' => true, 'order' => 4],

            // 用户管理
            ['id' => 4, 'name' => '用户管理', 'code' => 'UserMgt', 'type' => 'MENU', 'parent_id' => 2, 'path' => '/pms/user', 'icon' => 'i-fe:user', 'component' => '/src/views/pms/user/index.vue', 'keep_alive' => true, 'show' => true, 'enable' => true, 'order' => 3],
            ['id' => 18, 'name' => '新增用户', 'code' => 'AddUser', 'type' => 'BUTTON', 'parent_id' => 4, 'show' => true, 'enable' => true, 'order' => 1],
            ['id' => 19, 'name' => '编辑用户', 'code' => 'EditUser', 'type' => 'BUTTON', 'parent_id' => 4, 'show' => true, 'enable' => true, 'order' => 2],
            ['id' => 20, 'name' => '删除用户', 'code' => 'DeleteUser', 'type' => 'BUTTON', 'parent_id' => 4, 'show' => true, 'enable' => true, 'order' => 3],
            ['id' => 21, 'name' => '重置密码', 'code' => 'ResetPassword', 'type' => 'BUTTON', 'parent_id' => 4, 'show' => true, 'enable' => true, 'order' => 4],

            // 个人资料
            ['id' => 8, 'name' => '个人资料', 'code' => 'UserProfile', 'type' => 'MENU', 'path' => '/profile', 'icon' => 'i-fe:user', 'component' => '/src/views/profile/index.vue', 'show' => false, 'enable' => true, 'order' => 99],
        ];

        foreach ($permissions as $permission) {
            Permission::create($permission);
        }

        // 创建角色
        $superAdminRole = Role::create([
            'code' => 'SUPER_ADMIN',
            'name' => '超级管理员',
            'enable' => true,
        ]);

        $qaRole = Role::create([
            'code' => 'ROLE_QA',
            'name' => '质检员',
            'enable' => true,
        ]);

        // 分配所有权限给超级管理员
        $superAdminRole->permissions()->attach(Permission::pluck('id'));

        // 创建管理员用户
        $admin = User::create([
            'username' => 'admin',
            'password' => Hash::make('123456'),
            'enable' => true,
        ]);

        // 创建用户资料
        Profile::create([
            'user_id' => $admin->id,
            'nick_name' => 'Admin',
            'avatar' => 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif?imageView2/1/w/80/h/80',
        ]);

        // 分配角色
        $admin->roles()->attach([$superAdminRole->id, $qaRole->id]);
    }
}

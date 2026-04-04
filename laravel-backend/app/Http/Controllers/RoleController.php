<?php

namespace App\Http\Controllers;

use App\Models\Role;
use Illuminate\Http\Request;

class RoleController extends Controller
{
    public function permissionsTree(Request $request)
    {
        $user = auth()->user();
        $currentRoleCode = auth()->payload()->get('role');
        $currentRole = $user->roles()->where('code', $currentRoleCode)->first() 
                    ?? $user->roles()->first();

        if (!$currentRole) {
            return response()->json([
                'code' => 403,
                'message' => '未分配角色',
                'data' => [],
                'originUrl' => $request->path(),
            ], 403);
        }

        // 如果是超级管理员，返回所有权限
        if ($currentRole->code === 'SUPER_ADMIN') {
            $permissionController = new PermissionController();
            return $permissionController->tree($request);
        }

        // 其他角色返回该角色的权限树
        $permissions = $currentRole->permissions()->orderBy('order')->get();
        
        // 使用 PermissionController 中的 buildTree 逻辑（假设已定义为私有或移动到更合适的地方）
        // 这里我们直接在这个类里实现一个简易版或调用 PermissionController 的静态方法（如果支持）
        // 为了方便，我们在 PermissionController 中把 buildTree 改为 public
        $permissionController = new PermissionController();
        $allPermissions = $currentRole->permissions()->orderBy('order')->get();
        
        // 注意：这里需要递归构建树，且只包含该角色拥有的权限
        // 我们利用 allPermissions 构建
        $tree = $this->buildPermissionsTree($allPermissions);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $tree,
            'originUrl' => $request->path(),
        ]);
    }

    private function buildPermissionsTree($items, $parentId = null)
    {
        $res = [];
        foreach ($items as $item) {
            if ($item->parent_id == $parentId) {
                // 模拟 formatPermission
                $data = [
                    'id' => $item->id,
                    'name' => $item->name,
                    'code' => $item->code,
                    'type' => $item->type,
                    'parentId' => $item->parent_id,
                    'path' => $item->path,
                    'redirect' => $item->redirect,
                    'icon' => $item->icon,
                    'component' => $item->component,
                    'layout' => $item->layout,
                    'keepAlive' => $item->keep_alive,
                    'method' => $item->method,
                    'description' => $item->description,
                    'show' => $item->show,
                    'enable' => $item->enable,
                    'order' => $item->order,
                ];
                $children = $this->buildPermissionsTree($items, $item->id);
                if (!empty($children)) {
                    $data['children'] = $children;
                }
                $res[] = $data;
            }
        }
        return $res;
    }

    public function index(Request $request)
    {
        $roles = Role::where('enable', true)->get();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $roles->map(fn($role) => [
                'id' => $role->id,
                'code' => $role->code,
                'name' => $role->name,
                'enable' => $role->enable,
            ]),
            'originUrl' => $request->path(),
        ]);
    }

    public function page(Request $request)
    {
        $query = Role::with('permissions');

        if ($request->has('name')) {
            $query->where('name', 'like', '%' . $request->name . '%');
        }

        $pageNo = $request->input('pageNo', 1);
        $pageSize = $request->input('pageSize', 10);

        $total = $query->count();
        $roles = $query->skip(($pageNo - 1) * $pageSize)
                       ->take($pageSize)
                       ->get();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'pageData' => $roles->map(fn($role) => [
                    'id' => $role->id,
                    'code' => $role->code,
                    'name' => $role->name,
                    'enable' => $role->enable,
                    'permissionIds' => $role->permissions->pluck('id'),
                    'createTime' => $role->created_at,
                    'updateTime' => $role->updated_at,
                ]),
                'total' => $total,
            ],
            'originUrl' => $request->path(),
        ]);
    }

    public function store(Request $request)
    {
        $request->validate([
            'code' => 'required|string|unique:roles',
            'name' => 'required|string',
            'enable' => 'boolean',
            'permissionIds' => 'array',
        ]);

        $role = Role::create($request->only(['code', 'name', 'enable']));

        if ($request->has('permissionIds')) {
            $role->permissions()->attach($request->permissionIds);
        }

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $role,
            'originUrl' => $request->path(),
        ]);
    }

    public function update(Request $request, $id)
    {
        $role = Role::findOrFail($id);

        $request->validate([
            'code' => 'string|unique:roles,code,' . $id,
            'name' => 'string',
            'enable' => 'boolean',
            'permissionIds' => 'array',
        ]);

        $role->update($request->only(['code', 'name', 'enable']));

        if ($request->has('permissionIds')) {
            $role->permissions()->sync($request->permissionIds);
        }

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $role,
            'originUrl' => $request->path(),
        ]);
    }


    public function destroy(Request $request, $id)
    {
        $role = Role::findOrFail($id);
        $role->delete();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    public function addUsers(Request $request, $id)
    {
        $role = Role::findOrFail($id);

        $request->validate([
            'userIds' => 'required|array',
        ]);

        $role->users()->syncWithoutDetaching($request->userIds);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    public function removeUsers(Request $request, $id)
    {
        $role = Role::findOrFail($id);

        $request->validate([
            'userIds' => 'required|array',
        ]);

        $role->users()->detach($request->userIds);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }
}

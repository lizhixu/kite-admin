<?php

namespace App\Http\Controllers;

use App\Models\Permission;
use Illuminate\Http\Request;

class PermissionController extends Controller
{
    public function tree(Request $request)
    {
        $allPermissions = Permission::orderBy('order')->get();

        $tree = $this->buildTree($allPermissions);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $tree,
            'originUrl' => $request->path(),
        ]);
    }

    public function menuTree(Request $request)
    {
        $allMenuPermissions = Permission::where('type', 'MENU')
            ->orderBy('order')
            ->get();

        $tree = $this->buildTree($allMenuPermissions);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $tree,
            'originUrl' => $request->path(),
        ]);
    }

    private function buildTree($items, $parentId = null)
    {
        $res = [];
        foreach ($items as $item) {
            if ($item->parent_id == $parentId) {
                $data = $this->formatPermission($item);
                $children = $this->buildTree($items, $item->id);
                if (!empty($children)) {
                    $data['children'] = $children;
                }
                $res[] = $data;
            }
        }
        return $res;
    }


    public function buttonPermissions(Request $request, $parentId)
    {
        $permissions = Permission::where('type', 'BUTTON')
            ->where('parent_id', $parentId)
            ->orderBy('order')
            ->get();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $permissions->map(fn($p) => $this->formatPermission($p)),
            'originUrl' => $request->path(),
        ]);
    }

    public function store(Request $request)
    {
        $request->validate([
            'name' => 'required|string',
            'code' => 'required|string|unique:permissions',
            'type' => 'required|in:MENU,BUTTON',
        ]);

        $data = $request->all();
        // 映射驼峰到下划线（为了兼容前端发送驼峰命名的参数）
        $mappedData = $this->mapCamelToSnake($data);

        $permission = Permission::create($mappedData);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $permission,
            'originUrl' => $request->path(),
        ]);
    }

    public function update(Request $request, $id)
    {
        $permission = Permission::findOrFail($id);

        $request->validate([
            'name' => 'string',
            'code' => 'string|unique:permissions,code,' . $id,
            'type' => 'in:MENU,BUTTON',
        ]);

        $data = $request->all();
        $mappedData = $this->mapCamelToSnake($data);

        $permission->update($mappedData);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $permission,
            'originUrl' => $request->path(),
        ]);
    }

    private function mapCamelToSnake($data)
    {
        $map = [
            'parentId' => 'parent_id',
            'keepAlive' => 'keep_alive',
        ];

        $res = [];
        foreach ($data as $key => $value) {
            if (isset($map[$key])) {
                $res[$map[$key]] = $value;
            } else {
                // 如果已经是下划线或不需要转换的字段
                $res[$key] = $value;
            }
        }
        return $res;
    }

    public function destroy(Request $request, $id)
    {
        $permission = Permission::findOrFail($id);
        $permission->delete();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    private function formatPermission($permission)
    {
        return [
            'id' => $permission->id,
            'name' => $permission->name,
            'code' => $permission->code,
            'type' => $permission->type,
            'parentId' => $permission->parent_id,
            'path' => $permission->path,
            'redirect' => $permission->redirect,
            'icon' => $permission->icon,
            'component' => $permission->component,
            'layout' => $permission->layout,
            'keepAlive' => $permission->keep_alive,
            'method' => $permission->method,
            'description' => $permission->description,
            'show' => $permission->show,
            'enable' => $permission->enable,
            'order' => $permission->order,
        ];
    }
}


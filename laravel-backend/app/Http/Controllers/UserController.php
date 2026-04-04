<?php

namespace App\Http\Controllers;

use App\Models\User;
use App\Models\Profile;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Hash;

class UserController extends Controller
{
    public function detail(Request $request)
    {
        $user = auth()->user();
        $user->load(['profile', 'roles']);

        // 获取当前角色
        $currentRoleCode = auth()->payload()->get('role');
        $currentRole = $user->roles()->where('code', $currentRoleCode)->first() 
                    ?? $user->roles()->first();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'id' => $user->id,
                'username' => $user->username,
                'enable' => $user->enable,
                'createTime' => $user->created_at,
                'updateTime' => $user->updated_at,
                'profile' => $user->profile ? [
                    'id' => $user->profile->id,
                    'nickName' => $user->profile->nick_name,
                    'gender' => $user->profile->gender,
                    'avatar' => $user->profile->avatar,
                    'address' => $user->profile->address,
                    'email' => $user->profile->email,
                    'userId' => $user->profile->user_id,
                ] : null,
                'roles' => $user->roles->map(fn($role) => [
                    'id' => $role->id,
                    'code' => $role->code,
                    'name' => $role->name,
                    'enable' => $role->enable,
                ]),
                'currentRole' => $currentRole ? [
                    'id' => $currentRole->id,
                    'code' => $currentRole->code,
                    'name' => $currentRole->name,
                    'enable' => $currentRole->enable,
                ] : null,
            ],
            'originUrl' => $request->path(),
        ]);
    }

    public function index(Request $request)
    {
        $query = User::with(['roles', 'profile']);

        if ($request->has('username')) {
            $query->where('username', 'like', '%' . $request->username . '%');
        }

        $pageNo = $request->input('pageNo', 1);
        $pageSize = $request->input('pageSize', 10);

        $total = $query->count();
        $users = $query->skip(($pageNo - 1) * $pageSize)
                      ->take($pageSize)
                      ->get();

        $pageData = $users->map(function($user) {
            return [
                'id' => $user->id,
                'username' => $user->username,
                'enable' => $user->enable,
                'createTime' => $user->created_at,
                'updateTime' => $user->updated_at,
                'roles' => $user->roles->map(fn($role) => [
                    'id' => $role->id,
                    'code' => $role->code,
                    'name' => $role->name,
                    'enable' => $role->enable,
                ]),
                'gender' => $user->profile?->gender,
                'avatar' => $user->profile?->avatar,
                'address' => $user->profile?->address,
                'email' => $user->profile?->email,
            ];
        });

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'pageData' => $pageData,
                'total' => $total,
            ],
            'originUrl' => $request->path(),
        ]);
    }

    public function store(Request $request)
    {
        $request->validate([
            'username' => 'required|string|unique:users',
            'password' => 'required|string|min:6',
            'enable' => 'boolean',
            'roleIds' => 'array',
        ]);

        $user = User::create([
            'username' => $request->username,
            'password' => Hash::make($request->password),
            'enable' => $request->input('enable', true),
        ]);

        if ($request->has('roleIds')) {
            $user->roles()->attach($request->roleIds);
        }

        // 创建用户资料，增加对更多字段的支持
        Profile::create([
            'user_id' => $user->id,
            'nick_name' => $request->input('nickName', $request->username),
            'gender' => $request->input('gender'),
            'avatar' => $request->input('avatar'),
            'address' => $request->input('address'),
            'email' => $request->input('email'),
        ]);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $user,
            'originUrl' => $request->path(),
        ]);
    }

    public function update(Request $request, $id)
    {
        $user = User::findOrFail($id);

        $request->validate([
            'username' => 'string|unique:users,username,' . $id,
            'enable' => 'boolean',
            'roleIds' => 'array',
        ]);

        $user->update($request->only(['username', 'enable']));

        // 更新角色
        if ($request->has('roleIds')) {
            $user->roles()->sync($request->roleIds);
        }

        // 更新个人资料
        $profileData = $request->only(['nickName', 'gender', 'avatar', 'address', 'email']);
        if (!empty($profileData)) {
            // 映射字段名
            $mappedData = [];
            if (isset($profileData['nickName'])) $mappedData['nick_name'] = $profileData['nickName'];
            if (isset($profileData['gender'])) $mappedData['gender'] = $profileData['gender'];
            if (isset($profileData['avatar'])) $mappedData['avatar'] = $profileData['avatar'];
            if (isset($profileData['address'])) $mappedData['address'] = $profileData['address'];
            if (isset($profileData['email'])) $mappedData['email'] = $profileData['email'];

            $user->profile()->updateOrCreate(['user_id' => $user->id], $mappedData);
        }

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => $user,
            'originUrl' => $request->path(),
        ]);
    }

    public function destroy(Request $request, $id)
    {
        $user = User::findOrFail($id);
        $user->delete();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    public function resetPassword(Request $request, $id)
    {
        $request->validate([
            'password' => 'required|string|min:6',
        ]);

        $user = User::findOrFail($id);
        $user->update([
            'password' => Hash::make($request->password),
        ]);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    /**
     * 更新个人资料 (由个人中心调用)
     */
    public function updateProfile(Request $request, $id)
    {
        $user = User::findOrFail($id);

        // 这里只允许更新 profile 相关字段和 username (如果需求允许)
        $profileData = $request->only(['nickName', 'gender', 'avatar', 'address', 'email']);
        
        if ($request->has('username')) {
            $user->update(['username' => $request->username]);
        }

        if (!empty($profileData)) {
            $mappedData = [];
            if (isset($profileData['nickName'])) $mappedData['nick_name'] = $profileData['nickName'];
            if (isset($profileData['gender'])) $mappedData['gender'] = $profileData['gender'];
            if (isset($profileData['avatar'])) $mappedData['avatar'] = $profileData['avatar'];
            if (isset($profileData['address'])) $mappedData['address'] = $profileData['address'];
            if (isset($profileData['email'])) $mappedData['email'] = $profileData['email'];

            $user->profile()->updateOrCreate(['user_id' => $user->id], $mappedData);
        }

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }
}

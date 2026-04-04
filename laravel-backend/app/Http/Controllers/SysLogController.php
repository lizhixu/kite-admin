<?php

namespace App\Http\Controllers;

use App\Models\SysLog;
use Illuminate\Http\Request;

class SysLogController extends Controller
{
    /**
     * 获取系统日志列表
     *
     * @param  Request  $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $page = $request->query('pageNo', 1);
        $pageSize = $request->query('pageSize', 10);
        $username = $request->query('username');
        $method = $request->query('method');
        $statusCode = $request->query('statusCode');

        $query = SysLog::query();

        if ($username) {
            $query->where('username', 'like', '%' . $username . '%');
        }

        if ($method) {
            $query->where('method', $method);
        }

        if ($statusCode) {
            $query->where('status_code', $statusCode);
        }

        $total = $query->count();
        $logs = $query->orderBy('created_at', 'desc')
                      ->offset(($page - 1) * $pageSize)
                      ->limit($pageSize)
                      ->get();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'pageData' => $logs,
                'total' => $total,
            ],
            'originUrl' => $request->path(),
        ]);
    }
}

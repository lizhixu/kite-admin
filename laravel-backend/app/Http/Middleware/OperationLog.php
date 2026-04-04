<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;
use App\Models\SysLog;
use Illuminate\Support\Facades\Auth;

class OperationLog
{
    /**
     * Handle an incoming request.
     *
     * @param  \Closure(\Illuminate\Http\Request): (\Symfony\Component\HttpFoundation\Response)  $next
     */
    public function handle(Request $request, Closure $next): Response
    {
        $startTime = microtime(true);

        // 跳过日志自身的接口和获取验证码，避免干扰
        if ($request->is('syslog/list') || $request->is('api/syslog/list') || $request->is('api/auth/captcha')) {
            return $next($request);
        }

        $response = $next($request);

        $duration = (int)((microtime(true) - $startTime) * 1000);

        $user = Auth::user();
        
        $queryString = $request->getQueryString();
        $body = $request->except(['password', 'password_confirmation']);
        $bodyStr = !empty($body) ? json_encode($body) : '';

        // 拼接 Params，逻辑对标 Go 版 operalog.go
        $params = $queryString ?: '';
        if ($request->method() !== 'GET' && !empty($bodyStr)) {
            if ($params !== '') {
                $params = "{$params} | Body: {$bodyStr}";
            } else {
                $params = $bodyStr;
            }
        }

        if (strlen($params) > 5000) {
            $params = substr($params, 0, 5000) . '... (truncated)';
        }

        $responseBody = $response->getContent();
        if (strlen($responseBody) > 5000) {
            $responseBody = substr($responseBody, 0, 5000) . '... (truncated)';
        }

        $userAgent = $request->userAgent() ?: '';
        if (strlen($userAgent) > 250) {
            $userAgent = substr($userAgent, 0, 250);
        }

        $ip = $request->ip();
        if ($ip === '::1') {
            $ip = '127.0.0.1';
        }

        try {
            SysLog::create([
                'user_id' => $user ? $user->id : null,
                'username' => $user ? $user->username : 'guest',
                'action' => $this->getActionName($request),
                'method' => $request->method(),
                'path' => $request->path(), // Go 版使用的也是 c.Request.URL.Path
                'ip' => $ip,
                'user_agent' => $userAgent,
                'request_body' => $params, // 存入数据库的 request_body 字段
                'response_body' => $responseBody, // 存入数据库的 response_body 字段
                'status_code' => $response->getStatusCode(),
                'duration' => $duration,
            ]);
        } catch (\Exception $e) {
            \Illuminate\Support\Facades\Log::error('SysLog creation failed: ' . $e->getMessage());
        }

        return $response;
    }

    private function getActionName(Request $request)
    {
        $routeName = $request->route() ? $request->route()->getName() : null;
        if ($routeName) {
            return $routeName;
        }

        return $request->method() . ' ' . $request->path();
    }
}

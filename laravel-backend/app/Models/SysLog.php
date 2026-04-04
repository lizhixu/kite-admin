<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class SysLog extends Model
{
    use HasFactory;

    protected $table = 'sys_logs';

    protected $fillable = [
        'user_id',
        'username',
        'action',
        'method',
        'path',
        'ip',
        'user_agent',
        'request_body',
        'response_body',
        'status_code',
        'duration',
    ];

    protected $casts = [
        'status_code' => 'integer',
        'duration' => 'integer',
        'created_at' => 'datetime',
    ];

    /**
     * 转为数组以匹配 Go 版的 JSON 格式
     */
    public function toArray()
    {
        return [
            'id' => $this->id,
            'userId' => $this->user_id,
            'username' => $this->username,
            'method' => $this->method,
            'path' => $this->path,
            'params' => $this->request_body,
            'response' => $this->response_body,
            'ip' => $this->ip,
            'statusCode' => $this->status_code,
            'latency' => $this->duration,
            'userAgent' => $this->user_agent,
            'createTime' => $this->created_at ? $this->created_at->format('Y-m-d H:i:s') : null,
        ];
    }

    public function user()
    {
        return $this->belongsTo(User::class);
    }
}

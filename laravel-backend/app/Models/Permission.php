<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Permission extends Model
{
    use HasFactory;

    protected $fillable = [
        'name',
        'code',
        'type',
        'parent_id',
        'path',
        'redirect',
        'icon',
        'component',
        'layout',
        'keep_alive',
        'method',
        'description',
        'show',
        'enable',
        'order',
    ];

    protected $casts = [
        'show' => 'boolean',
        'enable' => 'boolean',
        'keep_alive' => 'boolean',
        'parent_id' => 'integer',
        'order' => 'integer',
    ];

    public function parent()
    {
        return $this->belongsTo(Permission::class, 'parent_id');
    }

    public function children()
    {
        return $this->hasMany(Permission::class, 'parent_id')->orderBy('order');
    }

    public function buttonChildren()
    {
        return $this->hasMany(Permission::class, 'parent_id')->where('type', 'BUTTON')->orderBy('order');
    }

    public function roles()
    {
        return $this->belongsToMany(Role::class, 'role_permissions');
    }
}

<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use Illuminate\Support\Facades\Route;

class AppServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        //
    }

    public function boot(): void
    {
        // 注册 JWT Auth Service Provider
        $this->app->register(\Tymon\JWTAuth\Providers\LaravelServiceProvider::class);
    }
}

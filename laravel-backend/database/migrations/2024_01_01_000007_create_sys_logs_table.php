<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('sys_logs', function (Blueprint $table) {
            $table->id();
            $table->foreignId('user_id')->nullable()->constrained()->onDelete('set null');
            $table->string('username')->nullable();
            $table->string('action')->nullable();
            $table->string('method')->nullable();
            $table->string('path')->nullable();
            $table->string('ip')->nullable();
            $table->text('user_agent')->nullable();
            $table->text('request_body')->nullable();
            $table->text('response_body')->nullable();
            $table->integer('status_code')->nullable();
            $table->integer('duration')->nullable();
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('sys_logs');
    }
};

<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('permissions', function (Blueprint $table) {
            $table->id();
            $table->string('name');
            $table->string('code')->unique();
            $table->enum('type', ['MENU', 'BUTTON'])->default('MENU');
            $table->foreignId('parent_id')->nullable()->constrained('permissions')->onDelete('cascade');
            $table->string('path')->nullable();
            $table->string('redirect')->nullable();
            $table->string('icon')->nullable();
            $table->string('component')->nullable();
            $table->string('layout')->nullable();
            $table->boolean('keep_alive')->nullable();
            $table->string('method')->nullable();
            $table->text('description')->nullable();
            $table->boolean('show')->default(true);
            $table->boolean('enable')->default(true);
            $table->integer('order')->default(0);
            $table->timestamps();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('permissions');
    }
};

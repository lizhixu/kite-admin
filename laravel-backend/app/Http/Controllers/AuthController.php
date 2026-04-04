<?php

namespace App\Http\Controllers;

use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;
use Tymon\JWTAuth\Exceptions\JWTException;

class AuthController extends Controller
{
    public function login(Request $request)
    {
        $request->validate([
            'username' => 'required|string',
            'password' => 'required|string',
            'captcha' => 'required|string',
        ]);

        // 验证验证码
        $sessionCaptcha = session('captcha');
        if (!$sessionCaptcha || strtolower($request->captcha) !== strtolower($sessionCaptcha)) {
            return response()->json([
                'code' => 400,
                'message' => '验证码错误',
                'data' => null,
                'originUrl' => $request->path(),
            ], 400);
        }

        $credentials = $request->only('username', 'password');

        if (!$token = auth()->attempt($credentials)) {
            return response()->json([
                'code' => 401,
                'message' => '用户名或密码错误',
                'data' => null,
                'originUrl' => $request->path(),
            ], 401);
        }

        $user = auth()->user();
        if (!$user->enable) {
            return response()->json([
                'code' => 403,
                'message' => '用户已被禁用',
                'data' => null,
                'originUrl' => $request->path(),
            ], 403);
        }

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'accessToken' => $token,
            ],
            'originUrl' => $request->path(),
        ]);
    }

    public function captcha(Request $request)
    {
        // 生成验证码字符
        $captcha = strtoupper(Str::random(4));
        session(['captcha' => $captcha]);

        $imageData = $this->generateCaptchaImage($captcha);

        return response($imageData)->header('Content-Type', 'image/png');
    }

    private function generateCaptchaImage($captcha)
    {
        // 检查GD扩展是否可用
        if (!function_exists('imagecreate')) {
            throw new \Exception('GD extension is not enabled');
        }

        // 尝试使用 TrueType 字体
        $useTTF = function_exists('imagettftext');

        if ($useTTF) {
            $width = 80;
            $height = 40;
            $image = imagecreatetruecolor($width, $height);

            // 随机浅色背景
            imagefilledrectangle($image, 0, 0, $width, $height, imagecolorallocate($image, mt_rand(230, 250), mt_rand(230, 250), mt_rand(230, 250)));

            // 绘制干扰线
            for ($i = 0; $i < 3; $i++) {
                imageline($image, mt_rand(0, $width), mt_rand(0, $height), mt_rand(0, $width), mt_rand(0, $height), imagecolorallocate($image, mt_rand(150, 200), mt_rand(150, 200), mt_rand(150, 200)));
            }

            // 绘制噪点
            for ($i = 0; $i < 80; $i++) {
                imagesetpixel($image, mt_rand(0, $width), mt_rand(0, $height), imagecolorallocate($image, mt_rand(100, 200), mt_rand(100, 200), mt_rand(100, 200)));
            }

            // 使用系统字体（Linux常见路径）
            $fontPaths = [
                '/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf',
                '/usr/share/fonts/dejavu/DejaVuSans.ttf',
                '/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf',
                '/System/Library/Fonts/Helvetica.ttc', // macOS
                'C:/Windows/Fonts/arial.ttf', // Windows
            ];

            $fontFile = null;
            foreach ($fontPaths as $path) {
                if (file_exists($path)) {
                    $fontFile = $path;
                    break;
                }
            }

            if ($fontFile) {
                $chars = str_split($captcha);
                $charWidth = ($width - 20) / 4;

                foreach ($chars as $index => $char) {
                    $x = 10 + $charWidth * $index + mt_rand(-2, 2);
                    $y = 28 + mt_rand(-2, 2);
                    $angle = mt_rand(-15, 15);
                    imagettftext($image, 22, $angle, $x, $y, imagecolorallocate($image, mt_rand(30, 80), mt_rand(30, 80), mt_rand(30, 80)), $fontFile, $char);
                }
            } else {
                // 没有找到字体，使用内置字体
                $useTTF = false;
            }
        }

        if (!$useTTF) {
            // 回退到内置字体
            $width = 80;
            $height = 40;
            $image = imagecreatetruecolor($width, $height);

            // 随机背景颜色
            $bgR = mt_rand(220, 255);
            $bgG = mt_rand(220, 255);
            $bgB = mt_rand(220, 255);
            $bgColor = imagecolorallocate($image, $bgR, $bgG, $bgB);
            imagefilledrectangle($image, 0, 0, $width, $height, $bgColor);

            // 绘制干扰线
            for ($i = 0; $i < 3; $i++) {
                $lineR = mt_rand(150, 200);
                $lineG = mt_rand(150, 200);
                $lineB = mt_rand(150, 200);
                $lineColor = imagecolorallocate($image, $lineR, $lineG, $lineB);
                imageline($image, mt_rand(0, $width), mt_rand(0, $height), mt_rand(0, $width), mt_rand(0, $height), $lineColor);
            }

            // 绘制噪点
            for ($i = 0; $i < 50; $i++) {
                $dotR = mt_rand(100, 200);
                $dotG = mt_rand(100, 200);
                $dotB = mt_rand(100, 200);
                $dotColor = imagecolorallocate($image, $dotR, $dotG, $dotB);
                imagesetpixel($image, mt_rand(0, $width), mt_rand(0, $height), $dotColor);
            }

            // 绘制验证码字符（使用内置最大字体）
            $chars = str_split($captcha);
            $charWidth = $width / 4;

            foreach ($chars as $index => $char) {
                $x = $charWidth * $index + mt_rand(2, 6);
                $y = mt_rand(12, 18);
                $textColor = imagecolorallocate($image, mt_rand(30, 100), mt_rand(30, 100), mt_rand(30, 100));
                imagestring($image, 5, $x, $y, $char, $textColor);
            }
        }

        // 添加边框
        $borderColor = imagecolorallocate($image, 180, 180, 180);
        imagerectangle($image, 0, 0, $width - 1, $height - 1, $borderColor);

        ob_start();
        imagepng($image);
        $imageData = ob_get_clean();
        imagedestroy($image);

        return $imageData;
    }


    public function switchRole(Request $request, $roleCode)
    {
        $user = auth()->user();
        $role = $user->roles()->where('code', $roleCode)->first();

        if (!$role) {
            return response()->json([
                'code' => 404,
                'message' => '角色不存在',
                'data' => null,
                'originUrl' => $request->path(),
            ], 404);
        }

        // 生成新的token，包含角色信息
        $token = auth()->claims(['role' => $roleCode])->fromUser($user);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => [
                'accessToken' => $token,
            ],
            'originUrl' => $request->path(),
        ]);
    }

    public function logout(Request $request)
    {
        auth()->logout();

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }

    /**
     * 修改当前登录用户的密码
     */
    public function changePassword(Request $request)
    {
        $request->validate([
            'oldPassword' => 'required|string',
            'newPassword' => 'required|string|min:6',
        ]);

        $user = auth()->user();

        if (!Hash::check($request->oldPassword, $user->password)) {
            return response()->json([
                'code' => 10005,
                'message' => '原密码错误',
                'data' => null,
                'originUrl' => $request->path(),
            ]);
        }

        $user->update([
            'password' => Hash::make($request->newPassword),
        ]);

        return response()->json([
            'code' => 0,
            'message' => 'OK',
            'data' => true,
            'originUrl' => $request->path(),
        ]);
    }
}

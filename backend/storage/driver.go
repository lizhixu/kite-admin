package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// Driver 存储驱动接口
type Driver interface {
	Upload(file multipart.File, fileName string) (storagePath string, url string, err error)
	Delete(storagePath string) error
}

// -------------------- Local Driver --------------------

type LocalDriver struct {
	BasePath string // 本地存储目录，如 "uploads"
	BaseURL  string // 可访问的基础 URL，如 "http://localhost:8080"
}

func (d *LocalDriver) Upload(file multipart.File, fileName string) (string, string, error) {
	if err := os.MkdirAll(d.BasePath, 0755); err != nil {
		return "", "", fmt.Errorf("创建目录失败: %w", err)
	}
	dst, err := os.Create(filepath.Join(d.BasePath, fileName))
	if err != nil {
		return "", "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", "", fmt.Errorf("写入文件失败: %w", err)
	}
	storagePath := filepath.Join(d.BasePath, fileName)
	url := d.BaseURL + "/uploads/" + fileName
	return storagePath, url, nil
}

func (d *LocalDriver) Delete(storagePath string) error {
	if storagePath == "" {
		return nil
	}
	err := os.Remove(storagePath)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// -------------------- S3 Driver --------------------

type S3Config struct {
	Endpoint  string
	Bucket    string
	Region    string
	AccessKey string
	SecretKey string
	PublicURL string // 对外访问基础 URL
}

type S3Driver struct {
	Config S3Config
}

func (d *S3Driver) Upload(file multipart.File, fileName string) (string, string, error) {
	// 使用标准 net/http 调用 S3 兼容 API（PUT Object）
	ctx := context.Background()
	_ = ctx

	data, err := io.ReadAll(file)
	if err != nil {
		return "", "", fmt.Errorf("读取文件失败: %w", err)
	}

	endpoint := d.Config.Endpoint
	bucket := d.Config.Bucket
	region := d.Config.Region
	accessKey := d.Config.AccessKey
	secretKey := d.Config.SecretKey

	if endpoint == "" || bucket == "" {
		return "", "", fmt.Errorf("S3 配置不完整")
	}

	// 使用 awsV4Signer 签名上传
	key := fileName
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, key)
	err = s3PutObject(ctx, url, region, accessKey, secretKey, data)
	if err != nil {
		return "", "", err
	}

	publicURL := d.Config.PublicURL
	if publicURL == "" {
		publicURL = fmt.Sprintf("%s/%s", endpoint, bucket)
	}
	return key, publicURL + "/" + key, nil
}

func (d *S3Driver) Delete(storagePath string) error {
	if storagePath == "" {
		return nil
	}
	ctx := context.Background()
	endpoint := d.Config.Endpoint
	bucket := d.Config.Bucket
	url := fmt.Sprintf("%s/%s/%s", endpoint, bucket, storagePath)
	return s3DeleteObject(ctx, url, d.Config.Region, d.Config.AccessKey, d.Config.SecretKey)
}

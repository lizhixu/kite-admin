package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// s3PutObject 使用 AWS Signature V4 执行 PUT 上传
func s3PutObject(ctx context.Context, url, region, accessKey, secretKey string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.ContentLength = int64(len(data))
	signV4(req, region, accessKey, secretKey, data)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("S3 PUT 请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("S3 PUT 失败，状态码: %d", resp.StatusCode)
	}
	return nil
}

// s3DeleteObject 使用 AWS Signature V4 执行 DELETE 删除
func s3DeleteObject(ctx context.Context, url, region, accessKey, secretKey string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	signV4(req, region, accessKey, secretKey, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("S3 DELETE 请求失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 && resp.StatusCode != 404 {
		return fmt.Errorf("S3 DELETE 失败，状态码: %d", resp.StatusCode)
	}
	return nil
}

// signV4 为请求添加 AWS Signature Version 4 Authorization 头
func signV4(req *http.Request, region, accessKey, secretKey string, body []byte) {
	t := time.Now().UTC()
	dateShort := t.Format("20060102")
	dateFull := t.Format("20060102T150405Z")

	// 计算 body hash
	bodyHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // empty
	if len(body) > 0 {
		h := sha256.Sum256(body)
		bodyHash = hex.EncodeToString(h[:])
	}

	req.Header.Set("x-amz-date", dateFull)
	req.Header.Set("x-amz-content-sha256", bodyHash)
	req.Header.Set("host", req.URL.Host)

	// 构造规范请求
	canonicalHeaders := fmt.Sprintf("host:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\n",
		req.URL.Host, bodyHash, dateFull)
	signedHeaders := "host;x-amz-content-sha256;x-amz-date"
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.Method, req.URL.Path, req.URL.RawQuery,
		canonicalHeaders, signedHeaders, bodyHash)

	// 构造签名字符串
	credentialScope := fmt.Sprintf("%s/%s/s3/aws4_request", dateShort, region)
	h := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", dateFull, credentialScope, hex.EncodeToString(h[:]))

	// 计算签名
	signingKey := deriveSigningKey(secretKey, dateShort, region, "s3")
	sig := hmacSHA256(signingKey, stringToSign)
	signature := hex.EncodeToString(sig)

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s,SignedHeaders=%s,Signature=%s",
		accessKey, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", authHeader)
}

func deriveSigningKey(secretKey, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

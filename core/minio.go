package core

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"petHealthToolApi/config"
	"time"
)

// minio集成

type MinioClient struct {
	client   *minio.Client
	bucket   string
	location string
}

func init() {
	endpoint := config.Config.Oss.Endpoint
	accessKeyID := config.Config.Oss.AccessKeyId
	secretAccessKey := config.Config.Oss.AccessKeySecret
	region := config.Config.Oss.Region
	bucket := config.Config.Oss.Bucket

	_, err := NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucket, region, false)
	if err != nil {
		log.Fatalf("Failed to initialize Minio client: %v", err)
	}
}

func NewMinioClient(endpoint, accessKeyID, secretAccessKey, bucket, location string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{Region: location})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Successfully created bucket: %s\n", bucket)
	}

	return &MinioClient{
		client:   client,
		bucket:   bucket,
		location: location,
	}, nil
}

// GeneratePresignedUploadURL 生成预签名的上传URL
// objectName: 要上传的对象名称(包含路径)
// expires: URL的有效期(秒)
// 返回: 预签名的上传URL和错误信息
func (m *MinioClient) GeneratePresignedUploadURL(objectName string, expires time.Duration) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	// 设置默认过期时间为1小时
	if expires <= 0 {
		expires = time.Hour
	}

	// 生成预签名的PUT URL用于上传
	presignedURL, err := m.client.PresignedPutObject(context.Background(), m.bucket, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return presignedURL.String(), nil
}

// GeneratePresignedUploadURLWithContentType 生成带有指定Content-Type的预签名上传URL
// objectName: 要上传的对象名称(包含路径)
// contentType: 文件的Content-Type
// expires: URL的有效期(秒)
// 返回: 预签名的上传URL和错误信息
func (m *MinioClient) GeneratePresignedUploadURLWithContentType(objectName, contentType string, expires time.Duration) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("minio client not initialized")
	}

	// 设置默认过期时间为1小时
	if expires <= 0 {
		expires = time.Hour
	}

	// 生成预签名的PUT URL用于上传
	reqParams := make(map[string][]string)
	reqParams["Content-Type"] = []string{contentType}

	presignedURL, err := m.client.PresignedPutObject(context.Background(), m.bucket, objectName, expires)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	// 添加Content-Type参数
	q := presignedURL.Query()
	for k, v := range reqParams {
		q.Add(k, v[0])
	}
	presignedURL.RawQuery = q.Encode()

	return presignedURL.String(), nil
}

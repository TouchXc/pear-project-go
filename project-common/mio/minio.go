package mio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strconv"
)

//此包中主要是minio的初始化配置

type MinioClient struct {
	c *minio.Client
}

//Put上传文件

func (c *MinioClient) Put(ctx context.Context, bucketName string, fileName string, data []byte, size int64, contentType string) (minio.UploadInfo, error) {
	uploadInfo, err := c.c.PutObject(ctx, bucketName, fileName, bytes.NewBuffer(data), size, minio.PutObjectOptions{ContentType: contentType})
	return uploadInfo, err
}

//Compose合并文件

func (c *MinioClient) Compose(ctx context.Context, bucketName string, fileName string, totalChunk int) (minio.UploadInfo, error) {
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: fileName,
	}
	var srcs []minio.CopySrcOptions
	for i := 1; i <= totalChunk; i++ {
		formatInt := strconv.FormatInt(int64(i), 10)
		src := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: fileName + "_" + formatInt,
		}
		srcs = append(srcs, src)
	}
	uploadInfo, err := c.c.ComposeObject(ctx, dst, srcs...)
	return uploadInfo, err
}
func NewMinioClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*MinioClient, error) {
	// 初始化minioClient实例
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	return &MinioClient{c: minioClient}, err
}

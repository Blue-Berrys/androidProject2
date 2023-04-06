package minio

import (
	"androidProject2/config"
	"context"
	"log"

	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Init() (*minio.Client, error) { //a
	endpoint := config.Miniourl
	accessKeyID := config.MinioaccessKey
	secretAccessKey := config.MiniosecretKey
	// useSSL := true

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		// Secure: useSSL,
	})

	return minioClient, err
}
func ImageToMinio(file *multipart.FileHeader, imageName string) error {
	// 初使化 minio client对象
	minioClient, err := Init()
	if err != nil {
		log.Fatalln("创建 MinIO 客户端失败", err)

	}
	log.Printf("创建 MinIO 客户端成功")
	// minioClient.PutObject()
	// 创建一个叫 mybucket 的存储桶。
	bucketName := config.BucketName
	location := minio.MakeBucketOptions{Region: "cn-north-1"}

	err = minioClient.MakeBucket(context.Background(), bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(context.Background(), bucketName)
		if err == nil && exists {
			log.Printf("存储桶 %s 已经存在", bucketName)
		} else {
			log.Fatalln("查询存储桶状态异常", err)
			return err
		}
	}
	log.Printf("创建存储桶 %s 成功", bucketName)

	// 使用PutObject上传一个文件
	image, err := file.Open()
	if err != nil {
		log.Printf("方法file.Open() 失败%v", err)
		return err
	}
	log.Printf("方法file.Open() 成功")
	_, err = minioClient.PutObject(context.Background(), bucketName, imageName+".jpg", image, -1, minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		log.Printf("upload image failed", err)
		return err
	}
	log.Println("putObject成功")
	return nil

}

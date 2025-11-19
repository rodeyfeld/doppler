package storage

import (
	"context"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type GarageStorage struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	region          string
	useSSL          bool
}

type GarageClient struct {
	Client *minio.Client
}

func newGarageStorage() *GarageStorage {
	// Parse SSL boolean from environment variable
	useSSL, err := strconv.ParseBool(os.Getenv("S3_USE_SSL"))
	if err != nil {
		useSSL = true
		log.Printf("S3_USE_SSL not set or invalid, defaulting to true")
	}

	return &GarageStorage{
		endpoint:        os.Getenv("S3_ENDPOINT"),
		accessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		secretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		region:          os.Getenv("S3_REGION"),
		useSSL:          useSSL,
	}
}

func NewGarageClient() *GarageClient {
	gs := newGarageStorage()

	client, err := minio.New(gs.endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(gs.accessKeyID, gs.secretAccessKey, ""),
		Secure: gs.useSSL,
		Region: gs.region,
	})
	if err != nil {
		log.Panicf("Minio client creation failed: %v", err)
	}

	return &GarageClient{
		Client: client,
	}
}

func (gc *GarageClient) StoreObject(bucketName string, filename string, filepath string) {
	ctx := context.Background()
	contentType := "image/png"

	_, err := gc.Client.FPutObject(ctx, bucketName, filename, filepath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Panicf("Failed to store object using S3: %v", err)
	}
}

func (gc *GarageClient) GetObject(bucketName string, filename string, filepath string) {
	ctx := context.Background()

	err := gc.Client.FGetObject(ctx, bucketName, filename, filepath, minio.GetObjectOptions{})
	if err != nil {
		log.Panicf("Failed to get object from S3: %v", err)
	}
}

// GetObjectStream fetches an object from S3 and returns a reader for streaming
func (gc *GarageClient) GetObjectStream(bucketName string, filename string) (io.ReadCloser, string, error) {
	ctx := context.Background()

	object, err := gc.Client.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", err
	}

	// Get object info to determine content type
	objInfo, err := object.Stat()
	if err != nil {
		object.Close()
		return nil, "", err
	}

	contentType := objInfo.ContentType
	if contentType == "" {
		contentType = "image/png" // Default content type
	}

	return object, contentType, nil
}

// GetPresignedURL generates a presigned URL for accessing an object
// The URL expires after the specified duration (default: 24 hours)
// If S3_PUBLIC_ENDPOINT is set, it replaces the internal endpoint with the public one
func (gc *GarageClient) GetPresignedURL(bucketName string, filename string, expiry time.Duration) (*url.URL, error) {
	ctx := context.Background()

	if expiry == 0 {
		expiry = 24 * time.Hour // Default to 24 hours
	}

	presignedURL, err := gc.Client.PresignedGetObject(ctx, bucketName, filename, expiry, nil)
	if err != nil {
		return nil, err
	}

	// Replace internal endpoint with public endpoint if configured
	publicEndpoint := os.Getenv("S3_PUBLIC_ENDPOINT")
	if publicEndpoint != "" {
		publicURL, err := url.Parse(publicEndpoint)
		if err != nil {
			log.Printf("Failed to parse S3_PUBLIC_ENDPOINT: %v", err)
			return presignedURL, nil
		}

		// Replace the scheme and host in the presigned URL with the public endpoint
		presignedURL.Scheme = publicURL.Scheme
		presignedURL.Host = publicURL.Host
	}

	return presignedURL, nil
}

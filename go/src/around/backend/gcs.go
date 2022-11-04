package backend

import (
	"context"
	"fmt"
	"io"
	"around/util"
	"cloud.google.com/go/storage"
)

var (
	GCSBackend *GoogleCloudStorageBackend
)

type GoogleCloudStorageBackend struct {
	client *storage.Client
	bucket string
}

// func InitGCSBackend(config *util.GCSInfo) {
func InitGCSBackend(config *util.ElasticsearchInfo) {
	client, err := storage.NewClient(context.Background())
	if err != nil {
		panic(err)
	}
	// fmt.Printf("config.Bucket is : %v \n", config.Bucket)

	GCSBackend = &GoogleCloudStorageBackend {
		client: client,
		bucket: config.Bucket,
	}
}

func (backend *GoogleCloudStorageBackend) SaveToGCS(r io.Reader, objectName string) (string, error) {
	ctx := context.Background()
	object := backend.client.Bucket(backend.bucket).Object(objectName)
	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	fmt.Printf("File is save to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}
// Package storage membungkus MinIO (S3-compatible). Satu bucket, prefix per
// tenant; object key SELALU dikonstruksi server-side dari tenant_id context —
// tidak pernah menerima path dari klien (docs/07-storage-media.md §1).
package storage

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	mc     *minio.Client
	bucket string
}

type Options struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func New(ctx context.Context, opt Options) (*Client, error) {
	mc, err := minio.New(opt.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opt.AccessKey, opt.SecretKey, ""),
		Secure: opt.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("storage: init minio: %w", err)
	}
	exists, err := mc.BucketExists(ctx, opt.Bucket)
	if err != nil {
		return nil, fmt.Errorf("storage: cek bucket: %w", err)
	}
	if !exists {
		if err := mc.MakeBucket(ctx, opt.Bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("storage: buat bucket: %w", err)
		}
	}
	return &Client{mc: mc, bucket: opt.Bucket}, nil
}

// ObjectKey membangun path "t/{tenant_id}/..." — satu-satunya pabrik key.
// RELASI: internal/modules/media — service presign/confirm memanggil ini dengan
// tenant_id dari context middleware, bukan dari payload klien.
func (c *Client) ObjectKey(tenantID uuid.UUID, parts ...string) string {
	return path.Join(append([]string{"t", tenantID.String()}, parts...)...)
}

// PresignPut memberi URL unggah langsung ke MinIO sehingga file tidak membebani
// API server (docs/07-storage-media.md §2).
func (c *Client) PresignPut(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	u, err := c.mc.PresignedPutObject(ctx, c.bucket, objectKey, expiry)
	if err != nil {
		return "", fmt.Errorf("storage: presign put %s: %w", objectKey, err)
	}
	return u.String(), nil
}

func (c *Client) PresignGet(ctx context.Context, objectKey string, expiry time.Duration) (string, error) {
	u, err := c.mc.PresignedGetObject(ctx, c.bucket, objectKey, expiry, url.Values{})
	if err != nil {
		return "", fmt.Errorf("storage: presign get %s: %w", objectKey, err)
	}
	return u.String(), nil
}

// Stat dipakai alur confirm untuk memverifikasi objek benar-benar terunggah
// beserta ukurannya sebelum baris `media` dicatat.
func (c *Client) Stat(ctx context.Context, objectKey string) (minio.ObjectInfo, error) {
	return c.mc.StatObject(ctx, c.bucket, objectKey, minio.StatObjectOptions{})
}

// Ping memverifikasi konektivitas untuk /readyz.
func (c *Client) Ping(ctx context.Context) error {
	_, err := c.mc.BucketExists(ctx, c.bucket)
	return err
}

// NEXT: lifecycle rule penghapusan prefix tenant `cancelled` setelah masa arsip
// 60 hari + job pembersih (docs/07-storage-media.md §5).

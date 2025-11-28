package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"lalan-be/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

/*
Konstanta batasan upload file global.
*/
const (
	MaxFileSize        = 10 * 1024 * 1024 // 10 MB
	MaxImageSize       = 5 * 1024 * 1024  // 5 MB
	MaxDocumentSize    = 10 * 1024 * 1024 // 10 MB
	PresignedURLExpiry = 15 * time.Minute
)

/*
AllowedImageTypes & AllowedDocumentTypes adalah whitelist MIME type yang diperbolehkan.
*/
var (
	AllowedImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/webp": true,
	}
	AllowedDocumentTypes = map[string]bool{
		"application/pdf": true,
	}
)

/*
FileMetadata berisi informasi lengkap file yang baru di-upload.
Digunakan sebagai return value UploadFile().
*/
type FileMetadata struct {
	FileName    string
	FileSize    int64
	ContentType string
	URL         string
	Path        string
	UploadedAt  time.Time
}

/*
Storage adalah kontrak (interface) untuk semua operasi object storage.
Implementasi saat ini: Supabase (S3-compatible).
*/
type Storage interface {
	Upload(ctx context.Context, file io.Reader, path string, contentType string) (string, error)
	UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*FileMetadata, error)
	Delete(ctx context.Context, path string) error
	Exists(ctx context.Context, path string) (bool, error)
	GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error)
}

/*
SupabaseStorage adalah implementasi Storage menggunakan Supabase Storage (S3-compatible).
*/
type SupabaseStorage struct {
	config config.StorageConfig
	client *s3.Client
}

/*
NewSupabaseStorage membuat instance storage dengan konfigurasi eksplisit.
*/
func NewSupabaseStorage(cfg config.StorageConfig) *SupabaseStorage {
	return &SupabaseStorage{config: cfg}
}

/*
NewSupabaseStorageFromEnv membuat instance storage dari environment/config terpusat.
*/
func NewSupabaseStorageFromEnv() *SupabaseStorage {
	cfg := config.LoadStorageConfig()
	return NewSupabaseStorage(cfg)
}

/*
getClient melakukan lazy initialization S3 client (hanya dibuat sekali).
*/
func (s *SupabaseStorage) getClient() (*s3.Client, error) {
	if s.client != nil {
		return s.client, nil
	}

	cfgAWS, err := awscfg.LoadDefaultConfig(context.TODO(),
		awscfg.WithRegion(s.config.Region),
		awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			s.config.AccessKey, s.config.SecretKey, "",
		)),
		awscfg.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: s.config.Endpoint}, nil
		})),
	)
	if err != nil {
		log.Printf("SupabaseStorage getClient: failed to load AWS config: %v", err)
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	s.client = s3.NewFromConfig(cfgAWS, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	return s.client, nil
}

/*
Upload mengunggah file dari io.Reader ke path tertentu di Supabase.

Alur kerja:
1. Lazy init client
2. Sanitasi path
3. PutObject ke bucket
4. Bangun public URL

Output sukses:
- string URL publik file
Output error:
- error → gagal init client / upload / network
*/
func (s *SupabaseStorage) Upload(ctx context.Context, file io.Reader, path string, contentType string) (string, error) {
	client, err := s.getClient()
	if err != nil {
		return "", err
	}

	path = sanitizePath(path)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(path),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("SupabaseStorage Upload: failed to upload %s: %v", path, err)
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	publicURL := s.buildPublicURL(path)
	log.Printf("SupabaseStorage Upload: success %s → %s", path, publicURL)
	return publicURL, nil
}

/*
UploadFile mengunggah file multipart dengan validasi ukuran & tipe.

Alur kerja:
1. Validasi ukuran ≤ MaxFileSize
2. Deteksi dan validasi Content-Type berdasarkan folder
3. Generate nama unik + ekstensi
4. Upload via Upload()
5. Return metadata lengkap

Output sukses:
- *FileMetadata
Output error:
- error → ukuran/tipe tidak valid / gagal buka file / gagal upload
*/
func (s *SupabaseStorage) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*FileMetadata, error) {
	if fileHeader.Size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	if err := validateContentType(contentType, folder); err != nil {
		return nil, err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	path := buildFilePath(folder, uniqueName)

	url, err := s.Upload(ctx, file, path, contentType)
	if err != nil {
		return nil, err
	}

	metadata := &FileMetadata{
		FileName:    fileHeader.Filename,
		FileSize:    fileHeader.Size,
		ContentType: contentType,
		URL:         url,
		Path:        path,
		UploadedAt:  time.Now(),
	}

	log.Printf("SupabaseStorage UploadFile: uploaded %s → %s", fileHeader.Filename, url)
	return metadata, nil
}

/*
Delete menghapus file dari Supabase (best-effort, tidak error jika sudah tidak ada).

Output sukses:
- nil (bahkan jika file sudah tidak ada)
Output error:
- error → gagal init client / network / HeadObject error selain 404
*/
func (s *SupabaseStorage) Delete(ctx context.Context, path string) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}

	path = sanitizePath(path)

	exists, err := s.Exists(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to check file existence: %w", err)
	}
	if !exists {
		log.Printf("SupabaseStorage Delete: file %s already gone, skipping", path)
		return nil
	}

	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Printf("SupabaseStorage Delete: failed to delete %s: %v", path, err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	log.Printf("SupabaseStorage Delete: successfully deleted %s", path)
	return nil
}

/*
Exists mengecek apakah objek ada di bucket.

Output sukses:
- (true, nil)  → file ada
- (false, nil) → file tidak ada
Output error:
- (false, error) → kegagalan jaringan / client
*/
func (s *SupabaseStorage) Exists(ctx context.Context, path string) (bool, error) {
	client, err := s.getClient()
	if err != nil {
		return false, err
	}

	path = sanitizePath(path)

	_, err = client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "404") {
			return false, nil
		}
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return true, nil
}

/*
GetPresignedURL menghasilkan URL sementara untuk download file.

Output sukses:
- string URL presigned (default expiry 15 menit)
Output error:
- error → gagal init client / generate presign
*/
func (s *SupabaseStorage) GetPresignedURL(ctx context.Context, path string, expiry time.Duration) (string, error) {
	client, err := s.getClient()
	if err != nil {
		return "", err
	}

	path = sanitizePath(path)
	if expiry == 0 {
		expiry = PresignedURLExpiry
	}

	presignClient := s3.NewPresignClient(client)
	result, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(path),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		log.Printf("SupabaseStorage GetPresignedURL: failed for %s: %v", path, err)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	log.Printf("SupabaseStorage GetPresignedURL: generated for %s (expires in %v)", path, expiry)
	return result.URL, nil
}

/*
buildPublicURL membangun URL publik untuk file yang sudah di-upload.
*/
func (s *SupabaseStorage) buildPublicURL(path string) string {
	return fmt.Sprintf("%s/%s/%s", s.config.Domain, s.config.Bucket, path)
}

// Helper functions

/*
sanitizePath membersihkan path dari karakter berbahaya dan leading/trailing slash.
*/
func sanitizePath(path string) string {
	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "../", "")
	path = strings.ReplaceAll(path, "..\\", "")
	return path
}

/*
buildFilePath menggabungkan folder dan nama file dengan benar.
*/
func buildFilePath(folder, fileName string) string {
	if folder == "" {
		return fileName
	}
	return fmt.Sprintf("%s/%s", strings.Trim(folder, "/"), fileName)
}

/*
validateContentType memastikan MIME type sesuai dengan folder tujuan.
*/
func validateContentType(contentType, folder string) error {
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	switch folder {
	case "ktp", "profile", "avatar":
		if !AllowedImageTypes[contentType] {
			return fmt.Errorf("invalid image type: %s. Allowed: jpg, jpeg, png, webp", contentType)
		}
	case "documents", "pdf":
		if !AllowedDocumentTypes[contentType] {
			return fmt.Errorf("invalid document type: %s. Allowed: pdf", contentType)
		}
	default:
		if !AllowedImageTypes[contentType] && !AllowedDocumentTypes[contentType] {
			return fmt.Errorf("unsupported file type: %s", contentType)
		}
	}
	return nil
}

/*
ExtractPathFromURL mengekstrak path relatif dari URL publik Supabase.
*/
func ExtractPathFromURL(url, storageDomain, bucket string) string {
	prefix := fmt.Sprintf("%s/%s/", storageDomain, bucket)
	if strings.HasPrefix(url, prefix) {
		return strings.TrimPrefix(url, prefix)
	}
	return url
}

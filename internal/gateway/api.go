package gateway

import (
	"context"
	"io"
	"reflect"

	"github.com/did-server/internal/gateway/pb"
)

type UserNode interface {
	CreateBucket(ctx context.Context, bucketName string, options pb.BucketOption) (BucketInfo, error)
	HeadBucket(ctx context.Context, bucketName string) (BucketInfo, error)

	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, opts PutObjectOptions) (MefsObjectInfo, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts DownloadObjectOptions) ([]byte, error)
	HeadObject(ctx context.Context, bucketName, objectName string) (MefsObjectInfo, error)
	DeleteObject(ctx context.Context, bucketName, objectName string) error

	ShowStorage(ctx context.Context) (uint64, error)
}

type UserNodeStruct struct {
	Internal struct {
		CreateBucket func(ctx context.Context, bucketName string, opts pb.BucketOption) (BucketInfo, error)                                    `perm:"write"`
		PutObject    func(ctx context.Context, bucketName, objectName string, reader io.Reader, opts PutObjectOptions) (MefsObjectInfo, error) `perm:"write"`
		DeleteObject func(ctx context.Context, bucketName, objectName string) error                                                            `perm:"write"`

		HeadBucket func(ctx context.Context, bucketName string) (BucketInfo, error) `perm:"read"`

		GetObject  func(ctx context.Context, bucketName, objectName string, opts DownloadObjectOptions) ([]byte, error) `perm:"read"`
		HeadObject func(ctx context.Context, bucketName, objectName string) (MefsObjectInfo, error)                     `perm:"read"`

		ShowStorage func(ctx context.Context) (uint64, error) `perm:"read"`
	}
}

func (s *UserNodeStruct) CreateBucket(ctx context.Context, bucketName string, options pb.BucketOption) (BucketInfo, error) {
	return s.Internal.CreateBucket(ctx, bucketName, options)
}

func (s *UserNodeStruct) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, opts PutObjectOptions) (MefsObjectInfo, error) {
	return s.Internal.PutObject(ctx, bucketName, objectName, reader, opts)
}

func (s *UserNodeStruct) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	return s.Internal.DeleteObject(ctx, bucketName, objectName)
}

func (s *UserNodeStruct) HeadBucket(ctx context.Context, bucketName string) (BucketInfo, error) {
	return s.Internal.HeadBucket(ctx, bucketName)
}

func (s *UserNodeStruct) GetObject(ctx context.Context, bucketName, objectName string, opts DownloadObjectOptions) ([]byte, error) {
	return s.Internal.GetObject(ctx, bucketName, objectName, opts)
}

func (s *UserNodeStruct) HeadObject(ctx context.Context, bucketName, objectName string) (MefsObjectInfo, error) {
	return s.Internal.HeadObject(ctx, bucketName, objectName)
}

func (s *UserNodeStruct) ShowStorage(ctx context.Context) (uint64, error) {
	return s.Internal.ShowStorage(ctx)
}

func GetInternalStructs(in interface{}) []interface{} {
	return getInternalStructs(reflect.ValueOf(in).Elem())
}

var _internalField = "Internal"

func getInternalStructs(rv reflect.Value) []interface{} {
	var out []interface{}

	internal := rv.FieldByName(_internalField)
	ii := internal.Addr().Interface()
	out = append(out, ii)

	for i := 0; i < rv.NumField(); i++ {
		if rv.Type().Field(i).Name == _internalField {
			continue
		}

		sub := getInternalStructs(rv.Field(i))

		out = append(out, sub...)
	}

	return out
}

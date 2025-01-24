package gateway

import (
	"context"
	"io"

	"net/http"
	"os"
	"time"

	"github.com/did-server/config"
	"github.com/go-kratos/kratos/v2/log"
)

type Mefs struct {
	addr    string
	headers http.Header
	logger  *log.Helper
}

func NewStorage(logger *log.Helper) (*Mefs, error) {
	repoDir := os.Getenv("MEFS_PATH")
	addr, headers, err := getMemoClientInfo(repoDir)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	napi, closer, err := newUserNode(context.Background(), addr, headers)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	defer closer()
	_, err = napi.ShowStorage(context.Background())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Mefs{
		addr:    addr,
		headers: headers,
		logger:  logger,
	}, nil
}

func NewStorageWithApiAndToken(sapi, token string, logger *log.Helper) (*Mefs, error) {
	addr, headers, err := createMemoClientInfo(sapi, token)
	if err != nil {
		return nil, err
	}

	napi, closer, err := newUserNode(context.Background(), addr, headers)
	if err != nil {
		return nil, err
	}

	defer closer()
	_, err = napi.ShowStorage(context.Background())
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &Mefs{
		addr:    addr,
		headers: headers,
	}, nil
}

func (m *Mefs) MakeBucketWithLocation(ctx context.Context, bucket string) error {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	defer closer()
	opts := DefaultBucketOptions()

	_, err = napi.CreateBucket(ctx, bucket, opts)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

func (m *Mefs) CheckBucket(ctx context.Context, bucket string) bool {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return false
	}
	defer closer()

	bi, err := napi.HeadBucket(ctx, bucket)
	if err != nil {
		m.logger.Error(err)
		return false
	}
	return bi.Confirmed
}

func (m *Mefs) PutObject(ctx context.Context, bucket, object string) (objInfo ObjectInfo, err error) {
	path := config.CachePath + "/" + bucket + object
	fi, err := os.Open(path)
	if err != nil {
		m.logger.Error(err)
		return objInfo, err
	}
	defer fi.Close()

	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return objInfo, err
	}
	defer closer()

	poo := CidUploadOption()

	moi, err := napi.PutObject(ctx, bucket, object, fi, poo)
	if err != nil {
		m.logger.Error(err)
		return objInfo, err
	}

	etag, _ := ToString(moi.ETag)

	return ObjectInfo{
		Bucket:     bucket,
		Name:       moi.Name,
		Size:       int64(moi.Size),
		Mid:        etag,
		CreateTime: time.Unix(moi.GetTime(), 0),
	}, nil
}

func (m *Mefs) GetObject(ctx context.Context, objectName string, writer io.Writer) error {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	defer closer()

	objInfo, err := napi.HeadObject(ctx, "", objectName)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	length := int64(objInfo.Size)

	stepLen := int64(DefaultSegSize * 16)
	stepAccMax := 16

	start := int64(0)
	end := length
	stepacc := 1
	for start < end {
		if stepacc > stepAccMax {
			stepacc = stepAccMax
		}

		readLen := stepLen*int64(stepacc) - (start % stepLen)
		if end-start < readLen {
			readLen = end - start
		}

		doo := DownloadObjectOptions{
			Start:  start,
			Length: readLen,
		}

		data, err := napi.GetObject(ctx, "", objectName, doo)
		if err != nil {
			//log.Println("received length err is:", start, readLen, stepLen, err)
			m.logger.Error(err)
			break
		}
		writer.Write(data)
		start += int64(readLen)
		stepacc *= 2
	}

	return nil

}

func (m *Mefs) GetObjectInfoByMid(ctx context.Context, mid string) (ObjectInfo, error) {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return ObjectInfo{}, err
	}
	defer closer()

	obi, err := napi.HeadObject(ctx, "", mid)
	if err != nil {
		m.logger.Error(err)
		return ObjectInfo{}, err
	}

	return ObjectInfo{
		Name: obi.Name,
		Size: int64(obi.Size),
	}, nil
}

func (m *Mefs) GetObjectInfo(ctx context.Context, bucket, object string) (ObjectInfo, error) {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return ObjectInfo{}, err
	}
	defer closer()

	obi, err := napi.HeadObject(ctx, bucket, object)
	if err != nil {
		m.logger.Error(err)
		return ObjectInfo{}, err
	}

	etag, _ := EtagToString(obi.ETag)

	return ObjectInfo{
		Bucket:     bucket,
		Name:       obi.Name,
		Mid:        etag,
		Size:       int64(obi.Size),
		CreateTime: time.Unix(obi.GetTime(), 0),
	}, nil
}

func (m *Mefs) DeleteObject(ctx context.Context, bucket, object string) error {
	napi, closer, err := newUserNode(ctx, m.addr, m.headers)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	defer closer()

	err = napi.DeleteObject(ctx, bucket, object)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	return nil
}

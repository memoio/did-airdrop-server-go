package gateway

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/did-server/internal/gateway/pb"

	"github.com/ipfs/go-cid"
)

type StreamType string

const (
	Null           StreamType = "null"
	PushStream     StreamType = "push"
	DefaultSegSize            = 248 * 1024 // byte

	RsPolicy       = 1
	DefaultTagFlag = 5
)

type BucketInfo struct {
	pb.BucketOption
	pb.BucketInfo
	Confirmed bool `json:"Confirmed"`
}

type PutObjectOptions struct {
	UserDefined map[string]string
}

func DefaultUploadOption() PutObjectOptions {
	poo := PutObjectOptions{
		UserDefined: make(map[string]string),
	}

	poo.UserDefined["encryption"] = "aes2"
	poo.UserDefined["etag"] = "md5"

	return poo
}

func CidUploadOption() PutObjectOptions {
	poo := PutObjectOptions{
		UserDefined: make(map[string]string),
	}

	poo.UserDefined["encryption"] = "aes2"
	poo.UserDefined["etag"] = "cid"
	return poo
}

type MefsObjectInfo struct {
	pb.ObjectInfo
	Parts       []*pb.ObjectPartInfo `json:"Parts"`
	Size        uint64               `json:"Size"`        // file size(sum of part.RawLength)
	StoredBytes uint64               `json:"StoredBytes"` // stored size(sum of part.Length)
	Mtime       int64                `json:"Mtime"`
	State       string               `json:"State"`
	ETag        []byte               `json:"MD5"`
}

type DownloadObjectOptions struct {
	UserDefined   map[string]string
	Start, Length int64
}

func DefaultBucketOptions() pb.BucketOption {
	return pb.BucketOption{
		Version:     1,
		Policy:      RsPolicy,
		DataCount:   5,
		ParityCount: 5,
		SegSize:     DefaultSegSize,
		TagFlag:     DefaultTagFlag,
	}
}

func ToString(etag []byte) (string, error) {
	if len(etag) == md5.Size {
		return hex.EncodeToString(etag), nil
	}

	_, ecid, err := cid.CidFromBytes(etag)
	if err != nil {
		return "", err
	}

	return ecid.String(), nil
}

type ObjectInfo struct {
	Bucket     string
	Name       string
	Size       int64
	Mid        string
	CreateTime time.Time
}

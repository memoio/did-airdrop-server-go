package gateway

import (
	"crypto/md5"
	"encoding/hex"

	"github.com/ipfs/go-cid"
)

func EtagToString(etag []byte) (string, error) {
	if len(etag) == md5.Size {
		return hex.EncodeToString(etag), nil
	}

	_, ecid, err := cid.CidFromBytes(etag)
	if err != nil {
		return "", err
	}

	return ecid.String(), nil
}

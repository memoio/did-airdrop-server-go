package gateway

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

const MemoPathVar = "MEFS_PATH"
const defaultRepoDir = "~/.memo"

func getRepoPath(override string) (string, error) {
	// override is first precedence
	if override != "" {
		return homedir.Expand(override)
	}
	// Environment variable is second precedence
	envRepoDir := os.Getenv(MemoPathVar)
	if envRepoDir != "" {
		return homedir.Expand(envRepoDir)
	}
	// Default is third precedence
	return homedir.Expand(defaultRepoDir)
}

func createMemoClientInfo(api, token string) (string, http.Header, error) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+strings.TrimSpace(token))

	apima, err := multiaddr.NewMultiaddr(strings.TrimSpace(api))
	if err != nil {
		return "", nil, err
	}
	_, addr, err := manet.DialArgs(apima)
	if err != nil {
		return "", nil, err
	}

	//addr = "http://" + addr + "/rpc/v0"
	return addr, headers, nil
}

func getMemoClientInfo(repoDir string) (string, http.Header, error) {
	repoPath, err := getRepoPath(repoDir)
	if err != nil {
		return "", nil, err
	}

	tokePath := path.Join(repoPath, "token")
	tokenBytes, err := os.ReadFile(tokePath)
	if err != nil {
		return "", nil, err
	}
	tokenBytes = bytes.TrimSpace(tokenBytes)
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+string(tokenBytes))

	rpcPath := path.Join(repoPath, "api")
	rpcBytes, err := os.ReadFile(rpcPath)
	if err != nil {
		return "", nil, err
	}
	rpcBytes = bytes.TrimSpace(rpcBytes)
	apima, err := multiaddr.NewMultiaddr(string(rpcBytes))
	if err != nil {
		return "", nil, err
	}
	_, addr, err := manet.DialArgs(apima)
	if err != nil {
		return "", nil, err
	}

	return addr, headers, nil
}

func newUserNode(ctx context.Context, addr string, requestHeader http.Header) (UserNode, jsonrpc.ClientCloser, error) {
	var res UserNodeStruct
	re := readerParamEncoder("http://" + addr + "/rpc/streams/v0/push")
	closer, err := jsonrpc.NewMergeClient(ctx, "ws://"+addr+"/rpc/v0", "Memoriae",
		GetInternalStructs(&res), requestHeader, re)

	return &res, closer, err
}

type ReaderStream struct {
	Type StreamType
	Info string
}

func readerParamEncoder(addr string) jsonrpc.Option {
	return jsonrpc.WithParamEncoder(new(io.Reader), func(value reflect.Value) (reflect.Value, error) {
		r := value.Interface().(io.Reader)

		reqID := uuid.New()
		u, err := url.Parse(addr)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("parsing push address: %w", err)
		}
		u.Path = path.Join(u.Path, reqID.String())

		go func() {
			resp, err := http.Post(u.String(), "application/octet-stream", r)
			if err != nil {

				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				b, _ := io.ReadAll(resp.Body)
				fmt.Println(string(b))
				return
			}

		}()

		return reflect.ValueOf(ReaderStream{Type: PushStream, Info: reqID.String()}), nil
	})
}

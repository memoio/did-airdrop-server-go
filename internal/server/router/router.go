package router

import (
	"os"

	"github.com/did-server/internal/did"
	"github.com/did-server/internal/gateway"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	klog "github.com/go-kratos/kratos/v2/log"
)

type handle struct {
	logger  *klog.Helper
	did     *did.MemoDID
	gateway *gateway.Mefs
}

func NewRouter(chain string, r *gin.Engine) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	loggers := klog.NewHelper(logger)
	did, err := did.NewMemoDID(chain, loggers)
	if err != nil {
		panic(err)
	}

	gateway, err := gateway.NewStorage(log.NewHelper(logger))
	if err != nil {
		log.NewHelper(logger).Error(err)
		return
	}

	h := &handle{
		did:     did,
		logger:  loggers,
		gateway: gateway,
	}

	loadDIDmoudles(r.Group("/did"), h)
	loadMfileDIDMoudles(r.Group("/mfile"), h)
	loadFileMoudles(r.Group("/file"), h)
}

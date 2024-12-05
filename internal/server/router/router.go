package router

import (
	"os"

	"github.com/did-server/internal/did"
	"github.com/gin-gonic/gin"
	klog "github.com/go-kratos/kratos/v2/log"
)

type handle struct {
	logger *klog.Helper
	did    *did.MemoDID
}

func NewRouter(r *gin.Engine) {
	logger := klog.With(klog.NewStdLogger(os.Stdout),
		"ts", klog.DefaultTimestamp,
		"caller", klog.DefaultCaller,
	)

	loggers := klog.NewHelper(logger)
	did, err := did.NewMemoDID("dev", loggers)
	if err != nil {
		panic(err)
	}

	h := &handle{
		did:    did,
		logger: loggers,
	}
	loadDIDmoudles(r.Group("/did"), h)
}

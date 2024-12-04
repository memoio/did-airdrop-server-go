package router

import "github.com/gin-gonic/gin"

type handle struct {
}

func NewRouter(r *gin.Engine) {
	h := &handle{}
	loadDIDmoudles(r.Group("/did"), h)
}

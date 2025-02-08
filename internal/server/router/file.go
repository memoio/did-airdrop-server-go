package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func loadFileMoudles(r *gin.RouterGroup, h *handle) {
	r.POST("/upload", h.fileUpload)
	r.GET("/download", h.fileDownload)
	r.GET("/list", h.fileList)
}

// @Summary file upload
// @Description upload file
// @Tags file
// @Accept json
// @Produce json
// @Param data body string true "data"
// @Param address body string true "address"
// @Success 200
// @Router /file/upload [post]
func (h *handle) fileUpload(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	data, ok := body["data"].(string)
	if !ok {
		err := fmt.Errorf("invalid data")
		h.logger.Error(err)
		c.JSON(http.StatusOK, ErrParamsInvalid)
		return
	}

	address := body["address"].(string)
	if address == "" {
		err := fmt.Errorf("invalid address")
		h.logger.Error(err)
		c.JSON(http.StatusOK, ErrParamsInvalid)
		return
	}

	err := h.gateway.MakeBucketWithLocation(c.Request.Context(), address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusOK, ErrParamsInvalid)
		return
	}

	h.gateway.PutObject(c.Request.Context(), address, data)
}

// @Summary file download
// @Description download file
// @Tags file
// @Accept json
// @Produce json
// @Param address body string true "address"
// @Success 200
// @Router /file/download [get]
func (h *handle) fileDownload(c *gin.Context) {}

// @Summary file list
// @Description list file
// @Tags file
// @Accept json
// @Produce json
// @Param address body string true "address"
// @Success 200
// @Router /file/list [get]
func (h *handle) fileList(c *gin.Context) {}

package router

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"mime"
	"net/http"

	"github.com/did-server/config"
	"github.com/did-server/internal/did"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

func loadMfileDIDMoudles(r *gin.RouterGroup, h *handle) {
	r.POST("/upload/create", h.uploadCreate)
	r.POST("/upload/confirm", h.uploadConfirm)
	r.GET("/download", h.download)
}

// @Summary		UploadCreate
// @Description	UploadCreate
// @Tags			mfile
// @Accept			json
// @Produce		json
// @Param			data	body		string	true	"data"
// @Param			address	body		string	true	"address"
// @Success		200		{string}	string	"ok"
// @Router			/mfile/upload/create [post]
func (h *handle) uploadCreate(c *gin.Context) {
	bucket := config.Bucket
	body := make(map[string]interface{})
	c.BindJSON(&body)
	data, ok := body["data"].(string)
	if !ok {
		err := fmt.Errorf("invalid data")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	address := body["address"].(string)
	if address == "" {
		err := fmt.Errorf("invalid address")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	did := body["did"].(string)
	if did == "" {
		err := fmt.Errorf("invalid did")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	pricef, ok := body["price"].(float64)
	if !ok {
		err := fmt.Errorf("invalid price")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	priceb := big.NewInt(0).SetInt64(int64(pricef))

	databyte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	object := address + hex.EncodeToString(crypto.Keccak256(databyte))

	var buf *bytes.Buffer = bytes.NewBuffer(databyte)

	err = h.createCacheFile(object, buf)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	info, err := h.gateway.PutObject(c.Request.Context(), bucket, object)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	mdid, err := h.did.CreateMfileDID(info.Mid)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	message, err := h.did.CreateMDIDMessage(info.Mid, did, priceb)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, ErrUploadFailed.Message)
		return
	}

	c.JSON(200, gin.H{"mdid": mdid.String(), "message": message})
}

// @Summary		UploadConfirm
// @Description	UploadConfirm
// @Tags			mfile
// @Accept			json
// @Produce		json
// @Param			sig		body		string	true	"sig"
// @Param			address	body		string	true	"address"
// @Success		200		{string}	string	"ok"
// @Router			/mfile/upload/confirm [post]
func (h *handle) uploadConfirm(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	sig, ok := body["sig"].(string)
	if !ok {
		err := fmt.Errorf("invalid sig")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, err.Error())
		return
	}

	mdid, ok := body["mdid"].(string)
	if !ok {
		err := fmt.Errorf("invaild address")
		h.logger.Error(err)
		c.JSON(ErrUploadFailed.Code, err)
		return
	}

	h.logger.Infof("sig %s message %s address %s", sig, mdid)

	// TODO: upload
	c.JSON(200, gin.H{"message": "success"})
}

// @Summary		Download
// @Description	Download
// @Tags			mfile
// @Accept			json
// @Produce		json
// @Param			mdid	query		string	true	"mdid"
// @Param			address	query		string	true	"address"
// @Success		200		{string}	string	"ok"
// @Router			/mfile/download [get]
func (h *handle) download(c *gin.Context) {
	mdid := c.Query("mdid")

	address := c.Query("address")

	if mdid == "" || address == "" {
		err := fmt.Errorf("invalid cid or address")
		h.logger.Error(err)
		c.JSON(ErrDownloadFailed.Code, err.Error())
		return
	}

	mfile, err := did.ParaseMfileDID(mdid)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDownloadFailed.Code, err.Error())
		return
	}

	info, err := h.gateway.GetObjectInfoByMid(c.Request.Context(), mfile.Identifier)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDownloadFailed.Code, err.Error())
		return
	}

	var w bytes.Buffer
	err = h.gateway.GetObject(c.Request.Context(), mfile.Identifier, &w)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDownloadFailed.Code, err.Error())
		return
	}

	head := fmt.Sprintf("attachment; filename=\"%s\"", info.Name)
	extraHeaders := map[string]string{
		"Content-Disposition": head,
	}

	ft := mime.TypeByExtension(info.Name)
	c.DataFromReader(http.StatusOK, info.Size, ft, &w, extraHeaders)

	c.JSON(200, "ok")
}

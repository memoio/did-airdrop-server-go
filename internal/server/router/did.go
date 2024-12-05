package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func loadDIDmoudles(r *gin.RouterGroup, h *handle) {
	r.GET("/createsigmsg", h.getCreateSigMsg)
	r.GET("/deletesigmsg", h.getDeleteSigMsg)
	r.POST("/create", h.createDID)
	r.GET("/info", h.getDIDInfo)
	r.POST("/delete", h.deleteDID)
	r.POST("/addverifyinfo", h.addVerifyInfo)
	r.POST("/changeverifyinfo", h.changeVerifyInfo)

}

// @ Summary GetCreateSigMsg
// @Description GetCreateSigMsg
// @Tags DID
// @Accept json
// @Produce json
// @Param chain query string true "The signature of the chain"
// @Param publicKey query string true "publicKey"
// @Success 200 {object} GetSigMsgResponse
// @Router /did/createsigmsg [get]
func (h *handle) getCreateSigMsg(c *gin.Context) {
	chain := c.Query("chain")

	publicKey := c.Query("publicKey")

	msg, err := h.did.GetCreateSignatureMassage(chain, publicKey)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetSignatureMessage.Code, ErrDIDGetSignatureMessage)
		return
	}
	c.JSON(200, GetSigMsgResponse{Msg: msg})
}

// @ Summary CreateDID
// @Description CreateDID
// @Tags DID
// @Accept json
// @Produce json
// @Param sig body string true "user signature"
// @Success 200 {object} CreateDIDResponse
// @Router /did/create [post]
// @Failure 502 {object} Error
// @Failure 503 {object} Error
func (h *handle) createDID(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	sig, ok := body["sig"].(string)
	if !ok {
		h.logger.Error("sig is not string", body)
		c.JSON(ErrSignatureNull.Code, ErrSignatureNull)
		return
	}

	fmt.Println(sig)

	c.JSON(200, CreateDIDResponse{DID: "did:memo:d687daa192ffa26373395872191e8502cc41fbfbf27dc07d3da3a35de57c2d96"})
}

// @ Summary GetDIDInfo
// @Description GetDIDInfo
// @Tags DID
// @Accept json
// @Produce json
// @Param did query string true "user did"
// @Success 200 {object} GetDIDInfoResponse
// @Router /did/info [get]
// @Failure 503 {object} Error
func (h *handle) getDIDInfo(c *gin.Context) {
	did := c.Query("did")
	if did == "" {
		c.JSON(ErrDIDNull.Code, ErrDIDNull)
		return
	}

	fmt.Println(did)
	c.JSON(200, GetDIDInfoResponse{DID: did, Info: []DIDInfo{
		{Chain: "memo", Balance: 100},
		{Chain: "eth", Balance: 0.01},
	}})
}

// @ Summary GetDeleteSigMsg
// @Description GetDeleteSigMsg
// @Tags DID
// @Accept json
// @Produce json
// @Param did query string true "user did"
// @Success 200 {object} GetSigMsgResponse
// @Router /did/deletesigmsg [get]
func (h *handle) getDeleteSigMsg(c *gin.Context) {
	did := c.Query("did")
	if did == "" {
		c.JSON(ErrDIDNull.Code, ErrDIDNull)
		return
	}

	msg, err := h.did.GetDeleteSignatureMassage(did)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetSignatureMessage.Code, ErrDIDGetSignatureMessage)
		return
	}

	c.JSON(200, GetSigMsgResponse{Msg: msg})
}

// @ Summary DeleteDID
// @Description DeleteDID
// @Tags DID
// @Accept json
// @Produce json
// @Param sig body string true "user signature"
// @Param did body string true "did"
// @Success 200 {object} DeleteDIDResponse
// @Router /did/delete [post]
// @Failure 502 {object} Error
// @Failure 504 {object} Error
func (h *handle) deleteDID(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	sig, ok := body["sig"].(string)
	if !ok {
		c.JSON(ErrSignatureNull.Code, ErrSignatureNull)
		return
	}

	did, ok := body["did"].(string)
	if !ok {
		c.JSON(ErrDIDNull.Code, ErrChainNull)
		return
	}

	fmt.Println(did, sig)

	c.JSON(200, DeleteDIDResponse{
		DID:    did,
		Status: "deactiaved",
	})
}

func (h *handle) addVerifyInfo(c *gin.Context) {
	c.JSON(200, AddVerifyInfoResponse{})
}

func (h *handle) changeVerifyInfo(c *gin.Context) {
	c.JSON(200, ChangeVerifyInfoResponse{})
}

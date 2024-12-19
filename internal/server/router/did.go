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
	r.GET("/exist", h.getDIDExist)
	r.GET("/number", h.getDIDNumber)

}

// @ Summary GetCreateSigMsg
// @Description GetCreateSigMsg
// @Tags DID
// @Accept json
// @Produce json
// @Param chain query string true "The signature of the chain"
// @Param address query string true "publicKey"
// @Success 200 {object} GetSigMsgResponse
// @Router /did/createsigmsg [get]
func (h *handle) getCreateSigMsg(c *gin.Context) {
	address := c.Query("address")

	msg, err := h.did.GetCreateSignatureMassageByAddress(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetSignatureMessage.Code, ErrDIDGetSignatureMessage)
		return
	}

	h.logger.Info(msg)
	c.JSON(200, GetSigMsgResponse{Msg: msg})
}

// @ Summary CreateDID
// @Description CreateDID
// @Tags DID
// @Accept json
// @Produce json
// @Param sig body string true "user signature"
// @Param address body string true "user address"
// @Success 200 {object} CreateDIDResponse
// @Router /did/create [post]
// @Failure 502 {object} Error
// @Failure 503 {object} Error
func (h *handle) createDID(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	_, ok := body["sig"].(string)
	if !ok {
		h.logger.Error("sig is not string", body)
		c.JSON(ErrSignatureNull.Code, ErrSignatureNull)
		return
	}

	address, ok := body["address"].(string)
	if !ok {
		h.logger.Error("address is not string", body)
		c.JSON(ErrAddressNull.Code, ErrAddressNull)
		return
	}

	// SigByte, err := hexutil.Decode(sig)
	// if err != nil {
	// 	h.logger.Error(err)
	// 	c.JSON(ErrSignature.Code, ErrSignature)
	// 	return
	// }

	did, err := h.did.RegisterDIDByAddress(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDCreateFailed.Code, ErrDIDCreateFailed)
		return
	}

	c.JSON(200, CreateDIDResponse{DID: did})
}

// @ Summary GetDIDInfo
// @Description GetDIDInfo
// @Tags DID
// @Accept json
// @Produce json
// @Param address query string true "user did"
// @Success 200 {object} GetDIDInfoResponse
// @Router /did/info [get]
// @Failure 503 {object} Error
func (h *handle) getDIDInfo(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		h.logger.Error("address is null", address)
		c.JSON(ErrAddressNull.Code, ErrAddressNull)
		return
	}

	did, number, err := h.did.GetDIDInfo(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetInfo.Code, gin.H{"message": ErrDIDGetInfo.Message, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"did":    did,
		"number": number,
	})
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

func (h *handle) getDIDExist(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		h.logger.Error("address is null", address)
		c.JSON(ErrAddressNull.Code, ErrAddressNull)
		return
	}

	number, err := h.did.GetDIDExist(address)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetInfo.Code, gin.H{"message": ErrDIDGetInfo.Message, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"exist": number,
	})

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

func (h *handle) getDIDNumber(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		c.JSON(ErrAddressNull.Code, ErrAddressNull)
		return
	}
	num, err := h.did.GetDIDNumber()
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetInfo.Code, gin.H{"message": ErrDIDGetInfo.Message, "error": err.Error()})
		return
	}

	err = h.did.AddDIDNumber(address, num)
	if err != nil {
		h.logger.Error(err)
		c.JSON(ErrDIDGetInfo.Code, gin.H{"message": ErrDIDGetInfo.Message, "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"number": num,
	})
}

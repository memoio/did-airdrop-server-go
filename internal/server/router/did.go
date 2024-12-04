package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func loadDIDmoudles(r *gin.RouterGroup, h *handle) {
	r.POST("/create", h.createDID)
	r.GET("/info", h.getDIDInfo)
	r.POST("/delete", h.deleteDID)
	r.POST("/addverifyinfo", h.addVerifyInfo)
	r.POST("/changeverifyinfo", h.changeVerifyInfo)

}

// @ Summary CreateDID
// @Description CreateDID
// @Tags DID
// @Accept json
// @Produce json
// @Param sig body string true "signature"
// @Param chain body string true "chain"
// @Success 200 {object} CreateDIDResponse
// @Router /did/create [post]
// @Failure 502 {object} Error
// @Failure 503 {object} Error
func (h *handle) createDID(c *gin.Context) {
	body := make(map[string]interface{})
	c.BindJSON(&body)
	sig, ok := body["sig"].(string)
	if !ok {
		c.JSON(ErrSignatureNull.Code, ErrSignatureNull)
	}

	chain, ok := body["chain"].(string)
	if !ok {
		c.JSON(ErrChainNull.Code, ErrChainNull)
	}

	fmt.Println(chain, sig)

	c.JSON(200, CreateDIDResponse{DID: "did:memo:d687daa192ffa26373395872191e8502cc41fbfbf27dc07d3da3a35de57c2d96"})
}

// @ Summary GetDIDInfo
// @Description GetDIDInfo
// @Tags DID
// @Accept json
// @Produce json
// @Param did query string true "did"
// @Success 200 {object} GetDIDInfoResponse
// @Router /did/info [get]
// @Failure 503 {object} Error
func (h *handle) getDIDInfo(c *gin.Context) {
	did := c.Query("did")
	if did == "" {
		c.JSON(ErrDIDNull.Code, ErrDIDNull)
	}

	fmt.Println(did)
	c.JSON(200, GetDIDInfoResponse{DID: did, Info: []DIDInfo{
		{Chain: "memo", Balance: 100},
		{Chain: "eth", Balance: 0.01},
	}})
}

// @ Summary DeleteDID
// @Description DeleteDID
// @Tags DID
// @Accept json
// @Produce json
// @Param sig body string true "signature"
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
	}

	did, ok := body["did"].(string)
	if !ok {
		c.JSON(ErrDIDNull.Code, ErrChainNull)
	}

	fmt.Println(did, sig)

	c.JSON(200, DeleteDIDResponse{})
}

func (h *handle) addVerifyInfo(c *gin.Context) {
	c.JSON(200, AddVerifyInfoResponse{})
}

func (h *handle) changeVerifyInfo(c *gin.Context) {
	c.JSON(200, ChangeVerifyInfoResponse{})
}

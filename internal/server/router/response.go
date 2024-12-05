package router

type CreateDIDResponse struct {
	DID string `json:"did"`
}

type GetDIDInfoResponse struct {
	DID  string    `json:"did"`
	Info []DIDInfo `json:"info"`
}

type DeleteDIDResponse struct {
	DID    string `json:"did"`
	Status string `json:"status"`
}

type DIDInfo struct {
	Chain   string  `json:"chain"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type AddVerifyInfoResponse struct {
}

type ChangeVerifyInfoResponse struct {
}

type GetSigMsgResponse struct {
	Msg string `json:"msg"`
}

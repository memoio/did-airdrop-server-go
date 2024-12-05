package router

var (
	ErrDIDCreated             = Error{Code: 500, Message: "DID created"}
	ErrDIDCreateFailed        = Error{Code: 501, Message: "DID create failed"}
	ErrSignatureNull          = Error{Code: 502, Message: "Signature is null"}
	ErrChainNull              = Error{Code: 503, Message: "Chain is null"}
	ErrDIDNull                = Error{Code: 504, Message: "DID is null"}
	ErrDIDGetSignatureMessage = Error{Code: 505, Message: "Get signature message failed"}
)

type Error struct {
	Code    int
	Message string
}

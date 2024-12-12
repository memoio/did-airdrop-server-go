package router

var (
	ErrDIDCreated             = Error{Code: 500, Message: "DID created"}
	ErrDIDCreateFailed        = Error{Code: 501, Message: "DID create failed"}
	ErrSignatureNull          = Error{Code: 502, Message: "Signature is null"}
	ErrChainNull              = Error{Code: 503, Message: "Chain is null"}
	ErrDIDNull                = Error{Code: 504, Message: "DID is null"}
	ErrDIDGetSignatureMessage = Error{Code: 505, Message: "Get signature message failed"}
	ErrAddressNull            = Error{Code: 506, Message: "Address is null"}
	ErrSignature              = Error{Code: 507, Message: "Signature failed"}
	ErrDIDGetInfo             = Error{Code: 508, Message: "Get DID info failed"}
)

type Error struct {
	Code    int
	Message string
}

package router

var (
	ErrDIDCreated             = Error{Code: 550, Message: "DID created"}
	ErrDIDCreateFailed        = Error{Code: 551, Message: "DID create failed"}
	ErrSignatureNull          = Error{Code: 552, Message: "Signature is null"}
	ErrChainNull              = Error{Code: 553, Message: "Chain is null"}
	ErrDIDNull                = Error{Code: 554, Message: "DID is null"}
	ErrDIDGetSignatureMessage = Error{Code: 555, Message: "Get signature message failed"}
	ErrAddressNull            = Error{Code: 556, Message: "Address is null"}
	ErrSignature              = Error{Code: 557, Message: "Signature failed"}
	ErrDIDGetInfo             = Error{Code: 558, Message: "Get DID info failed"}
	ErrUploadFailed           = Error{Code: 559, Message: "Mfile Upload failed"}
	ErrDownloadFailed         = Error{Code: 560, Message: "Mfile download failed"}
)

type Error struct {
	Code    int
	Message string
}

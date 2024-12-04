package contract

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Controller struct {
	Endpoint      string
	ProxyAddr     common.Address
	DIDTransactor *bind.TransactOpts
}

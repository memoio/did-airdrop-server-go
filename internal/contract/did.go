package contract

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/did-solidity/go-contracts/proxy"
	"github.com/memoio/go-did/types"
)

func (c *Controller) CreateDID(publicKey, sig string) (*types.MemoDID, error) {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	proxyIns, err := proxy.NewProxy(c.proxyAddr, client)
	if err != nil {
		return nil, err
	}

	log.Println(proxyIns)

	return nil, err
}

func (c *Controller) GetNonce() (uint64, error) {
	return 0, nil
}

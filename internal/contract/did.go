package contract

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/did-solidity/go-contracts/proxy"
	"github.com/memoio/go-did/memo"
	"github.com/memoio/go-did/types"
)

func GetDIDInfo(did string) string {
	resolver, err := memo.NewMemoDIDResolver("dev")
	if err != nil {
		panic(err.Error())
	}

	document, err := resolver.Resolve(did)
	if err != nil {
		panic(err.Error())
	}

	data, err := json.Marshal(document)
	if err != nil {
		panic(err.Error())
	}

	return string(data)
}

func (c *Controller) CreateDID(publicKey, sig string) (*types.MemoDID, error) {
	client, err := ethclient.DialContext(context.TODO(), c.Endpoint)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	proxyIns, err := proxy.NewProxy(c.ProxyAddr, client)
	if err != nil {
		return nil, err
	}

	log.Println(proxyIns)

	return nil, err
}

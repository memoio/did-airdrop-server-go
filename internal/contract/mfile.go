package contract

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/did-solidity/go-contracts/proxy"
)

func (c *Controller) RegisterMfile(mfileI, didI string, price *big.Int, sig []byte) error {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	proxyIns, err := proxy.NewProxy(c.proxyAddr, client)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	tx, err := proxyIns.RegisterMfileDid(c.didTransactor, mfileI, "cid", 0, didI, price, nil, sig)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	return c.CheckTx(tx.Hash(), "RegisterDID")
}

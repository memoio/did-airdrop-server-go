package contract

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/memoio/did-solidity/go-contracts/proxy"
	"golang.org/x/xerrors"

	etypes "github.com/ethereum/go-ethereum/core/types"
	com "github.com/memoio/contractsv2/common"
	"github.com/memoio/go-did/types"
)

var (
	checkTxSleepTime = 6 // 先等待6s（出块时间加1）
	nextBlockTime    = 5 // 出块时间5s
)

func (c *Controller) RegisterDID(did, method string, publickey, sig []byte) error {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		c.logger.Error(err)
		return err
	}
	defer client.Close()

	proxyIns, err := proxy.NewProxy(c.proxyAddr, client)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	tx, err := proxyIns.CreateDID(c.didTransactor, did, method, publickey, sig)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	return c.CheckTx(tx.Hash(), "RegisterDID")
}

func (c *Controller) DeleteDID(did string, sig []byte) error {
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

	tx, err := proxyIns.DeactivateDID(c.didTransactor, did, true, sig)
	if err != nil {
		c.logger.Error(err)
		return err
	}

	return c.CheckTx(tx.Hash(), "DeactivateDID")
}

func (c *Controller) GetDIDStatus(didStr string) (bool, error) {
	did, err := types.ParseMemoDID(didStr)
	if err != nil {
		c.logger.Error(err)
		return false, err
	}

	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		c.logger.Error(err)
		return false, err
	}

	accountIns, err := proxy.NewIAccountDid(c.accountAddr, client)
	if err != nil {
		c.logger.Error(err)
		return false, err
	}

	dactivated, err := accountIns.IsDeactivated(&bind.CallOpts{}, did.Identifier)
	if err != nil {
		c.logger.Error(err)
		return false, err
	}

	return dactivated, nil
}

func (c *Controller) GetNonce(did string) (uint64, error) {
	client, err := ethclient.DialContext(context.TODO(), c.endpoint)
	if err != nil {
		c.logger.Error(err)
		return 0, err
	}

	proxyCaller, err := proxy.NewProxyCaller(c.proxyAddr, client)
	if err != nil {
		c.logger.Error(err)
		return 0, err
	}

	nonce, err := proxyCaller.GetNonce(&bind.CallOpts{}, did)
	if err != nil {
		c.logger.Error(err)
		return 0, err
	}

	return nonce, nil
}

func (c *Controller) CheckTx(txHash common.Hash, name string) error {
	var receipt *etypes.Receipt

	t := checkTxSleepTime
	for i := 0; i < 10; i++ {
		time.Sleep(time.Duration(t) * time.Second)
		receipt = com.GetTransactionReceipt(c.endpoint, txHash)
		if receipt != nil {
			break
		}
		t = nextBlockTime
	}

	if receipt == nil {
		err := xerrors.Errorf("%s: cann't get transaction(%s) receipt, not packaged", name, txHash)
		c.logger.Error(err)
		return err
	}

	// 0 means fail
	if receipt.Status == 0 {
		if receipt.GasUsed != receipt.CumulativeGasUsed {
			err := xerrors.Errorf("%s: transaction(%s) exceed gas limit", name, txHash)
			c.logger.Error(err)
			return err
		}
		err := xerrors.Errorf("%s: transaction(%s) mined but execution failed, please check your tx input", name, txHash)
		c.logger.Error(err)
		return err
	}
	return nil
}

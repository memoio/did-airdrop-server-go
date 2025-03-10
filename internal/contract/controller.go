package contract

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/did-server/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-kratos/kratos/v2/log"
	com "github.com/memoio/contractsv2/common"
	inst "github.com/memoio/contractsv2/go_contracts/instance"
)

type Controller struct {
	instanceAddr  common.Address
	endpoint      string
	privateKey    *ecdsa.PrivateKey
	didTransactor *bind.TransactOpts
	proxyAddr     common.Address
	logger        *log.Helper
	accountAddr   common.Address
}

func NewController(chain string, logger *log.Helper) (*Controller, error) {
	return NewControllerWithDID(chain, logger)
}

func NewControllerWithDID(chain string, logger *log.Helper) (*Controller, error) {
	instanceAddr, endpoint := com.GetInsEndPointByChain(chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		logger.Error(err)
		chainID = big.NewInt(985)
	}

	privateKey, err := crypto.HexToECDSA(config.Privatekey)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	instanceIns, err := inst.NewInstance(instanceAddr, client)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	// get proxyAddr
	proxyAddr, err := instanceIns.Instances(&bind.CallOpts{}, com.TypeDidProxy)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	accountAddr, err := instanceIns.Instances(&bind.CallOpts{}, com.TypeAccountDid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = big.NewInt(200)

	return &Controller{
		instanceAddr:  instanceAddr,
		endpoint:      endpoint,
		privateKey:    privateKey,
		didTransactor: auth,
		proxyAddr:     proxyAddr,
		logger:        logger,
		accountAddr:   accountAddr,
	}, nil

}

func (c *Controller) Proxy() common.Address {
	return c.proxyAddr
}

func (c *Controller) Account() common.Address {
	return c.accountAddr
}

func (c *Controller) EndPoint() string {
	return c.endpoint
}

func (c *Controller) Instance() common.Address {
	return c.instanceAddr
}

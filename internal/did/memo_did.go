package did

import (
	"context"
	"encoding/binary"
	"encoding/hex"

	"github.com/did-server/internal/contract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-kratos/kratos/v2/log"

	com "github.com/memoio/contractsv2/common"
	"github.com/memoio/go-did/types"
)

type MemoDID struct {
	Controller *contract.Controller
	chain      string
	logger     *log.Helper
}

func NewMemoDID(chain string, logger *log.Helper) (*MemoDID, error) {
	controller, err := contract.NewController(chain, logger)
	if err != nil {
		return nil, err
	}

	return &MemoDID{
		Controller: controller,
		chain:      chain,
		logger:     logger,
	}, nil
}

// Create unregistered DID
func (m *MemoDID) CreateDIDByPubKey(publicKeyStr string) (*types.MemoDID, error) {
	_, endpoint := com.GetInsEndPointByChain(m.chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}
	defer client.Close()

	publicKeyECDSA, err := m.publickeyFromString(publicKeyStr)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.TODO(), address)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	identifier := hex.EncodeToString(crypto.Keccak256(binary.AppendUvarint(address.Bytes(), nonce)))

	return &types.MemoDID{
		Method:      "memo",
		Identifier:  identifier,
		Identifiers: []string{identifier},
	}, nil
}

func (m *MemoDID) CreateDIDByAaddress(addressStr string) (*types.MemoDID, error) {
	_, endpoint := com.GetInsEndPointByChain(m.chain)

	client, err := ethclient.DialContext(context.TODO(), endpoint)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}
	defer client.Close()

	address := common.HexToAddress(addressStr)
	nonce, err := client.PendingNonceAt(context.TODO(), address)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	identifier := hex.EncodeToString(crypto.Keccak256(binary.AppendUvarint(address.Bytes(), nonce)))

	return &types.MemoDID{
		Method:      "memo",
		Identifier:  identifier,
		Identifiers: []string{identifier},
	}, nil
}

func (m *MemoDID) RegisterDIDByAddress(addressStr string, sig []byte) (string, error) {
	did, err := m.CreateDIDByAaddress(addressStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	address := common.HexToAddress(addressStr)
	err = m.Controller.RegisterDIDByAddress(did.Identifier, m.getMethodType("address"), address.Bytes(), sig)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}
	return did.String(), nil
}

func (m *MemoDID) RegisterDID(did, vtype string, publicKey, sig []byte) (string, error) {
	err := m.Controller.RegisterDID(did, m.getMethodType(vtype), publicKey, sig)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}
	return did, nil
}

func (m *MemoDID) GetDIDStatus() {

}

func (m *MemoDID) GetDIDInfo(address string) (*types.MemoDIDDocument, error) {
	did, err := m.CreateDIDByAaddress(address)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	return m.Controller.GetDIDInfo(did.String())
}

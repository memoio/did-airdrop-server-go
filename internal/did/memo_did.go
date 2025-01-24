package did

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/did-server/internal/contract"
	"github.com/did-server/internal/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	db         *database.DataBase
}

func NewMemoDID(chain string, logger *log.Helper) (*MemoDID, error) {
	controller, err := contract.NewController(chain, logger)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	db, err := database.CreateDB(logger)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &MemoDID{
		Controller: controller,
		chain:      chain,
		logger:     logger,
		db:         db,
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

func (m *MemoDID) CreateDIDByAddress(addressStr string) (*types.MemoDID, error) {
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

func (m *MemoDID) GetDIDNumber() (int, error) {
	num, err := m.db.GetNumber()
	if err != nil {
		m.logger.Error(err)
		return 0, err
	}

	return num, nil
}

func (m *MemoDID) AddDIDNumber(address string, num int) error {
	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	return m.db.AddNumber(did.String(), num)
}
func (m *MemoDID) RegisterDIDByAddress(addressStr string, sig []byte) (string, error) {
	did, err := m.CreateDIDByAddress(addressStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	address := common.HexToAddress(addressStr)

	num, err := m.db.GetNumber()
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	m.logger.Info("register did: ", did.String(), " number: ", num)

	err = m.Controller.RegisterDID(did.Identifier, m.getMethodType("address"), address.Bytes(), sig, big.NewInt(int64(num)))
	if err != nil {
		if strings.Contains(err.Error(), "existed") {
			return did.String(), nil
		}
		m.logger.Error(err)
		return "", err
	}

	err = m.db.AddNumber(did.String(), num)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return did.String(), nil
}

func (m *MemoDID) RegisterDIDByAddressByAdmin(addressStr string) (string, error) {
	did, err := m.CreateDIDByAddress(addressStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	address := common.HexToAddress(addressStr)

	num, err := m.db.GetNumber()
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	m.logger.Info("register did: ", did.String(), " number: ", num)

	err = m.Controller.RegisterDIDByAdmin(did.Identifier, m.getMethodType("address"), address.Bytes(), big.NewInt(int64(num)))
	if err != nil {
		if strings.Contains(err.Error(), "existed") {
			return did.String(), nil
		}
		m.logger.Error(err)
		return "", err
	}

	err = m.db.AddNumber(did.String(), num)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return did.String(), nil
}

func (m *MemoDID) RegisterDIDByTomAdmin(addressStr string) (string, error) {
	did, err := m.CreateDIDByAddress(addressStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	address := common.HexToAddress(addressStr)

	err = m.Controller.RegisterDIDByTonAdmin(did.Identifier, m.getMethodType("ton"), address.Bytes())
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return did.String(), nil
}

func (m *MemoDID) RegisterDIDByPublic(publicKeyStr string, sig []byte) (string, error) {
	did, err := m.CreateDIDByPubKey(publicKeyStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	publicKeyByte, err := hexutil.Decode(publicKeyStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	num, err := m.db.GetNumber()
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	err = m.Controller.RegisterDID(did.Identifier, m.getMethodType("pubkey"), publicKeyByte, sig, big.NewInt(int64(num)))
	if err != nil {
		m.logger.Error(err)
		return "", err
	}
	return did.String(), nil
}

func (m *MemoDID) GetDIDStatus() {

}

func (m *MemoDID) GetDIDInfo(address string) (string, string, error) {
	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return "", "", err
	}

	number, err := m.Controller.GetDIDInfo(did.Identifier)
	if err != nil {
		m.logger.Error(err)
		return "", "", err
	}

	return did.String(), number, nil
}

func (m *MemoDID) GetDIDExist(address string) (int, error) {
	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return 0, err
	}

	number, err := m.Controller.GetDIDVerify(did.Identifier)
	if err != nil {
		m.logger.Error(err)
		return 0, err
	}

	return number, nil
}

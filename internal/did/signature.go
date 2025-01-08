package did

import (
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (m *MemoDID) GetCreateSignatureMassageByPubKey(publickey string) (string, error) {
	did, err := m.CreateDIDByPubKey(publickey)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	nonce, err := m.Controller.GetNonce(did.String())
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return m.getCreateDIDHashPubkey(did.Identifier, publickey, nonce)
}

func (m *MemoDID) GetCreateSignatureMassageByAddress(address string) (string, error) {
	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	nonce, err := m.Controller.GetNonce(did.String())
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return m.CreateDIDMessageByAddress(did.Identifier, address, nonce)
}

func (m *MemoDID) VerifySign(sign, address string) (bool, error) {
	sig := hexutil.MustDecode(sign)
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}

	nonce, err := m.Controller.GetNonce(did.Identifier)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}
	hash, err := m.getCreateDIDHashByAddress(did.Identifier, address, nonce)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}
	hashB, err := hexutil.Decode(hash)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}

	signB, err := hexutil.Decode(sign)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}
	pubkey, err := crypto.SigToPub(hashB, signB)
	if err != nil {
		m.logger.Error(err)
		return false, err
	}

	addrP := crypto.PubkeyToAddress(*pubkey)
	if strings.Compare(addrP.Hex(), address) != 0 {
		return false, nil
	}
	return true, nil
}

func (m *MemoDID) GetDeleteSignatureMassage(did string) (string, error) {
	nonce, err := m.Controller.GetNonce(did)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return m.getDeleteDIDHash(did, nonce)
}

func (m *MemoDID) getCreateDIDHashPubkey(did, publickeyStr string, nonce uint64) (string, error) {
	tmp8 := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp8, nonce)

	pubKey, err := hexutil.Decode(publickeyStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	m.logger.Info(len(pubKey))

	createDID := []byte("createDID")
	didByte := []byte(did)
	method := []byte(m.getMethodType("pubkey"))

	hash := crypto.Keccak256(
		createDID,
		didByte,
		method,
		pubKey,
		tmp8,
	)
	return hexutil.Encode(hash), nil
}

func (m *MemoDID) getDeleteDIDHash(did string, nonce uint64) (string, error) {
	tmp8 := make([]byte, 8)

	binary.BigEndian.PutUint64(tmp8, nonce)

	deleteDID := []byte("deleteDID")
	didByte := []byte(did)
	deactivate := []byte{1}
	hash := crypto.Keccak256(
		deleteDID,
		didByte,
		deactivate,
		tmp8,
	)

	return hexutil.Encode(hash), nil
}

func (m *MemoDID) CreateDIDMessageByAddress(didI, addressStr string, nonce uint64) (string, error) {
	m.logger.Infof("didI = %s, addressStr = %s, nonce = %d", didI, addressStr, nonce)
	tmp8 := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp8, nonce)

	address := common.HexToAddress(addressStr)

	createDID := []byte("createDID")
	didByte := []byte(didI)
	method := []byte(m.getMethodType("address"))

	message := append(createDID, didByte...)
	message = append(message, method...)
	message = append(message, address.Bytes()...)
	message = append(message, tmp8...)

	return hexutil.Encode(message), nil
}

func (m *MemoDID) getCreateDIDHashByAddress(didI, addressStr string, nonce uint64) (string, error) {
	tmp8 := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp8, nonce)

	address := common.HexToAddress(addressStr)

	createDID := []byte("createDID")
	didByte := []byte(didI)
	method := []byte(m.getMethodType("address"))

	message := append(createDID, didByte...)
	message = append(message, method...)
	message = append(message, address.Bytes()...)
	message = append(message, tmp8...)

	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)

	hash := crypto.Keccak256(
		[]byte(msg),
	)
	return hexutil.Encode(hash), nil
}

func (m *MemoDID) getMethodType(mtype string) string {
	switch mtype {
	case "address":
		return "EcdsaSecp256k1RecoveryMethod2020"
	case "pubkey":
		return "EcdsaSecp256k1VerificationKey2019"
	default:
		return "EcdsaSecp256k1RecoveryMethod2020"
	}

}

func (m *MemoDID) publickeyFromString(publickey string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := hexutil.Decode(publickey)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	publicKeyECDSA, err := crypto.DecompressPubkey(pubKeyBytes)
	if err != nil {
		m.logger.Error(err)
		return nil, err
	}

	return publicKeyECDSA, nil
}

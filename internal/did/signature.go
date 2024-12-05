package did

import (
	"crypto/ecdsa"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (m *MemoDID) GetCreateSignatureMassage(vchain, publickey string) (string, error) {
	did, err := m.CreateDID(m.chain, publickey)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	nonce, err := m.Controller.GetNonce()
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return m.getCreateDIDHash(did.String(), vchain, publickey, nonce)
}

func (m *MemoDID) GetDeleteSignatureMassage(did string) (string, error) {
	nonce, err := m.Controller.GetNonce()
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	return m.getDeleteDIDHash(did, nonce)
}

func (m *MemoDID) getCreateDIDHash(did, vchain, publickeyStr string, nonce uint64) (string, error) {
	tmp8 := make([]byte, 8)
	binary.BigEndian.PutUint64(tmp8, nonce)

	pubKey, err := hexutil.Decode(publickeyStr)
	if err != nil {
		m.logger.Error(err)
		return "", err
	}

	createDID := common.LeftPadBytes([]byte("createDID"), 32)
	didByte := common.LeftPadBytes([]byte(did), 32)
	method := common.LeftPadBytes([]byte(m.getMethodType(vchain)), 32)
	pubKeyBytes := common.LeftPadBytes(pubKey, 32)

	hash := crypto.Keccak256(
		createDID,
		didByte,
		method,
		pubKeyBytes,
		tmp8,
	)
	return hexutil.Encode(hash), nil
}

func (m *MemoDID) getDeleteDIDHash(did string, nonce uint64) (string, error) {
	tmp8 := make([]byte, 8)
	deactivate := make([]byte, 32)
	binary.BigEndian.PutUint64(tmp8, nonce)

	deleteDID := common.LeftPadBytes([]byte("deleteDID"), 32)
	didByte := common.LeftPadBytes([]byte(did), 32)

	hash := crypto.Keccak256(
		deleteDID,
		didByte,
		deactivate,
		tmp8,
	)

	return hexutil.Encode(hash), nil
}

func (m *MemoDID) getMethodType(vtype string) string {
	return "EcdsaSecp256k1VerificationKey2019"
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

package did

import (
	"math/big"

	"github.com/memoio/go-did/types"
)

func (m *MemoDID) CreateMfileDID(cid string) (*types.MfileDID, error) {
	return &types.MfileDID{
		Method:     "mfile",
		HashMethod: "cid",
		Identifier: cid,
	}, nil
}

func (m *MemoDID) CreateMfileInfo(address, cid string, price *big.Int) error {
	did, err := m.CreateDIDByAddress(address)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	mfile, err := m.CreateMfileDID(cid)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	return m.db.CreateMfileInfo(address, did.String(), mfile.String(), price)
} 

func (m *MemoDID) RegisterMfileDID(mdidString string, sig []byte) error {
	minfo, err := m.db.GetMfileInfo(mdidString)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	did, err := types.ParseMemoDID(minfo.DID)
	if err != nil {
		m.logger.Error(err)
		return err
	}
	mfile, err := types.ParseMfileDID(mdidString)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	err = m.Controller.RegisterMfile(mfile.Identifier, did.Identifier, minfo.Price, sig)
	if err != nil {
		m.logger.Error(err)
		return err
	}

	return nil
}

func ParaseMfileDID(didString string) (*types.MfileDID, error) {
	return types.ParseMfileDID(didString)
}

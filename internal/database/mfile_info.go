package database

import (
	"math/big"

	"gorm.io/gorm"
)

type MfileInfo struct {
	gorm.Model
	Address string
	DID     string
	MDID    string `gorm:"uniqueIndex:mfile_composite;"`
	Price   *big.Int
}

func (d *DataBase) CreateMfileInfo(address, did, mdid string, price *big.Int) error {
	result := d.db.Create(&MfileInfo{Address: address, DID: did, MDID: mdid, Price: price})
	if result.Error != nil {
		err := result.Error
		d.logger.Error(err)
		return err
	}

	return nil
}

func (d *DataBase) GetMfileInfo(mdid string) (*MfileInfo, error) {
	var info MfileInfo
	result := d.db.Where("mdid = ?", mdid).First(&info)
	if result.Error != nil {
		err := result.Error
		d.logger.Error(err)
		return nil, err
	}
	return &info, nil
}

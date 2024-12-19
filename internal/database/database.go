package database

import (
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DataBase struct {
	db     *gorm.DB
	logger *log.Helper
}

type Number struct {
	gorm.Model
	Did string `gorm:"unique"`
	Num int
}

func CreateDB(logger *log.Helper) (*DataBase, error) {
	db, err := gorm.Open(sqlite.Open("numbers.db"), &gorm.Config{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	db.AutoMigrate(&Number{})
	return &DataBase{db: db, logger: logger}, nil
}

func (d *DataBase) GetNumber() (int, error) {
	var maxNumber Number
	err := d.db.Order("num desc").Limit(1).First(&maxNumber).Error

	// 如果查询没有结果，返回默认值 100001
	if err != nil || maxNumber.ID == 0 {
		d.logger.Info(maxNumber.ID)
		maxNumber.Num = 100001
		return maxNumber.Num, nil
	}

	return maxNumber.Num + 1, nil
}

func (d *DataBase) AddNumber(did string, num int) error {
	result := d.db.Create(&Number{Did: did, Num: num})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

package models

type BlockScan struct {
	Id          int64 `gorm:"type:int(11) UNSIGNED AUTO_INCREMENT;primary_key" json:"id"`
	BlockNumber int64 `gorm:"type:int(64) UNSIGNED not null COMMENT '同步的区块高度'" json:"block_number"`
}

func (b BlockScan) Create(blockScan BlockScan) error {
	return db.Create(&blockScan).Error
}
func (b *BlockScan) GetNumber() int64 {
	var bscScan BlockScan
	err := db.Select("*").Where("id >0").First(&bscScan).Error
	if err != nil {
		return 0
	}
	return bscScan.BlockNumber
}

func (b *BlockScan) Edit(data map[string]interface{}) error {
	return db.Model(&b).Updates(data).Error
}

package models


// import (
// 	"errors"
// 	"github.com/jinzhu/gorm"
// )

type TradeHistory struct {
	ID 		uint `gorm:"primary_key"`
	Hash    string `gorm:"column:hash"`
	From    string `gorm:"column:from"`
	To      string `gorm:"column:to"`
	Amount  uint64 `gorm:"column:amount"`
}


func (t TradeHistory) CreateUser(tradeHistory TradeHistory) error {
	return db.Create(&tradeHistory).Error
}

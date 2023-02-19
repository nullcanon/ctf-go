package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type UserTable struct {
	Timestamp      uint64 `gorm:"column:timestamp"`
	TotalReward    string `gorm:"column:total_reward; default:'0'"`
	ReceivedReward string `gorm:"column:received_reward; default:'0'"`
	LpRewards      string `gorm:"column:lp_rewards; default:'0'"`
	TradeVolume    string `gorm:"column:trade_volume; default:'0'"`
	TradeVolLowers string `gorm:"column:trade_vol_lowers; default:'0'"`
	TradeVolAll    string `gorm:"column:trade_vol_all; default:'0'"`
	Role           int    `gorm:"column:role"`
	Upper          string `gorm:"column:upper; default:'0'"`
	Self           string `gorm:"column:self; default:'0'"`
}

func (u UserTable) CreateUser(userinfo UserTable) error {
	return db.Create(&userinfo).Error
}

func (u UserTable) UpdateUserInv() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update("upper", u.Upper)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

func (u UserTable) UpdateTradeVolume() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update(map[string]interface{}{"trade_volume": u.TradeVolume, "role": u.Role})
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

func (u UserTable) UpdateTradeVolLowers() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update(map[string]interface{}{"trade_vol_all": u.TradeVolAll})
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

func (u UserTable) UpdateTotalReward() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update(map[string]interface{}{"total_reward": u.TotalReward})
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

func (u UserTable) UpdateTradeVolAll() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update(map[string]interface{}{"total_reward": u.TradeVolume, "role": u.Role})
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

// func (u UserTable) UpdateTotalReward(userinfo UserTable) error {
// }

// func (u UserTable) UpdateReceivedReward(userinfo UserTable) error {
// }

func (u UserTable) UpdateLpRewards() error {
	var userinfo UserTable
	result := db.First(&userinfo, "self = ?", u.Self)
	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", u.Self).Update("lp_rewards", u.LpRewards)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(u)
	} else {
		return result.Error
	}
}

// func (u UserTable) FirstOrCreateUser(self string, userinfo UserTable) error {
// 	return db.Where(User{Name: "new_name"}).Attrs(User{Age: 18}).FirstOrCreate(&user)
// }

func (u UserTable) FetchUserInfo(userinfo *[]UserTable) {
	// db.Table("user_table").Where("self = ?", "your_self").First(&userinfo)
	db.Find(&userinfo)
}

type UserLowersTable struct {
	Self  string `gorm:"column:self"`
	Lower string `gorm:"column:lowers"`
}

func (u UserLowersTable) CreateUserLowers(userlowers UserLowersTable) error {
	return db.Create(&userlowers).Error
}

func (u UserLowersTable) FatchLowers(self string, lowers *[]string) {
	db.Table("user_lowers_table").Where("self = ?", self).Pluck("lowers", lowers)
}

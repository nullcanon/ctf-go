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

func (u UserTable) UpdateUserInv(userinfo UserTable) error {
	result := db.First(&userinfo, "self = ?", userinfo.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", userinfo.Self).Update("upper", userinfo.Upper)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(userinfo)
	} else {
		return result.Error
	}
	return nil
}

func (u UserTable) UpdateTradeVolume(userinfo UserTable) error {
	result := db.First(&userinfo, "self = ?", userinfo.Self)

	if result.Error == nil {
		db.Model(&UserTable{}).Where("self = ?", userinfo.Self).Update("upper", userinfo.TradeVolume)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u.CreateUser(userinfo)
	} else {
		return result.Error
	}
	return nil
}

// func (u UserTable) UpdateTotalReward(userinfo UserTable) error {
// }

// func (u UserTable) UpdateReceivedReward(userinfo UserTable) error {
// }

// func (u UserTable) UpdateLpReward(userinfo UserTable) error {
// }

// func (u UserTable) FirstOrCreateUser(self string, userinfo UserTable) error {
// 	return db.Where(User{Name: "new_name"}).Attrs(User{Age: 18}).FirstOrCreate(&user)
// }

func (u UserTable) FetchUserInfo(userinfo *[]UserTable) {
	// db.Table("user_table").Where("self = ?", "your_self").First(&userinfo)
	db.Find(&userinfo)
}

// TODO db.Find() 具体参考 chatgpt
// var users []User
//     db.Find(&users)

//     for _, user := range users {
//         inviter.Userinfos[user.Upper] = &user
//     }

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

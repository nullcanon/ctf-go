package models


import (
	"errors"
	"github.com/jinzhu/gorm"
)

type UserTable struct {
	Timestamp      uint64 `gorm:"column:timestamp"`
	TotalReward    uint64 `gorm:"column:total_reward"`
	ReceivedReward uint64 `gorm:"column:received_reward"`
	LpReward       uint64 `gorm:"column:lp_reward"`
	TradeVolume    uint64 `gorm:"column:trade_volume"`
	Upper          string `gorm:"column:upper"`
	Self           string `gorm:"column:self"`
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

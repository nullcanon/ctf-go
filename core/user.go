package core

import (
	//   "github.com/kanocz/goginjsonrpc"
	//   "github.com/gin-gonic/gin"
	"ctf/models"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	// "hash"
	"time"

	common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	// crypto1 "github.com/ubiq/go-ubiq/crypto"
)

// TODO 改为User模块
// TODO 优化 mysql 数据插入为异步

var InviterHandle *Inviter

var selfRate = decimal.NewFromFloat(0.2)
var lowerRate = decimal.NewFromFloat(0.5)
var treeRate = decimal.NewFromFloat(0.3)
var feeRate = decimal.NewFromFloat(0.01)
var subCoinPrice = decimal.NewFromFloat(0.158)

type Trole int

const (
	NOMAL       = 0
	PERSONAL    = 1
	TEAM_LEADER = 2
)

type User struct {
	timestamp      uint64
	totalReward    decimal.Decimal // 总奖励
	receivedReward decimal.Decimal // 已领取奖励
	tradeVolume    decimal.Decimal // 个人交易量
	lpRewards      decimal.Decimal // lp奖励
	tradeVolLowers decimal.Decimal // 直推交易量
	tradeVolAll    decimal.Decimal // 伞下所有交易量
	upper          string
	self           string
	role           Trole
	lowers         []string
	// nonce uint64
}

type Inviter struct {
	userinfos map[string]*User
	mutex     sync.Mutex
}

func NewInviter() (*Inviter, error) {
	inviter := &Inviter{
		userinfos: make(map[string]*User),
	}
	// Recovering data from a database
	if err := inviter.fatchData(); err != nil {
		return nil, err
	}
	return inviter, nil
}

func (t *Inviter) fatchData() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var users []models.UserTable
	userinfo := models.UserTable{}
	userinfo.FetchUserInfo(&users)
	for _, user := range users {
		var lowers []string
		userlower := models.UserLowersTable{}
		userlower.FatchLowers(user.Self, &lowers)
		t.userinfos[user.Self] = &User{
			timestamp:      user.Timestamp,
			totalReward:    decimal.RequireFromString(user.TotalReward),
			receivedReward: decimal.RequireFromString(user.ReceivedReward),
			tradeVolume:    decimal.RequireFromString(user.TradeVolume),
			lpRewards:      decimal.RequireFromString(user.LpRewards),
			tradeVolLowers: decimal.RequireFromString(user.TradeVolLowers),
			tradeVolAll:    decimal.RequireFromString(user.TradeVolAll),
			upper:          user.Upper,
			self:           user.Self,
			lowers:         lowers,
		}

		if _, ok := t.userinfos[user.Upper]; !ok {
			var ulowers []string
			userlower.FatchLowers(user.Upper, &ulowers)
			t.userinfos[user.Upper] = &User{
				timestamp:      0,
				totalReward:    decimal.NewFromInt(0),
				receivedReward: decimal.NewFromInt(0),
				tradeVolume:    decimal.NewFromInt(0),
				lpRewards:      decimal.NewFromInt(0),
				tradeVolLowers: decimal.NewFromInt(0),
				tradeVolAll:    decimal.NewFromInt(0),
				upper:          "",
				self:           user.Upper,
				lowers:         ulowers,
			}
		}
	}

	return nil
}

func RecoverPublicKeyAddress(data string, signature string) (string, error) {
	hexdata := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)

	hash := crypto.Keccak256Hash([]byte(hexdata))

	hexmessage, err := hex.DecodeString(signature[2:])
	if err != nil {
		return "", err
	}
	hexmessage[64] -= 27

	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), hexmessage)
	if err != nil {
		return "", err
	}

	address := crypto.PubkeyToAddress(*sigPublicKeyECDSA)

	return address.String(), nil
}

func (t *Inviter) BindInvCode(upper string, user string, singerMessage string) (string, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if !common.IsHexAddress(upper) {
		return "", errors.New("Upper not ethereum address")
	}

	upperChecksum := common.HexToAddress(upper).Hex()

	// Verify the signature message
	sigPublicKey, err := RecoverPublicKeyAddress(upperChecksum, singerMessage)
	if err != nil {
		fmt.Println("sigPublicKey error:", err.Error())
		return "", errors.New("SingerMessage error")
	}
	userChecksum := common.HexToAddress(sigPublicKey).Hex()

	if common.HexToAddress(user).Hex() != userChecksum {
		return "", errors.New("SingerMessage not match user")
	}

	if userChecksum == upperChecksum {
		return "", errors.New("Can't invite yourself")
	}

	// Chek repeated invitations
	for tmpUser := userChecksum; ; {
		// TODO 互相绑定时这一行崩溃
		if u, ok := t.userinfos[tmpUser]; ok {
			if u.upper == userChecksum {
				return "", errors.New("Repeated invitations")
			}
			tmpUser = u.upper
		} else {
			break
		}
	}

	var (
		userInfo User
	)
	if userInfo, ok := t.userinfos[userChecksum]; ok {
		if userInfo.upper != "" {
			return "", errors.New("Upper is exist")
		}
	}
	time := uint64(time.Now().UnixMilli())
	userInfo.upper = upperChecksum
	userInfo.self = userChecksum
	userInfo.timestamp = time
	t.userinfos[userChecksum] = &userInfo

	if upperInfo, ok := t.userinfos[upperChecksum]; ok {
		upperInfo.lowers = append(upperInfo.lowers, userChecksum)
	} else {
		upperInfo = &User{
			lowers: []string{userChecksum},
		}
		userInfo.self = upperChecksum
		t.userinfos[upperChecksum] = upperInfo
	}

	// up to database
	usertable := models.UserTable{
		Timestamp:      time,
		TotalReward:    "0",
		ReceivedReward: "0",
		Upper:          upperChecksum,
		Self:           userChecksum,
	}
	usertable.UpdateUserInv()

	userlowers := models.UserLowersTable{
		Self:  upperChecksum,
		Lower: userChecksum,
	}
	userlowers.CreateUserLowers(userlowers)
	return "", nil
}

// 可领取子币数量= 返利类型对应用户的币种交易金额 * 1% * 不同类型返佣比例 / 0.158
// 0.158为子币开盘价格
// 个人返佣比例：20%
// 直推返佣比例：50%
// 社区节点返佣比例：30%

func getSubCoinAmount(tradeAmount decimal.Decimal, rate decimal.Decimal) decimal.Decimal {
	return tradeAmount.Mul(feeRate).Mul(rate).Div(subCoinPrice)
}

func (t *Inviter) UpdateTradeVolume(user string, amount decimal.Decimal) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	var userInfo *User
	userChecksum := common.HexToAddress(user).Hex()
	logrus.Infof("UpdateTradeVolume")
	userInfo, ok := t.userinfos[userChecksum]
	if ok {
		// 先更新自己的交易量
		userInfo.tradeVolume = userInfo.tradeVolume.Add(amount)
		userInfo.totalReward = getSubCoinAmount(amount, selfRate)
		// 如果自己交易量达到3000刀，更新角色为 PERSONAL
		if userInfo.tradeVolume.GreaterThanOrEqual(decimal.NewFromInt(3000)) {
			userInfo.role = PERSONAL
		}

		if userInfo.upper != "" {
			// 更新上级的直推奖励 tradeVolLowers
			if upperInfo, ok := t.userinfos[userInfo.upper]; ok {
				upperInfo.tradeVolLowers = upperInfo.tradeVolLowers.Add(amount)
				upperInfo.totalReward = getSubCoinAmount(amount, lowerRate)
				if upperInfo.tradeVolLowers.GreaterThanOrEqual(decimal.NewFromInt(30000)) {
					upperInfo.role = TEAM_LEADER
				}

				// 更新数据库
				uppertable := models.UserTable{
					Self:        upperInfo.self,
					TradeVolume: upperInfo.tradeVolLowers.String(),
					Role:        int(upperInfo.role),
				}
				uppertable.UpdateTradeVolLowers()

				// 更新上级链条的伞下收益
				// TODO 一笔交易量计入直推奖励的同时，计入散装奖励吗，这里的做法暂时计入
				tmpAddress := userInfo.upper
				for true {
					info := t.userinfos[tmpAddress]
					if info.upper == "" || info.role == TEAM_LEADER {
						break
					}
					info.tradeVolAll = info.tradeVolAll.Add(amount)
					info.totalReward = getSubCoinAmount(amount, treeRate)
					tmpAddress = info.upper

					tmpuppertable := models.UserTable{
						Self:        info.self,
						TradeVolAll: info.tradeVolAll.String(),
					}
					tmpuppertable.UpdateTradeVolAll()
				}
				// 如果上面todo不计入 upperInfo.tradeVolAll -= amount
			}
		}

	} else {
		// 只更新自己的交易量
		userInfo = &User{
			self:        userChecksum,
			tradeVolume: amount,
		}
		userInfo.totalReward = getSubCoinAmount(amount, selfRate)
		t.userinfos[userChecksum] = userInfo
	}
	// 更新数据库
	usertable := models.UserTable{
		Self:        userChecksum,
		TradeVolume: userInfo.tradeVolume.String(),
		Role:        int(userInfo.role),
	}
	return usertable.UpdateTradeVolume()

}

func (t *Inviter) UpdateLpRewards(user string, amount decimal.Decimal) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	userChecksum := common.HexToAddress(user).Hex()
	logrus.Infof("UpdateLpRewards")
	logrus.Infof("UpdateLpRewards amount %v", amount.String())
	userInfo, ok := t.userinfos[userChecksum]
	if ok {
		// 更新自己的lp奖励
		userInfo.lpRewards = userInfo.lpRewards.Add(amount)
		logrus.Infof("UpdateLpRewards update amount %v", userInfo.lpRewards.String())
	} else {
		userInfo = &User{
			self:      userChecksum,
			lpRewards: amount,
		}
		t.userinfos[userChecksum] = userInfo
	}
	// 更新数据库
	logrus.Infof("UpdateLpRewards %v", userInfo.lpRewards.String())
	usertable := models.UserTable{
		Self:      userChecksum,
		LpRewards: userInfo.lpRewards.String(),
	}
	logrus.Infof("usertable %v", usertable.LpRewards)

	return usertable.UpdateLpRewards()
}

type Users struct {
	Users []string `json:"users"`
}

func (t *Inviter) ProcessPresellUsersRewards() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	// 1、读取json用户列表
	jsonFile, err := os.Open("./files/presell_users.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}

	var users Users
	if err := json.Unmarshal(byteValue, &users); err != nil {
		fmt.Println(err)
	}

	// 2、为json列表用户的upper更新信息
	for _, user := range users.Users {
		checkSumUser := common.HexToAddress(user).String()
		if userInfo, ok := t.userinfos[checkSumUser]; ok {
			if upperInfo, ok := t.userinfos[userInfo.upper]; ok {
				upperInfo.totalReward = upperInfo.totalReward.Add(decimal.NewFromInt(200))
				// 3、记录数据库
				usertable := models.UserTable{
					Self:        userInfo.upper,
					TotalReward: upperInfo.totalReward.String(),
				}
				usertable.UpdateTotalReward()
			}

		}
		fmt.Println(user)
	}

}

func (t *Inviter) GetTotalRewardRank(offset uint64, limit uint64) {

}

// 添加排序代码
// 参考
// type User struct {
// 	ID       int
// 	Name     string
// 	Score    int
//   }

//   func main() {
// 	users := []User{
// 	  {ID: 1, Name: "Tom", Score: 90},
// 	  {ID: 2, Name: "John", Score: 80},
// 	  {ID: 3, Name: "Jane", Score: 100},
// 	}

// 	sort.Slice(users, func(i, j int) bool {
// 	  return users[i].Score > users[j].Score
// 	})

// 	for i, user := range users {
// 	  fmt.Printf("Rank %d: %s, Score: %d\n", i+1, user.Name, user.Score)
// 	}

type ReferalsData struct {
	Address string
	Time    uint64
	Volume  decimal.Decimal
}

func (t *Inviter) GetLowers(address string, offset uint64, limit uint64) ([]ReferalsData, uint64, error) {
	userChecksum := common.HexToAddress(address).Hex()
	user, ok := t.userinfos[userChecksum]
	if !ok {
		return nil, 0, fmt.Errorf("user %s not found", userChecksum)
	}
	total := uint64(len(user.lowers))
	if offset >= total {
		return nil, 0, fmt.Errorf("offset out of range")
	}
	end := offset + limit
	if end > total {
		end = total
	}

	rd := []ReferalsData{}

	for _, value := range user.lowers[offset:end] {
		if lowuser, ok := t.userinfos[value]; ok {
			rd = append(rd, ReferalsData{
				Address: value,
				Time:    lowuser.timestamp,
				Volume:  lowuser.tradeVolume,
			})
		}
	}

	return rd, total, nil
}

func (t *Inviter) GetMyUpper(lower string) (string, error) {
	user, ok := t.userinfos[lower]
	if !ok {
		return "", nil
	}
	return user.upper, nil
}

func (t *Inviter) GetLowersAmount(address string) uint64 {
	userChecksum := common.HexToAddress(address).Hex()
	if upperInfo, ok := t.userinfos[userChecksum]; ok {
		return uint64(len(upperInfo.lowers))
	}
	return 0
}

// func (t *Inviter) GetTotalLpRewards() uint64 {

// }

// func (t *Inviter) GetLpRewardsRank(offset uint64, limit uint64) (uint64, []string) {
//	total :=
// }

// func (t *Inviter) getReLowersTradeVolume(address string) uint {

// }

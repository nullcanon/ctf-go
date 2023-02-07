package core

import (
	//   "github.com/kanocz/goginjsonrpc"
	//   "github.com/gin-gonic/gin"
	"ctf/models"
	"encoding/hex"
	"errors"
	"fmt"

	// "hash"
	"time"

	common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	// crypto1 "github.com/ubiq/go-ubiq/crypto"
)

// TODO 改为User模块
// TODO 优化 mysql 数据插入为异步

var InviterHandle *Inviter

type Trole int

const (
	NOMAL       = 0
	PERSONAL    = 1
	TEAM_LEADER = 2
)

type User struct {
	timestamp      uint64
	totalReward    uint64
	receivedReward uint64
	tradeVolume    uint64
	upper          string
	self           string
	Role           Trole
	lowers         []string
	// nonce uint64
}

type Inviter struct {
	userinfos map[string]*User
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
	var users []models.UserTable
	userinfo := models.UserTable{}
	userinfo.FetchUserInfo(&users)
	for _, user := range users {
		var lowers []string
		userlower := models.UserLowersTable{}
		userlower.FatchLowers(user.Self, &lowers)
		fmt.Println("fatchData user", user.Self)
		t.userinfos[user.Self] = &User{
			timestamp:      user.Timestamp,
			totalReward:    user.TotalReward,
			receivedReward: user.ReceivedReward,
			tradeVolume:    user.TradeVolume,
			upper:          user.Upper,
			self:           user.Self,
			lowers:         lowers,
		}

		if _, ok := t.userinfos[user.Upper]; !ok {
			var ulowers []string
			userlower.FatchLowers(user.Upper, &ulowers)
			t.userinfos[user.Upper] = &User{
				timestamp:      0,
				totalReward:    0,
				receivedReward: 0,
				tradeVolume:    user.TradeVolume,
				upper:          "",
				self:           user.Upper,
				lowers:         ulowers,
			}
		}
	}

	for self, user := range t.userinfos {
		fmt.Printf("Self: %s, Timestamp: %d, Total Reward: %d, Received Reward: %d, Upper: %s, s: %s, Lowers: %v\n",
			self, user.timestamp, user.totalReward, user.receivedReward, user.upper, user.self, user.lowers)
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
		TotalReward:    0,
		ReceivedReward: 0,
		Upper:          upperChecksum,
		Self:           userChecksum,
	}
	usertable.UpdateUserInv(usertable)

	userlowers := models.UserLowersTable{
		Self:  upperChecksum,
		Lower: userChecksum,
	}
	userlowers.CreateUserLowers(userlowers)
	return "", nil
}

type ReferalsData struct {
	Address string
	Time    uint64
	Volume  uint64
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

// func (t *Inviter) getReLowersTradeVolume(address string) uint {

// }

package api

import (
	"ctf/core"
	"ctf/models"
)

type Api struct {
}

func (a *Api) BindInvCode(upper string, user string, singerMessage string) (string, error) {
	return core.InviterHandle.BindInvCode(upper, user, singerMessage)
}

type MyInfo struct {
	PersonalTime uint64
	CommuniyTime uint64
	UpperAddress string
	TradeVolume  string
}

func (a *Api) GetMyInfo(user string) (MyInfo, error) {
	upper, err := core.InviterHandle.GetMyUpper(user)
	if err != nil {
		return MyInfo{}, err
	}
	volume, _ := core.InviterHandle.GetTradeVolumeAndRewards(user)

	return MyInfo{
		TradeVolume:  volume,
		PersonalTime: 0,
		CommuniyTime: 0,
		UpperAddress: upper,
	}, nil
}

type InvData struct {
	InvAmount      uint64
	LowTradeAmount uint64
	Received       string
	Available      string
	LowTradeVolume string
}

func (a *Api) GetInvData(user string, role int) (InvData, error) {
	invAmount := core.InviterHandle.GetLowersAmount(user)
	volume, commission := core.InviterHandle.GetTradeVolumeAndRewards(user)
	lowerTradeAmount := core.InviterHandle.GetLowersTradeAmount(user)
	return InvData{
		Received:       "0",
		Available:      commission,
		InvAmount:      invAmount,
		LowTradeAmount: lowerTradeAmount,
		LowTradeVolume: volume,
	}, nil
}

type RankDetail struct {
	Id      uint64
	Address string
	Amount  string
}

type TradingData struct {
	Address    string
	Volume     string
	Commission string
}

type Commission struct {
	Time          uint64
	ExtractAmount string
}

type DataDetail struct {
	Total   uint64
	Offset  uint64
	Limit   uint64
	Records interface{}
}

func (a *Api) GetInvDataDetail(user string, role int, sub int, offset uint64, limit uint64) (DataDetail, error) {
	var (
		records    interface{}
		total      uint64
		err        error
		volume     string
		commission string
	)
	if sub == 0 {
		var data []core.ReferalsData
		data, total, err = core.InviterHandle.GetLowers(user, offset, limit)
		if err != nil {
			return DataDetail{}, err
		}
		records = data
	} else if sub == 1 {
		if role == 0 {
			volume, commission = core.InviterHandle.GetLowersTradeVolumeAndRewards(user)
		} else {
			volume, commission = core.InviterHandle.GetLowersTreeTradeVolumeAndRewards(user)
		}
		records = []TradingData{
			{
				Address:    "",
				Volume:     volume,
				Commission: commission,
			},
		}
	} else {
		records = []Commission{
			{
				Time:          0,
				ExtractAmount: "0",
			},
		}
	}

	return DataDetail{
		Total:   total,
		Offset:  offset,
		Limit:   limit,
		Records: records,
	}, nil
}

func (a *Api) GetNodesRewards() (uint64, error) {
	return 0, nil
}

type BuyBackDataDetail struct {
	Id        uint64
	Timestamp uint64
	Amount    string
	Link      string
}

func (a *Api) GetBuyBackDetil() ([]BuyBackDataDetail, error) {
	return []BuyBackDataDetail{
		{
			Id:        1,
			Timestamp: 1676903961000,
			Amount:    "4398.42",
			Link:      "https://bscscan.com/tx/0x23f28d6a51e7214319186723bdce4fe9bfd0c7eb51f40e01a433990948280ae3",
		},
		{
			Id:		  2,
			Timestamp: 1677210923000,
			Amount:    "1284.64",
			Link:      "https://bscscan.com/tx/0xd13e3ea91e13a8b9a39d1feab52df93cea7b3ff666299ec1dc8e8f594a8c4bbb",
		},
		{
			Id:		  3,
			Timestamp: 1678356323000,
			Amount:    "845.48",
			Link:      "https://bscscan.com/tx/0xd9c747d16174309a519bac8b2b87a4a97fcd4ce7d044039f209ef4a7e285d756",
		},
		{
			Id:		  4,
			Timestamp: 1678696873000,
			Amount:    "573.34",
			Link:      "https://bscscan.com/tx/0xe70d6db39f48c98d0e38653e5ae1a1ef623772c3e455b6e8cd4892c01b2df313",
		}
	}, nil
}

func (a *Api) GetLpAwardedBonus() (string, error) {
	return models.UserTable{}.GetLpRewradsTotal()
}

func (a *Api) GetLpBonusRank(offset uint64, limit uint64) (DataDetail, error) {

	sortedData, total, _ := models.UserTable{}.GetLpRewardsRank(offset, limit)

	rankDetails := make([]RankDetail, len(sortedData))
	for i, data := range sortedData {
		rankDetails[i] = RankDetail{
			Id:      offset + uint64(i+1), // 排名从 1 开始
			Address: data.Self,
			Amount:  data.LpRewards,
		}
	}

	return DataDetail{
		Total:   uint64(total),
		Offset:  offset,
		Limit:   limit,
		Records: rankDetails,
	}, nil
}

func (a *Api) GetCTFCoinTradeVolume() (string, error) {
	return models.UserTable{}.GetTradeVolumeTotal()
}

func (a *Api) GetACoinRewardRank(offset uint64, limit uint64) (DataDetail, error) {

	sortedData, total, err := models.UserTable{}.GetTotalRewardRanks(offset, limit)

	rankDetails := make([]RankDetail, len(sortedData))
	for i, data := range sortedData {
		rankDetails[i] = RankDetail{
			Id:      offset + uint64(i+1), // 排名从 1 开始
			Address: data.Self,
			Amount:  data.TotalReward,
		}
	}

	return DataDetail{
		Total:   uint64(total),
		Offset:  offset,
		Limit:   limit,
		Records: rankDetails,
	}, err
}

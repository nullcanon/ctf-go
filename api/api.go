package api

import (
	"ctf/core"
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

	return MyInfo{
		TradeVolume:  "0",
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
	return InvData{
		Received:       "0",
		Available:      "0",
		InvAmount:      invAmount,
		LowTradeAmount: 0,
		LowTradeVolume: "0",
	}, nil
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
		records interface{}
		total   uint64
		err     error
	)
	if sub == 0 {
		var data []core.ReferalsData
		data, total, err = core.InviterHandle.GetLowers(user, offset, limit)
		if err != nil {
			return DataDetail{}, err
		}
		records = data
	} else if sub == 1 {
		records = []TradingData{
			TradingData{
				Address:    "",
				Volume:     "0",
				Commission: "0",
			},
		}
	} else {
		records = []Commission{
			Commission{
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
	id        uint64
	timestamp uint64
	amount    string
	link      string
}

func (a *Api) GetBuyBackDetil() ([]BuyBackDataDetail, error) {
	return []BuyBackDataDetail{
		BuyBackDataDetail{
			id:        0,
			timestamp: 0,
			amount:    "0",
			link:      "https://bscscan.com/tx/0x0e2491b362b750edac36c61a215425d18d7ac8d37c9373bd4110f1005ab046a2",
		},
	}, nil
}

func (a *Api) GetLpAwardedBonus() (string, error) {
	return "0", nil
}

type RankDetail struct {
	id      uint64
	address string
	amount  string
}

func (a *Api) GetLpBonusRank(offset uint64, limit uint64) (DataDetail, error) {
	records := []RankDetail{
		RankDetail{
			id:      0,
			address: " ",
			amount:  "0",
		},
	}

	return DataDetail{
		Total:   0,
		Offset:  offset,
		Limit:   limit,
		Records: records,
	}, nil
}

func (a *Api) GetCTFCoinTradeVolume() (string, error) {
	return "332849.32", nil
}

func (a *Api) GetACoinRewardRank(offset uint64, limit uint64) (DataDetail, error) {
	records := []RankDetail{
		RankDetail{
			id:      0,
			address: "",
			amount:  "0",
		},
	}

	return DataDetail{
		Total:   0,
		Offset:  offset,
		Limit:   limit,
		Records: records,
	}, nil
}

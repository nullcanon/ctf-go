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
	TradeVolume  uint64
	PersonalTime uint64
	CommuniyTime uint64
	UpperAddress string
}

func (a *Api) GetMyInfo(user string) (MyInfo, error) {
	upper, err := core.InviterHandle.GetMyUpper(user)
	if err != nil {
		return MyInfo{}, err
	}

	return MyInfo{
		TradeVolume:  0,
		PersonalTime: 0,
		CommuniyTime: 0,
		UpperAddress: upper,
	}, nil
}

type InvData struct {
	Received       uint64
	Available      uint64
	InvAmount      uint64
	LowTradeAmount uint64
	LowTradeVolume uint64
}

func (a *Api) GetInvData(user string, role int) (InvData, error) {
	return InvData{}, nil
}

type ReferalsData struct {
	Address string
	Time    uint64
	Volume  uint64
}

type TradingData struct {
	Address    string
	Volume     uint64
	Commission uint64
}

type Commission struct {
	Time          uint64
	ExtractAmount uint64
}

type DataDetail struct {
	Total   uint64
	Offset  uint64
	Limit   uint64
	Records interface{}
}

func (a *Api) GetInvDataDetail(user string, role int, sub int, offset uint64, limit uint64) (DataDetail, error) {
	var records interface{}
	if sub == 0 {
		records = []ReferalsData{
			ReferalsData{
				Address: "",
				Time:    0,
				Volume:  0,
			},
		}
	} else if sub == 1 {
		records = []TradingData{
			TradingData{
				Address:    "",
				Volume:     0,
				Commission: 0,
			},
		}
	} else {
		records = []Commission{
			Commission{
				Time:          0,
				ExtractAmount: 0,
			},
		}
	}

	return DataDetail{
		Total:   0,
		Offset:  offset,
		Limit:   limit,
		Records: records,
	}, nil
}

func (a *Api) GetNodesRewards() (uint64, error) {
	return 0, nil
}

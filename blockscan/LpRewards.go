// 收到lp奖励
package blockscan

import (
	"context"
	"ctf/core"
	"ctf/utils"

	ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
	"github.com/HydroProtocol/ethereum-watcher/blockchain"
	"github.com/HydroProtocol/ethereum-watcher/plugin"
	"github.com/HydroProtocol/ethereum-watcher/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

// const (
// 	CTF_CONTRACT = "0x29a3DAa1bf8DE08a8afE3e29E36Fa2797f7a37b5"
// 	USDT_BLACK = "0xb5151A092593ca413fC8782Cd92904F4f3d5F2f8" //代币合约中间地址
// )

var ReceiveLpFee = []string{"0xae84679ccd369379ebf43f95e15856a8f3ca5d85bd241d4d4426c0dd60e68693"}

// var ethrpc *rpc.EthBlockChainRPC

func lpRewardsHandle(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {
	logrus.Infof("recv len: %v", len(receiptLogs))

	for _, log := range receiptLogs {

		logrus.Infof("recv BlockNum: %v", log.GetBlockNum())
		logrus.Infof("hash: %s", log.GetTransactionHash())

		Topics := log.GetTopics()
		// 获取 from topic[1] 和 to topic[2] 地址
		recv := common.HexToAddress(Topics[1]).String()
		logrus.Infof("recv: %s", recv)

		weiamount, _ := plugin.HexToDecimal(log.GetData())
		logrus.Infof("recv amount: %v", utils.WeiToEth(weiamount))
		if err := core.InviterHandle.UpdateLpRewards(recv, utils.WeiToEth(weiamount)); err != nil {
			logrus.Errorf("Lp rewards Write to database rewards: %v", err)
		}
		blockscan.ScanType = 1
		blockscan.BlockNumber = int64(log.GetBlockNum() + 1)
		if err := blockscan.UptadeBlockNumber(); err != nil {
			logrus.Errorf("Lp rewards Write to database block: %v", err)
		}
	}

	return nil
}

func ScanLpRewards() {

	api := "https://blissful-damp-owl.bsc-testnet.discover.quiknode.pro/04dff7903bb2526b98ec1a883d9dbc6b45bb3b6e/"
	startBlockNum := 26846722
	blockNumber := int(blockscan.GetNumber(1))
	if blockNumber < startBlockNum {
		blockNumber = startBlockNum
	}
	contractAdx := CTF_CONTRACT
	ethrpc = rpc.NewEthRPC(api)
	logrus.Info("ScanLpRewards start blocknumber %v", blockNumber)
	receiptLogWatcher := ethereum_watcher.NewReceiptLogWatcher(
		context.Background(),
		api,
		blockNumber,
		contractAdx,
		ReceiveLpFee,
		lpRewardsHandle,
		ethereum_watcher.ReceiptLogWatcherConfig{
			StepSizeForBigLag:               50,
			IntervalForPollingNewBlockInSec: 3,
			RPCMaxRetry:                     3,
			ReturnForBlockWithNoReceiptLog:  false,
		},
	)

	receiptLogWatcher.Run()

}

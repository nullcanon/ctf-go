// 统计交易量

package blockscan

// 同步链上数据，记录数据库

import (
	"context"
	"ctf/utils"
	// "fmt"
	// "github.com/HydroProtocol/ethereum-watcher/plugin"
	ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
	"github.com/sirupsen/logrus"

	// "github.com/shopspring/decimal"
	"github.com/HydroProtocol/ethereum-watcher/blockchain"
	"github.com/HydroProtocol/ethereum-watcher/plugin"
	"github.com/HydroProtocol/ethereum-watcher/rpc"
	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/ethdb/leveldb"
)

const (
	CONTRACT = "0x29a3DAa1bf8DE08a8afE3e29E36Fa2797f7a37b5"
	CTF_CONTRACT = "0x00000000000000000000000029a3daa1bf8de08a8afe3e29e36fa2797f7a37b5"
	USDT_BLACK = "0x000000000000000000000000b5151a092593ca413fc8782cd92904f4f3d5f2f8"
	USDT = "0x8538d1641ad855db9e36fc1c7dc84236f104bb4a"
	CTF_USDT_PAIR = "0x0000000000000000000000000baf10bcf766f47f5f35877799b419792be1cb5f"
)

var ERC20_transfer = []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}

var ethrpc *rpc.EthBlockChainRPC

func tradeVolumeHandle(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {
	logrus.Infof("len: %v", len(receiptLogs))
	
	for _, log := range receiptLogs {

		logrus.Infof("BlockNum: %v", log.GetBlockNum())
		
		Topics := log.GetTopics()
		// 获取 from topic[1] 和 to topic[2] 地址
		from := Topics[1]
		to := Topics[2]
		var user string
		if (from == CTF_USDT_PAIR || to == CTF_USDT_PAIR) && !(to == CTF_CONTRACT || from == CTF_CONTRACT){
			logrus.Infof("topic: %s", Topics[0])
			logrus.Infof("hash: %s", log.GetTransactionHash())
			logrus.Infof("from: %s", common.HexToAddress(Topics[1]))
			logrus.Infof("to: %s", common.HexToAddress(Topics[2]))
			ctfWeiamount, _ := plugin.HexToDecimal(log.GetData())
			logrus.Infof("amount: %v", utils.WeiToEth(ctfWeiamount))

			// 获取usdt单位交易量
			txReceipt, err := ethrpc.GetTransactionReceipt(log.GetTransactionHash())
			if err == nil {
				subLogs := txReceipt.GetLogs()
				for _, subLog := range subLogs {
					if USDT == subLog.GetAddress() {
						subTopics := subLog.GetTopics()
						if (subTopics[1] == CTF_USDT_PAIR || subTopics[2] == CTF_USDT_PAIR) && subTopics[2] != USDT_BLACK{
							usdtWeiamount, _ := plugin.HexToDecimal(subLog.GetData())
							logrus.Infof("Usdt amount: %v", utils.WeiToEth(usdtWeiamount))
						}
					}
				}
			}
			if from == CTF_USDT_PAIR {
				user = to
			} else {
				user = from
			}
			logrus.Infof("user: %v", common.HexToAddress(user))

		}


		// 当 from 地址是 pair 地址，那么to地址是用户地址，则该笔交易为买入，获取买入交易的ctf数额
		// 获取该笔交易hash，解析usdt转账的log，获取转入pair地址的数量为当前交易量

	}
	return nil
}

func ScanTradeVolume() {

	api := "https://blissful-damp-owl.bsc-testnet.discover.quiknode.pro/04dff7903bb2526b98ec1a883d9dbc6b45bb3b6e/"
	startBlockNum := 26845584
	contractAdx := CONTRACT
	ethrpc = rpc.NewEthRPC(api)

	receiptLogWatcher := ethereum_watcher.NewReceiptLogWatcher(
		context.Background(),
		api,
		startBlockNum,
		contractAdx,
		ERC20_transfer,
		tradeVolumeHandle,
		ethereum_watcher.ReceiptLogWatcherConfig{
			StepSizeForBigLag:               50,
			IntervalForPollingNewBlockInSec: 3,
			RPCMaxRetry:                     3,
			ReturnForBlockWithNoReceiptLog:  false,
		},
	)

	receiptLogWatcher.Run()

}

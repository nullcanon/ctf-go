// 统计交易量

package blockscan

// 同步链上数据，记录数据库

import (
	"context"
	"ctf/core"
	"ctf/models"
	"ctf/utils"

	"github.com/shopspring/decimal"

	ethereum_watcher "github.com/HydroProtocol/ethereum-watcher"
	"github.com/HydroProtocol/ethereum-watcher/blockchain"
	"github.com/HydroProtocol/ethereum-watcher/plugin"
	"github.com/HydroProtocol/ethereum-watcher/rpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

const (
	CTF_CONTRACT  = "0x29a3DAa1bf8DE08a8afE3e29E36Fa2797f7a37b5"
	USDT_BLACK    = "0xb5151A092593ca413fC8782Cd92904F4f3d5F2f8" //代币合约中间地址
	USDT          = "0x8538d1641ad855db9e36fc1c7dc84236f104bb4a"
	CTF_USDT_PAIR = "0x0bAF10bCF766f47F5F35877799B419792bE1cB5f"
)

var ERC20_transfer = []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}

var ethrpc *rpc.EthBlockChainRPC
var blockscan models.BlockScan

func tradeVolumeHandle(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {
	logrus.Infof("len: %v", len(receiptLogs))
	var usdtAmount decimal.Decimal

	for _, log := range receiptLogs {

		logrus.Infof("BlockNum: %v", log.GetBlockNum())

		Topics := log.GetTopics()
		// 获取 from topic[1] 和 to topic[2] 地址
		from := common.HexToAddress(Topics[1]).String()
		to := common.HexToAddress(Topics[2]).String()
		var user string
		if (from == CTF_USDT_PAIR || to == CTF_USDT_PAIR) && !(to == CTF_CONTRACT || from == CTF_CONTRACT) {
			logrus.Infof("topic: %s", Topics[0])
			logrus.Infof("hash: %s", log.GetTransactionHash())
			logrus.Infof("from: %s", from)
			logrus.Infof("to: %s", to)
			ctfWeiamount, _ := plugin.HexToDecimal(log.GetData())
			logrus.Infof("amount: %v", utils.WeiToEth(ctfWeiamount))

			// 获取usdt单位交易量
			txReceipt, err := ethrpc.GetTransactionReceipt(log.GetTransactionHash())
			if err == nil {
				subLogs := txReceipt.GetLogs()
				for _, subLog := range subLogs {
					if USDT == subLog.GetAddress() {
						subTopics := subLog.GetTopics()
						subFrom := common.HexToAddress(subTopics[1]).String()
						subTo := common.HexToAddress(subTopics[2]).String()
						if (subFrom == CTF_USDT_PAIR || subTo == CTF_USDT_PAIR) && subTo != USDT_BLACK {
							usdtWeiamount, _ := plugin.HexToDecimal(subLog.GetData())
							usdtAmount = utils.WeiToEth(usdtWeiamount)
							logrus.Infof("Usdt amount: %v", usdtAmount)
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

			//  同步数据
			if err := core.InviterHandle.UpdateTradeVolume(common.HexToAddress(user).String(), usdtAmount); err != nil {
				logrus.Errorf("Trade Volume Write to database: %v", err)
			}

			blockscan.ScanType = 0
			blockscan.BlockNumber = int64(log.GetBlockNum() + 1)
			if err := blockscan.UptadeBlockNumber(); err != nil {
				logrus.Errorf("Trade Volume Write to database block: %v", err)
			}
		}
	}

	return nil
}

func ScanTradeVolume() {

	api := "https://blissful-damp-owl.bsc-testnet.discover.quiknode.pro/04dff7903bb2526b98ec1a883d9dbc6b45bb3b6e/"
	startBlockNum := 26845584
	blockNumber := int(blockscan.GetNumber(0))
	if blockNumber < startBlockNum {
		blockNumber = startBlockNum
	}

	contractAdx := CTF_CONTRACT
	ethrpc = rpc.NewEthRPC(api)

	receiptLogWatcher := ethereum_watcher.NewReceiptLogWatcher(
		context.Background(),
		api,
		blockNumber,
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
